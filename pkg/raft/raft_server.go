/*
Copyright 2017 The GoStor Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package raft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/goraft/raft"
	"github.com/gorilla/mux"
	"github.com/gostor/gofs/pkg/cache"
)

type RaftServer struct {
	peers      []string // initial peers to join with
	raftServer raft.Server
	dataDir    string
	httpAddr   string
	router     *mux.Router
	httpServer *http.Server
}

func NewRaftServer(peers []string, httpAddr string, dataDir string, pulseSeconds int, r *mux.Router, httpServer *http.Server, cache cache.Cache) (*RaftServer, error) {
	s := &RaftServer{
		peers:      peers,
		httpAddr:   httpAddr,
		dataDir:    dataDir,
		router:     r,
		httpServer: httpServer,
	}

	if log.GetLevel() == log.DebugLevel {
		raft.SetLogLevel(2)
	}

	var err error
	transporter := raft.NewHTTPTransporter("/cluster", 0)
	transporter.Transport.MaxIdleConnsPerHost = 1024
	log.Debugf("Starting RaftServer with IP:%v:", httpAddr)

	// Clear old cluster configurations if peers are changed
	if oldPeers, changed := isPeersChanged(s.dataDir, httpAddr, s.peers); changed {
		log.Infof("Peers Change: %v => %v", oldPeers, s.peers)
		os.RemoveAll(path.Join(s.dataDir, "conf"))
		os.RemoveAll(path.Join(s.dataDir, "log"))
		os.RemoveAll(path.Join(s.dataDir, "snapshot"))
	}

	s.raftServer, err = raft.NewServer(s.httpAddr, s.dataDir, transporter, nil, cache, "")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	transporter.Install(s.raftServer, s)
	s.raftServer.SetHeartbeatInterval(500 * time.Millisecond)
	s.raftServer.SetElectionTimeout(time.Duration(pulseSeconds) * 500 * time.Millisecond)
	s.raftServer.Start()

	if len(s.peers) > 0 {
		// Join to leader if specified.
		for {
			log.Infoln("Joining cluster:", strings.Join(s.peers, ","))
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			firstJoinError := s.Join(s.peers)
			if firstJoinError != nil {
				log.Infoln("No existing server found. Starting as leader in the new cluster.")
				_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
					Name:             s.raftServer.Name(),
					ConnectionString: "http://" + s.httpAddr,
				})
				if err != nil {
					log.Infoln(err)
				} else {
					break
				}
			} else {
				break
			}
		}
	} else if s.raftServer.IsLogEmpty() {
		// Initialize the server by joining itself.
		log.Infoln("Initializing new cluster")

		_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
			Name:             s.raftServer.Name(),
			ConnectionString: "http://" + s.httpAddr,
		})

		if err != nil {
			log.Infoln(err)
			return nil, err
		}

	} else {
		log.Infoln("Old conf,log,snapshot should have been removed.")
	}

	return s, nil
}

func (s *RaftServer) Do(command raft.Command) (interface{}, error) {
	return s.raftServer.Do(command)
}

func (s *RaftServer) Leader() (string, error) {
	if s.raftServer == nil {
		return "", fmt.Errorf("RaftServer is not ready")
	}
	return s.raftServer.Leader(), nil
}

func (s *RaftServer) Peers() (members []string) {
	peers := s.raftServer.Peers()

	for _, p := range peers {
		members = append(members, strings.TrimPrefix(p.ConnectionString, "http://"))
	}

	return
}

// This is a hack around Gorilla mux not providing the correct net/http
// HandleFunc() interface.
func (s *RaftServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

// Returns the connection string.
func (s *RaftServer) connectionString() string {
	//return fmt.Sprintf("http://%s:%d", s.host, s.port)
	return s.httpAddr
}

func isPeersChanged(dir string, self string, peers []string) (oldPeers []string, changed bool) {
	confPath := path.Join(dir, "conf")
	// open conf file
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		return oldPeers, true
	}
	conf := &raft.Config{}
	if err = json.Unmarshal(b, conf); err != nil {
		return oldPeers, true
	}

	for _, p := range conf.Peers {
		oldPeers = append(oldPeers, strings.TrimPrefix(p.ConnectionString, "http://"))
	}
	oldPeers = append(oldPeers, self)

	sort.Strings(peers)
	sort.Strings(oldPeers)

	return oldPeers, !reflect.DeepEqual(peers, oldPeers)

}

// Join joins an existing cluster.
func (s *RaftServer) Join(peers []string) error {
	command := &raft.DefaultJoinCommand{
		Name:             s.raftServer.Name(),
		ConnectionString: "http://" + s.httpAddr,
	}

	var err error
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(command)
	for _, m := range peers {
		if m == s.httpAddr {
			continue
		}
		target := fmt.Sprintf("http://%s/cluster/join", strings.TrimSpace(m))
		log.Infoln("Attempting to connect to:", target)

		err = postFollowingOneRedirect(target, "application/json", b)

		if err != nil {
			log.Infoln("Post returned error: ", err.Error())
			if _, ok := err.(*url.Error); ok {
				// If we receive a network error try the next member
				continue
			}
		} else {
			return nil
		}
	}

	return errors.New("Could not connect to any cluster peers")
}

// a workaround because http POST following redirection misses request body
func postFollowingOneRedirect(target string, contentType string, b bytes.Buffer) error {
	backupReader := bytes.NewReader(b.Bytes())
	resp, err := http.Post(target, contentType, &b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	data, _ := ioutil.ReadAll(resp.Body)
	reply := string(data)

	if strings.HasPrefix(reply, "\"http") {
		urlStr := reply[1 : len(reply)-1]

		log.Infoln("Post redirected to ", urlStr)
		resp2, err2 := http.Post(urlStr, contentType, backupReader)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		data, _ = ioutil.ReadAll(resp2.Body)
		statusCode = resp2.StatusCode
	}

	log.Infoln("Post returned status: ", statusCode, string(data))
	if statusCode != http.StatusOK {
		return errors.New(string(data))
	}

	return nil
}
