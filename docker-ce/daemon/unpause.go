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

// ContainerUnpause unpauses a container
func (daemon *Daemon) ContainerUnpause(name string) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	if err := daemon.containerUnpause(container); err != nil {
		return err
	}

	return nil
}

// containerUnpause resumes the container execution after the container is paused.
func (daemon *Daemon) containerUnpause(container *container.Container) error {
	container.Lock()
	defer container.Unlock()

	// We cannot unpause the container which is not paused
	if !container.Paused {
		return fmt.Errorf("Container %s is not paused", container.ID)
	}

	if err := daemon.containerd.Resume(container.ID); err != nil {
		return fmt.Errorf("Cannot unpause container %s: %s", container.ID, err)
	}

	return nil
}
