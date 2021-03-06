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

package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/opts"
	"github.com/spf13/pflag"
)

var (
	defaultPidFile  = "/var/run/docker.pid"
	defaultGraph    = "/var/lib/docker"
	defaultExecRoot = "/var/run/docker"
)

// installUnixConfigFlags adds command-line options to the top-level flag parser for
// the current process that are common across Unix platforms.
func installUnixConfigFlags(conf *config.Config, flags *pflag.FlagSet) {
	conf.Runtimes = make(map[string]types.Runtime)

	flags.StringVarP(&conf.SocketGroup, "group", "G", "docker", "Group for the unix socket")
	flags.StringVar(&conf.BridgeConfig.IP, "bip", "", "Specify network bridge IP")
	flags.StringVarP(&conf.BridgeConfig.Iface, "bridge", "b", "", "Attach containers to a network bridge")
	flags.StringVar(&conf.BridgeConfig.FixedCIDR, "fixed-cidr", "", "IPv4 subnet for fixed IPs")
	flags.Var(opts.NewIPOpt(&conf.BridgeConfig.DefaultGatewayIPv4, ""), "default-gateway", "Container default gateway IPv4 address")
	flags.Var(opts.NewIPOpt(&conf.BridgeConfig.DefaultGatewayIPv6, ""), "default-gateway-v6", "Container default gateway IPv6 address")
	flags.BoolVar(&conf.BridgeConfig.InterContainerCommunication, "icc", true, "Enable inter-container communication")
	flags.Var(opts.NewIPOpt(&conf.BridgeConfig.DefaultIP, "0.0.0.0"), "ip", "Default IP when binding container ports")
	flags.Var(opts.NewNamedRuntimeOpt("runtimes", &conf.Runtimes, config.StockRuntimeName), "add-runtime", "Register an additional OCI compatible runtime")
	flags.StringVar(&conf.DefaultRuntime, "default-runtime", config.StockRuntimeName, "Default OCI runtime for containers")

}
