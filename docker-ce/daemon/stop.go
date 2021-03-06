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
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/errors"
	"github.com/docker/docker/container"
)

// ContainerStop looks for the given container and terminates it,
// waiting the given number of seconds before forcefully killing the
// container. If a negative number of seconds is given, ContainerStop
// will wait for a graceful termination. An error is returned if the
// container is not found, is already stopped, or if there is a
// problem stopping the container.
func (daemon *Daemon) ContainerStop(name string, seconds *int) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}
	if !container.IsRunning() {
		err := fmt.Errorf("Container %s is already stopped", name)
		return errors.NewErrorWithStatusCode(err, http.StatusNotModified)
	}
	if seconds == nil {
		stopTimeout := container.StopTimeout()
		seconds = &stopTimeout
	}
	if err := daemon.containerStop(container, *seconds); err != nil {
		return fmt.Errorf("Cannot stop container %s: %v", name, err)
	}
	return nil
}

// containerStop halts a container by sending a stop signal, waiting for the given
// duration in seconds, and then calling SIGKILL and waiting for the
// process to exit. If a negative duration is given, Stop will wait
// for the initial signal forever. If the container is not running Stop returns
// immediately.
func (daemon *Daemon) containerStop(container *container.Container, seconds int) error {
	if !container.IsRunning() {
		return nil
	}

	daemon.stopHealthchecks(container)

	stopSignal := container.StopSignal()
	// 1. Send a stop signal
	if err := daemon.killPossiblyDeadProcess(container, stopSignal); err != nil {
		// While normally we might "return err" here we're not going to
		// because if we can't stop the container by this point then
		// it's probably because it's already stopped. Meaning, between
		// the time of the IsRunning() call above and now it stopped.
		// Also, since the err return will be environment specific we can't
		// look for any particular (common) error that would indicate
		// that the process is already dead vs something else going wrong.
		// So, instead we'll give it up to 2 more seconds to complete and if
		// by that time the container is still running, then the error
		// we got is probably valid and so we force kill it.
		if _, err := container.WaitStop(2 * time.Second); err != nil {
			logrus.Infof("Container failed to stop after sending signal %d to the process, force killing", stopSignal)
			if err := daemon.killPossiblyDeadProcess(container, 9); err != nil {
				return err
			}
		}
	}

	// 2. Wait for the process to exit on its own
	if _, err := container.WaitStop(time.Duration(seconds) * time.Second); err != nil {
		logrus.Infof("Container %v failed to exit within %d seconds of signal %d - using the force", container.ID, seconds, stopSignal)
		// 3. If it doesn't, then send SIGKILL
		if err := daemon.Kill(container); err != nil {
			container.WaitStop(-1 * time.Second)
			logrus.Warn(err) // Don't return error because we only care that container is stopped, not what function stopped it
		}
	}

	daemon.LogContainerEvent(container, "stop")
	return nil
}
