/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package system

import (
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/daemon/cluster"
)

// systemRouter provides information about the Docker system overall.
// It gathers information about host, daemon and container events.
type systemRouter struct {
	backend Backend
	cluster *cluster.Cluster
	routes  []router.Route
}

// NewRouter initializes a new system router
func NewRouter(b Backend, c *cluster.Cluster) router.Router {
	r := &systemRouter{
		backend: b,
		cluster: c,
	}

	r.routes = []router.Route{
		router.NewOptionsRoute("/{anyroute:.*}", optionsHandler),
		router.NewGetRoute("/_ping", pingHandler),
		router.Cancellable(router.NewGetRoute("/events", r.getEvents)),
		router.NewGetRoute("/info", r.getInfo),
		router.NewGetRoute("/version", r.getVersion),
		router.NewGetRoute("/system/df", r.getDiskUsage),
		router.NewPostRoute("/auth", r.postAuth),
	}

	return r
}

// Routes returns all the API routes dedicated to the docker system
func (s *systemRouter) Routes() []router.Route {
	return s.routes
}
