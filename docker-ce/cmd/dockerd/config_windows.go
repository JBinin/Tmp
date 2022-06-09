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
package main

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/daemon/config"
	"github.com/spf13/pflag"
)

var (
	defaultPidFile string
	defaultGraph   = filepath.Join(os.Getenv("programdata"), "docker")
)

// installConfigFlags adds flags to the pflag.FlagSet to configure the daemon
func installConfigFlags(conf *config.Config, flags *pflag.FlagSet) {
	// First handle install flags which are consistent cross-platform
	installCommonConfigFlags(conf, flags)

	// Then platform-specific install flags.
	flags.StringVar(&conf.BridgeConfig.FixedCIDR, "fixed-cidr", "", "IPv4 subnet for fixed IPs")
	flags.StringVarP(&conf.BridgeConfig.Iface, "bridge", "b", "", "Attach containers to a virtual switch")
	flags.StringVarP(&conf.SocketGroup, "group", "G", "", "Users or groups that can access the named pipe")
}
