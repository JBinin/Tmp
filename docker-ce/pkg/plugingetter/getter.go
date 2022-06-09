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
package plugingetter

import "github.com/docker/docker/pkg/plugins"

const (
	// Lookup doesn't update RefCount
	Lookup = 0
	// Acquire increments RefCount
	Acquire = 1
	// Release decrements RefCount
	Release = -1
)

// CompatPlugin is an abstraction to handle both v2(new) and v1(legacy) plugins.
type CompatPlugin interface {
	Client() *plugins.Client
	Name() string
	BasePath() string
	IsV1() bool
}

// CountedPlugin is a plugin which is reference counted.
type CountedPlugin interface {
	Acquire()
	Release()
	CompatPlugin
}

// PluginGetter is the interface implemented by Store
type PluginGetter interface {
	Get(name, capability string, mode int) (CompatPlugin, error)
	GetAllByCap(capability string) ([]CompatPlugin, error)
	GetAllManagedPluginsByCap(capability string) []CompatPlugin
	Handle(capability string, callback func(string, *plugins.Client))
}
