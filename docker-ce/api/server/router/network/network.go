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
package network

import (
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/daemon/cluster"
)

// networkRouter is a router to talk with the network controller
type networkRouter struct {
	backend Backend
	cluster *cluster.Cluster
	routes  []router.Route
}

// NewRouter initializes a new network router
func NewRouter(b Backend, c *cluster.Cluster) router.Router {
	r := &networkRouter{
		backend: b,
		cluster: c,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the network controller
func (r *networkRouter) Routes() []router.Route {
	return r.routes
}

func (r *networkRouter) initRoutes() {
	r.routes = []router.Route{
		// GET
		router.NewGetRoute("/networks", r.getNetworksList),
		router.NewGetRoute("/networks/", r.getNetworksList),
		router.NewGetRoute("/networks/{id:.+}", r.getNetwork),
		// POST
		router.NewPostRoute("/networks/create", r.postNetworkCreate),
		router.NewPostRoute("/networks/{id:.*}/connect", r.postNetworkConnect),
		router.NewPostRoute("/networks/{id:.*}/disconnect", r.postNetworkDisconnect),
		router.NewPostRoute("/networks/prune", r.postNetworksPrune),
		// DELETE
		router.NewDeleteRoute("/networks/{id:.*}", r.deleteNetwork),
	}
}
