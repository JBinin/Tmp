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
package authorization

import (
	"sync"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
)

// Plugin allows third party plugins to authorize requests and responses
// in the context of docker API
type Plugin interface {
	// Name returns the registered plugin name
	Name() string

	// AuthZRequest authorizes the request from the client to the daemon
	AuthZRequest(*Request) (*Response, error)

	// AuthZResponse authorizes the response from the daemon to the client
	AuthZResponse(*Request) (*Response, error)
}

// newPlugins constructs and initializes the authorization plugins based on plugin names
func newPlugins(names []string) []Plugin {
	plugins := []Plugin{}
	pluginsMap := make(map[string]struct{})
	for _, name := range names {
		if _, ok := pluginsMap[name]; ok {
			continue
		}
		pluginsMap[name] = struct{}{}
		plugins = append(plugins, newAuthorizationPlugin(name))
	}
	return plugins
}

var getter plugingetter.PluginGetter

// SetPluginGetter sets the plugingetter
func SetPluginGetter(pg plugingetter.PluginGetter) {
	getter = pg
}

// GetPluginGetter gets the plugingetter
func GetPluginGetter() plugingetter.PluginGetter {
	return getter
}

// authorizationPlugin is an internal adapter to docker plugin system
type authorizationPlugin struct {
	plugin *plugins.Client
	name   string
	once   sync.Once
}

func newAuthorizationPlugin(name string) Plugin {
	return &authorizationPlugin{name: name}
}

func (a *authorizationPlugin) Name() string {
	return a.name
}

func (a *authorizationPlugin) AuthZRequest(authReq *Request) (*Response, error) {
	if err := a.initPlugin(); err != nil {
		return nil, err
	}

	authRes := &Response{}
	if err := a.plugin.Call(AuthZApiRequest, authReq, authRes); err != nil {
		return nil, err
	}

	return authRes, nil
}

func (a *authorizationPlugin) AuthZResponse(authReq *Request) (*Response, error) {
	if err := a.initPlugin(); err != nil {
		return nil, err
	}

	authRes := &Response{}
	if err := a.plugin.Call(AuthZApiResponse, authReq, authRes); err != nil {
		return nil, err
	}

	return authRes, nil
}

// initPlugin initializes the authorization plugin if needed
func (a *authorizationPlugin) initPlugin() error {
	// Lazy loading of plugins
	var err error
	a.once.Do(func() {
		if a.plugin == nil {
			var plugin plugingetter.CompatPlugin
			var e error

			if pg := GetPluginGetter(); pg != nil {
				plugin, e = pg.Get(a.name, AuthZApiImplements, plugingetter.Lookup)
			} else {
				plugin, e = plugins.Get(a.name, AuthZApiImplements)
			}
			if e != nil {
				err = e
				return
			}
			a.plugin = plugin.Client()
		}
	})
	return err
}
