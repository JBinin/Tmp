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
package volume

import "github.com/docker/docker/api/server/router"

// volumeRouter is a router to talk with the volumes controller
type volumeRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new volume router
func NewRouter(b Backend) router.Router {
	r := &volumeRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the volumes controller
func (r *volumeRouter) Routes() []router.Route {
	return r.routes
}

func (r *volumeRouter) initRoutes() {
	r.routes = []router.Route{
		// GET
		router.NewGetRoute("/volumes", r.getVolumesList),
		router.NewGetRoute("/volumes/{name:.*}", r.getVolumeByName),
		// POST
		router.NewPostRoute("/volumes/create", r.postVolumesCreate),
		router.NewPostRoute("/volumes/prune", r.postVolumesPrune),
		// DELETE
		router.NewDeleteRoute("/volumes/{name:.*}", r.deleteVolumes),
	}
}
