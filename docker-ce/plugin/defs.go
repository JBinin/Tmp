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

import (
	"sync"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/docker/plugin/v2"
)

// Store manages the plugin inventory in memory and on-disk
type Store struct {
	sync.RWMutex
	plugins map[string]*v2.Plugin
	/* handlers are necessary for transition path of legacy plugins
	 * to the new model. Legacy plugins use Handle() for registering an
	 * activation callback.*/
	handlers map[string][]func(string, *plugins.Client)
}

// NewStore creates a Store.
func NewStore(libRoot string) *Store {
	return &Store{
		plugins:  make(map[string]*v2.Plugin),
		handlers: make(map[string][]func(string, *plugins.Client)),
	}
}
