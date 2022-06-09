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
package daemon

import (
	"fmt"

	"github.com/docker/docker/container"
)

// ContainerPause pauses a container
func (daemon *Daemon) ContainerPause(name string) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	if err := daemon.containerPause(container); err != nil {
		return err
	}

	return nil
}

// containerPause pauses the container execution without stopping the process.
// The execution can be resumed by calling containerUnpause.
func (daemon *Daemon) containerPause(container *container.Container) error {
	container.Lock()
	defer container.Unlock()

	// We cannot Pause the container which is not running
	if !container.Running {
		return errNotRunning{container.ID}
	}

	// We cannot Pause the container which is already paused
	if container.Paused {
		return fmt.Errorf("Container %s is already paused", container.ID)
	}

	// We cannot Pause the container which is restarting
	if container.Restarting {
		return errContainerIsRestarting(container.ID)
	}

	if err := daemon.containerd.Pause(container.ID); err != nil {
		return fmt.Errorf("Cannot pause container %s: %s", container.ID, err)
	}

	return nil
}
