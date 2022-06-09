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
package runconfig

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/sysinfo"
)

// DefaultDaemonNetworkMode returns the default network stack the daemon should
// use.
func DefaultDaemonNetworkMode() container.NetworkMode {
	return container.NetworkMode("nat")
}

// IsPreDefinedNetwork indicates if a network is predefined by the daemon
func IsPreDefinedNetwork(network string) bool {
	return !container.NetworkMode(network).IsUserDefined()
}

// ValidateNetMode ensures that the various combinations of requested
// network settings are valid.
func ValidateNetMode(c *container.Config, hc *container.HostConfig) error {
	if hc == nil {
		return nil
	}
	parts := strings.Split(string(hc.NetworkMode), ":")
	if len(parts) > 1 {
		return fmt.Errorf("invalid --net: %s", hc.NetworkMode)
	}
	return nil
}

// ValidateIsolation performs platform specific validation of the
// isolation in the hostconfig structure. Windows supports 'default' (or
// blank), 'process', or 'hyperv'.
func ValidateIsolation(hc *container.HostConfig) error {
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil {
		return nil
	}
	if !hc.Isolation.IsValid() {
		return fmt.Errorf("invalid --isolation: %q. Windows supports 'default', 'process', or 'hyperv'", hc.Isolation)
	}
	return nil
}

// ValidateQoS performs platform specific validation of the Qos settings
func ValidateQoS(hc *container.HostConfig) error {
	return nil
}

// ValidateResources performs platform specific validation of the resource settings
func ValidateResources(hc *container.HostConfig, si *sysinfo.SysInfo) error {
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil {
		return nil
	}

	if hc.Resources.CPURealtimePeriod != 0 {
		return fmt.Errorf("invalid --cpu-rt-period: Windows does not support this feature")
	}
	if hc.Resources.CPURealtimeRuntime != 0 {
		return fmt.Errorf("invalid --cpu-rt-runtime: Windows does not support this feature")
	}
	return nil
}
