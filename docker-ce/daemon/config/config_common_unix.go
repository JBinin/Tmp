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
// +build solaris linux freebsd

package config

import (
	"net"

	"github.com/docker/docker/api/types"
)

// CommonUnixConfig defines configuration of a docker daemon that is
// common across Unix platforms.
type CommonUnixConfig struct {
	ExecRoot          string                   `json:"exec-root,omitempty"`
	ContainerdAddr    string                   `json:"containerd,omitempty"`
	Runtimes          map[string]types.Runtime `json:"runtimes,omitempty"`
	DefaultRuntime    string                   `json:"default-runtime,omitempty"`
	DefaultInitBinary string                   `json:"default-init,omitempty"`
}

type commonUnixBridgeConfig struct {
	DefaultIP                   net.IP `json:"ip,omitempty"`
	IP                          string `json:"bip,omitempty"`
	DefaultGatewayIPv4          net.IP `json:"default-gateway,omitempty"`
	DefaultGatewayIPv6          net.IP `json:"default-gateway-v6,omitempty"`
	InterContainerCommunication bool   `json:"icc,omitempty"`
}

// GetRuntime returns the runtime path and arguments for a given
// runtime name
func (conf *Config) GetRuntime(name string) *types.Runtime {
	conf.Lock()
	defer conf.Unlock()
	if rt, ok := conf.Runtimes[name]; ok {
		return &rt
	}
	return nil
}

// GetDefaultRuntimeName returns the current default runtime
func (conf *Config) GetDefaultRuntimeName() string {
	conf.Lock()
	rt := conf.DefaultRuntime
	conf.Unlock()

	return rt
}

// GetAllRuntimes returns a copy of the runtimes map
func (conf *Config) GetAllRuntimes() map[string]types.Runtime {
	conf.Lock()
	rts := conf.Runtimes
	conf.Unlock()
	return rts
}

// GetExecRoot returns the user configured Exec-root
func (conf *Config) GetExecRoot() string {
	return conf.ExecRoot
}

// GetInitPath returns the configure docker-init path
func (conf *Config) GetInitPath() string {
	conf.Lock()
	defer conf.Unlock()
	if conf.InitPath != "" {
		return conf.InitPath
	}
	return conf.DefaultInitBinary
}
