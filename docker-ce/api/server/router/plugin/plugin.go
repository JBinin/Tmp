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
package plugin

import "github.com/docker/docker/api/server/router"

// pluginRouter is a router to talk with the plugin controller
type pluginRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router {
	r := &pluginRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the plugin controller
func (r *pluginRouter) Routes() []router.Route {
	return r.routes
}

func (r *pluginRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRoute("/plugins", r.listPlugins),
		router.NewGetRoute("/plugins/{name:.*}/json", r.inspectPlugin),
		router.NewGetRoute("/plugins/privileges", r.getPrivileges),
		router.NewDeleteRoute("/plugins/{name:.*}", r.removePlugin),
		router.NewPostRoute("/plugins/{name:.*}/enable", r.enablePlugin), // PATCH?
		router.NewPostRoute("/plugins/{name:.*}/disable", r.disablePlugin),
		router.Cancellable(router.NewPostRoute("/plugins/pull", r.pullPlugin)),
		router.Cancellable(router.NewPostRoute("/plugins/{name:.*}/push", r.pushPlugin)),
		router.Cancellable(router.NewPostRoute("/plugins/{name:.*}/upgrade", r.upgradePlugin)),
		router.NewPostRoute("/plugins/{name:.*}/set", r.setPlugin),
		router.NewPostRoute("/plugins/create", r.createPlugin),
	}
}
