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
package graphdriver

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/plugin/v2"
)

type pluginClient interface {
	// Call calls the specified method with the specified arguments for the plugin.
	Call(string, interface{}, interface{}) error
	// Stream calls the specified method with the specified arguments for the plugin and returns the response IO stream
	Stream(string, interface{}) (io.ReadCloser, error)
	// SendFile calls the specified method, and passes through the IO stream
	SendFile(string, io.Reader, interface{}) error
}

func lookupPlugin(name string, pg plugingetter.PluginGetter, config Options) (Driver, error) {
	if !config.ExperimentalEnabled {
		return nil, fmt.Errorf("graphdriver plugins are only supported with experimental mode")
	}
	pl, err := pg.Get(name, "GraphDriver", plugingetter.Acquire)
	if err != nil {
		return nil, fmt.Errorf("Error looking up graphdriver plugin %s: %v", name, err)
	}
	return newPluginDriver(name, pl, config)
}

func newPluginDriver(name string, pl plugingetter.CompatPlugin, config Options) (Driver, error) {
	home := config.Root
	if !pl.IsV1() {
		if p, ok := pl.(*v2.Plugin); ok {
			if p.PropagatedMount != "" {
				home = p.PluginObj.Config.PropagatedMount
			}
		}
	}
	proxy := &graphDriverProxy{name, pl}
	return proxy, proxy.Init(filepath.Join(home, name), config.DriverOptions, config.UIDMaps, config.GIDMaps)
}
