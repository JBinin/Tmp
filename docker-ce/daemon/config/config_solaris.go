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
	"github.com/spf13/pflag"
)

var (
	defaultPidFile = "/system/volatile/docker/docker.pid"
	defaultGraph   = "/var/lib/docker"
	defaultExec    = "zones"
)

// Config defines the configuration of a docker daemon.
// These are the configuration settings that you pass
// to the docker daemon when you launch it with say: `docker -d -e lxc`
type Config struct {
	CommonConfig

	// These fields are common to all unix platforms.
	CommonUnixConfig
}

// BridgeConfig stores all the bridge driver specific
// configuration.
type BridgeConfig struct {
	commonBridgeConfig

	// Fields below here are platform specific.
	commonUnixBridgeConfig
}

// IsSwarmCompatible defines if swarm mode can be enabled in this config
func (conf *Config) IsSwarmCompatible() error {
	return nil
}
