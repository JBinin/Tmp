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
// +build windows

package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/libnetwork"
)

func (daemon *Daemon) setupLinkedContainers(container *container.Container) ([]string, error) {
	return nil, nil
}

// getSize returns real size & virtual size
func (daemon *Daemon) getSize(containerID string) (int64, int64) {
	// TODO Windows
	return 0, 0
}

func (daemon *Daemon) setupIpcDirs(container *container.Container) error {
	return nil
}

// TODO Windows: Fix Post-TP5. This is a hack to allow docker cp to work
// against containers which have volumes. You will still be able to cp
// to somewhere on the container drive, but not to any mounted volumes
// inside the container. Without this fix, docker cp is broken to any
// container which has a volume, regardless of where the file is inside the
// container.
func (daemon *Daemon) mountVolumes(container *container.Container) error {
	return nil
}

func detachMounted(path string) error {
	return nil
}

func killProcessDirectly(container *container.Container) error {
	return nil
}

func isLinkable(child *container.Container) bool {
	return false
}

func enableIPOnPredefinedNetwork() bool {
	return true
}

func (daemon *Daemon) isNetworkHotPluggable() bool {
	return false
}

func setupPathsAndSandboxOptions(container *container.Container, sboxOptions *[]libnetwork.SandboxOption) error {
	return nil
}

func initializeNetworkingPaths(container *container.Container, nc *container.Container) {
}
