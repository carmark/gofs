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
package metadata

import (
	"context"
	"net/http"

	"github.com/gostor/gofs/pkg/apiserver/httputils"
	"github.com/gostor/gofs/pkg/apiserver/router"
	"github.com/gostor/gofs/pkg/master"
)

// mdRouter is a router to talk with the metadata controller
type mdRouter struct {
	routes []router.Route
	master *master.Master
}

// NewRouter initializes a new container router
func NewRouter(master *master.Master) router.Router {
	r := &mdRouter{master: master}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the container controller
func (r *mdRouter) Routes() []router.Route {
	return r.routes
}

// initRoutes initializes the routes in metadata router
func (r *mdRouter) initRoutes() {
	r.routes = []router.Route{
		// GET
		router.NewGetRoute("/{path:.*}", r.getMetadataOperation),
		// POST
		router.NewPostRoute("/{path:.*}", r.postMetadataOperation),
		// PUT
		router.NewPostRoute("/{path:.*}", r.putMetadataOperation),
		// DELETE
		router.NewDeleteRoute("/{path:.*}", r.deleteMetadataOperation),
	}
}

func (r *mdRouter) getMetadataOperation(ctx context.Context, w http.ResponseWriter, req *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	path := vars["path"]
	operation := r.Form.Get("op")

	resp, err := r.master.GetPathHandler(path, operation)
	if err != nil {
		return err
	}
	httputils.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (r *mdRouter) postMetadataOperation(ctx context.Context, w http.ResponseWriter, req *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	path := vars["path"]
	operation := r.Form.Get("op")

	err := r.master.PostPathHandler(path, operation)
	if err != nil {
		return err
	}
	httputils.WriteJSON(w, http.StatusOK, nil)
	return nil
}

func (r *mdRouter) putMetadataOperation(ctx context.Context, w http.ResponseWriter, req *http.Request, vars map[string]string) error {
	return nil
}

func (r *mdRouter) deleteMetadataOperation(ctx context.Context, w http.ResponseWriter, req *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	path := vars["path"]
	operation := r.Form.Get("op")
	recursive := httputils.BoolValueOrDefault(req, "recursive", false)

	err := r.master.DeletePathHandler(path, operation, recursive)
	if err != nil {
		return err
	}
	httputils.WriteJSON(w, http.StatusOK, nil)
	return nil
}
