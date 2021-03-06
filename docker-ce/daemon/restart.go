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

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/container"
)

// ContainerRestart stops and starts a container. It attempts to
// gracefully stop the container within the given timeout, forcefully
// stopping it if the timeout is exceeded. If given a negative
// timeout, ContainerRestart will wait forever until a graceful
// stop. Returns an error if the container cannot be found, or if
// there is an underlying error at any stage of the restart.
func (daemon *Daemon) ContainerRestart(name string, seconds *int) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}
	if seconds == nil {
		stopTimeout := container.StopTimeout()
		seconds = &stopTimeout
	}
	if err := daemon.containerRestart(container, *seconds); err != nil {
		return fmt.Errorf("Cannot restart container %s: %v", name, err)
	}
	return nil

}

// containerRestart attempts to gracefully stop and then start the
// container. When stopping, wait for the given duration in seconds to
// gracefully stop, before forcefully terminating the container. If
// given a negative duration, wait forever for a graceful stop.
func (daemon *Daemon) containerRestart(container *container.Container, seconds int) error {
	// Avoid unnecessarily unmounting and then directly mounting
	// the container when the container stops and then starts
	// again
	if err := daemon.Mount(container); err == nil {
		defer daemon.Unmount(container)
	}

	if container.IsRunning() {
		// set AutoRemove flag to false before stop so the container won't be
		// removed during restart process
		autoRemove := container.HostConfig.AutoRemove

		container.HostConfig.AutoRemove = false
		err := daemon.containerStop(container, seconds)
		// restore AutoRemove irrespective of whether the stop worked or not
		container.HostConfig.AutoRemove = autoRemove
		// containerStop will write HostConfig to disk, we shall restore AutoRemove
		// in disk too
		if toDiskErr := container.ToDiskLocking(); toDiskErr != nil {
			logrus.Errorf("Write container to disk error: %v", toDiskErr)
		}

		if err != nil {
			return err
		}
	}

	if err := daemon.containerStart(container, "", "", true); err != nil {
		return err
	}

	daemon.LogContainerEvent(container, "restart")
	return nil
}
