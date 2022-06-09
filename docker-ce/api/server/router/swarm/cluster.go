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
package swarm

import "github.com/docker/docker/api/server/router"

// swarmRouter is a router to talk with the build controller
type swarmRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new build router
func NewRouter(b Backend) router.Router {
	r := &swarmRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the swarm controller
func (sr *swarmRouter) Routes() []router.Route {
	return sr.routes
}

func (sr *swarmRouter) initRoutes() {
	sr.routes = []router.Route{
		router.NewPostRoute("/swarm/init", sr.initCluster),
		router.NewPostRoute("/swarm/join", sr.joinCluster),
		router.NewPostRoute("/swarm/leave", sr.leaveCluster),
		router.NewGetRoute("/swarm", sr.inspectCluster),
		router.NewGetRoute("/swarm/unlockkey", sr.getUnlockKey),
		router.NewPostRoute("/swarm/update", sr.updateCluster),
		router.NewPostRoute("/swarm/unlock", sr.unlockCluster),
		router.NewGetRoute("/services", sr.getServices),
		router.NewGetRoute("/services/{id}", sr.getService),
		router.NewPostRoute("/services/create", sr.createService),
		router.NewPostRoute("/services/{id}/update", sr.updateService),
		router.NewDeleteRoute("/services/{id}", sr.removeService),
		router.Experimental(router.Cancellable(router.NewGetRoute("/services/{id}/logs", sr.getServiceLogs))),
		router.NewGetRoute("/nodes", sr.getNodes),
		router.NewGetRoute("/nodes/{id}", sr.getNode),
		router.NewDeleteRoute("/nodes/{id}", sr.removeNode),
		router.NewPostRoute("/nodes/{id}/update", sr.updateNode),
		router.NewGetRoute("/tasks", sr.getTasks),
		router.NewGetRoute("/tasks/{id}", sr.getTask),
		router.NewGetRoute("/secrets", sr.getSecrets),
		router.NewPostRoute("/secrets/create", sr.createSecret),
		router.NewDeleteRoute("/secrets/{id}", sr.removeSecret),
		router.NewGetRoute("/secrets/{id}", sr.getSecret),
		router.NewPostRoute("/secrets/{id}/update", sr.updateSecret),
	}
}
