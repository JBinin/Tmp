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
// +build !windows

package daemon

import (
	"fmt"

	"github.com/docker/docker/container"
	"github.com/docker/docker/libcontainerd"
)

func (daemon *Daemon) getLibcontainerdCreateOptions(container *container.Container) ([]libcontainerd.CreateOption, error) {
	createOptions := []libcontainerd.CreateOption{}

	// Ensure a runtime has been assigned to this container
	if container.HostConfig.Runtime == "" {
		container.HostConfig.Runtime = daemon.configStore.GetDefaultRuntimeName()
		container.ToDisk()
	}

	rt := daemon.configStore.GetRuntime(container.HostConfig.Runtime)
	if rt == nil {
		return nil, fmt.Errorf("no such runtime '%s'", container.HostConfig.Runtime)
	}
	if UsingSystemd(daemon.configStore) {
		rt.Args = append(rt.Args, "--systemd-cgroup=true")
	}
	createOptions = append(createOptions, libcontainerd.WithRuntime(rt.Path, rt.Args))

	return createOptions, nil
}
