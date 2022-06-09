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
// +build solaris

package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/runconfig"
	"github.com/docker/libnetwork"
)

func (daemon *Daemon) setupLinkedContainers(container *container.Container) ([]string, error) {
	return nil, nil
}

func (daemon *Daemon) setupIpcDirs(container *container.Container) error {
	return nil
}

func killProcessDirectly(container *container.Container) error {
	return nil
}

func detachMounted(path string) error {
	return nil
}

func isLinkable(child *container.Container) bool {
	// A container is linkable only if it belongs to the default network
	_, ok := child.NetworkSettings.Networks[runconfig.DefaultDaemonNetworkMode().NetworkName()]
	return ok
}

func enableIPOnPredefinedNetwork() bool {
	return false
}

func (daemon *Daemon) isNetworkHotPluggable() bool {
	return false
}

func setupPathsAndSandboxOptions(container *container.Container, sboxOptions *[]libnetwork.SandboxOption) error {
	return nil
}

func initializeNetworkingPaths(container *container.Container, nc *container.Container) {
}
