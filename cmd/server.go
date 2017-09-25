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

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/gostor/gofs/pkg/apiserver"
	"github.com/gostor/gofs/pkg/master"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	var host string
	var driver string
	var logLevel string
	var peers string
	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Setup a server",
		Long:  `Setup the GoFS's metadata server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createDaemon(host, driver, logLevel, peers)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&logLevel, "log", "info", "Log level")
	flags.StringVar(&host, "host", "tcp://127.0.0.1:9876", "Host for GoFS server")
	flags.StringVar(&peers, "join", "", "Peers")
	return cmd
}

func createDaemon(host, driver, level, peers string) error {
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "panic", "fatal", "error":
		log.SetLevel(log.ErrorLevel)
	default:
		return fmt.Errorf("unknown log level: %v", level)
	}
	hosts := []string{}
	if host != "" {
		hosts = append(hosts, host)
	}
	serverConfig := &apiserver.Config{
		Addrs: []apiserver.Addr{},
	}
	for _, protoAddr := range hosts {
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 {
			err := fmt.Errorf("bad format %s, expected PROTO://ADDR", protoAddr)
			log.Error(err)
			return err
		}
		serverConfig.Addrs = append(serverConfig.Addrs, apiserver.Addr{Proto: protoAddrParts[0], Addr: protoAddrParts[1]})
	}

	s, err := apiserver.New(serverConfig)
	if err != nil {
		log.Error(err)
		return err
	}
	os.Mkdir(filepath.Join(os.TempDir(), "gofs"), 0700)
	cfg := master.MasterConfig{
		Peers:       strings.Split(peers, ","),
		HttpAddr:    host,
		DataDir:     filepath.Join(os.TempDir(), "gofs"),
		PuleSeconds: 2,
		Name:        host,
		Router:      s.GetMuxRouter(),
		HttpServers: s.GetHttpServer(),
		CacheType:   "memory",
		CacheDir:    filepath.Join(os.TempDir(), "gofs", "cache"),
	}
	master, err := master.NewMaster(&cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	s.InitRouters(master)
	// The serve API routine never exits unless an error occurs
	// We need to start it as a goroutine and wait on it so
	// daemon doesn't exit
	serveAPIWait := make(chan error)
	go s.Wait(serveAPIWait)

	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGINT, syscall.SIGTERM)

	// Daemon is fully initialized and handling API traffic
	// Wait for serve API job to complete
	select {
	case errAPI := <-serveAPIWait:
		// If we have an error here it is unique to API (as daemonErr would have
		// exited the daemon process above)
		if errAPI != nil {
			log.Warnf("Shutting down due to ServeAPI error: %v", errAPI)
		}
	case <-stopAll:
		break
	}
	s.Close()
	return nil
}
