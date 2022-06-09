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
package config

import (
	"github.com/docker/docker/api/types"
)

// BridgeConfig stores all the bridge driver specific
// configuration.
type BridgeConfig struct {
	commonBridgeConfig
}

// Config defines the configuration of a docker daemon.
// These are the configuration settings that you pass
// to the docker daemon when you launch it with say: `dockerd -e windows`
type Config struct {
	CommonConfig

	// Fields below here are platform specific. (There are none presently
	// for the Windows daemon.)
}

// GetRuntime returns the runtime path and arguments for a given
// runtime name
func (conf *Config) GetRuntime(name string) *types.Runtime {
	return nil
}

// GetInitPath returns the configure docker-init path
func (conf *Config) GetInitPath() string {
	return ""
}

// GetDefaultRuntimeName returns the current default runtime
func (conf *Config) GetDefaultRuntimeName() string {
	return StockRuntimeName
}

// GetAllRuntimes returns a copy of the runtimes map
func (conf *Config) GetAllRuntimes() map[string]types.Runtime {
	return map[string]types.Runtime{}
}

// GetExecRoot returns the user configured Exec-root
func (conf *Config) GetExecRoot() string {
	return ""
}

// IsSwarmCompatible defines if swarm mode can be enabled in this config
func (conf *Config) IsSwarmCompatible() error {
	return nil
}
