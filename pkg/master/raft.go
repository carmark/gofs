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
package master

import (
	log "github.com/Sirupsen/logrus"
	"github.com/goraft/raft"
	"github.com/gostor/gofs/pkg/api"
)

// Handles incoming RAFT joins.
func (m *Master) JoinLeader(command *raft.DefaultJoinCommand) error {
	leader, err := m.RaftServer.Leader()
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Infof("Processing incoming join. Current Leader: %v, Self: %v, Peers: %v", leader, m.Name, m.RaftServer.Peers())
	log.Infof("join command from Name[%v], Connection[%v]", command.Name, command.ConnectionString)

	if _, err := m.RaftServer.Do(command); err != nil {
		return err
	}
	return nil
}

func (m *Master) RaftStatus() *api.RaftClusterStatusResponse {
	leader, err := m.RaftServer.Leader()
	if err != nil {
		log.Error(err)
		return nil
	}
	return &api.RaftClusterStatusResponse{
		IsLeader: m.IsLeader(),
		Leader:   leader,
		Peers:    m.RaftServer.Peers(),
	}
}
