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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/goraft/raft"
	"github.com/gostor/gofs/pkg/apiserver/httputils"
)

// Handles incoming RAFT joins.
func (r *raftRouter) joinHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, vars map[string]string) error {
	command := &raft.DefaultJoinCommand{}

	commandText, _ := ioutil.ReadAll(req.Body)
	log.Info("Command:", string(commandText))
	if err := json.NewDecoder(strings.NewReader(string(commandText))).Decode(&command); err != nil {
		log.Infof("Error decoding json message[command: %s]: %v", string(commandText), err)
		return err
	}

	err := r.master.JoinLeader(command)
	if err != nil {
		switch err {
		case raft.NotLeaderError:
			r.redirectToLeader(w, req)
		default:
			log.Infoln("Error processing join:", err)
			return err
		}
	}
	return nil
}

func (r *raftRouter) redirectToLeader(w http.ResponseWriter, req *http.Request) {
	if leader, err := r.master.RaftServer.Leader(); err == nil {
		learderLocation := "http://" + leader + req.URL.Path
		log.Infoln("Redirecting to", learderLocation)
		httputils.WriteJSON(w, http.StatusOK, learderLocation)
	} else {
		log.Infof("Error: Leader Unknown, %v", err)
		httputils.WriteError(w, fmt.Errorf("Leader unknown, %v", err))
	}
}

func (r *raftRouter) statusHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, vars map[string]string) error {
	ret := r.master.RaftStatus()
	return httputils.WriteJSON(w, http.StatusOK, ret)
}
