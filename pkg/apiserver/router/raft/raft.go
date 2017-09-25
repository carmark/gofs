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
	"github.com/gostor/gofs/pkg/apiserver/router"
	"github.com/gostor/gofs/pkg/master"
)

// raftRouter is a router to talk with the raft controller
type raftRouter struct {
	routes []router.Route
	master *master.Master
}

// NewRouter initializes a new container router
func NewRouter(master *master.Master) router.Router {
	r := &raftRouter{master: master}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the container controller
func (r *raftRouter) Routes() []router.Route {
	return r.routes
}

// initRoutes initializes the routes in metadata router
func (r *raftRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewPostRoute("/cluster/join", r.joinHandler),
		router.NewGetRoute("/cluster/status", r.statusHandler),
	}
}
