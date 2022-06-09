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
package checkpoint

import (
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
)

// checkpointRouter is a router to talk with the checkpoint controller
type checkpointRouter struct {
	backend Backend
	decoder httputils.ContainerDecoder
	routes  []router.Route
}

// NewRouter initializes a new checkpoint router
func NewRouter(b Backend, decoder httputils.ContainerDecoder) router.Router {
	r := &checkpointRouter{
		backend: b,
		decoder: decoder,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the checkpoint controller
func (r *checkpointRouter) Routes() []router.Route {
	return r.routes
}

func (r *checkpointRouter) initRoutes() {
	r.routes = []router.Route{
		router.Experimental(router.NewGetRoute("/containers/{name:.*}/checkpoints", r.getContainerCheckpoints)),
		router.Experimental(router.NewPostRoute("/containers/{name:.*}/checkpoints", r.postContainerCheckpoint)),
		router.Experimental(router.NewDeleteRoute("/containers/{name}/checkpoints/{checkpoint}", r.deleteContainerCheckpoint)),
	}
}
