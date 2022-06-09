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
package container

import (
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	loggerCloseTimeout = 10 * time.Second
)

// Reset puts a container into a state where it can be restarted again.
func (container *Container) Reset(lock bool) {
	if lock {
		container.Lock()
		defer container.Unlock()
	}

	if err := container.CloseStreams(); err != nil {
		logrus.Errorf("%s: %s", container.ID, err)
	}

	// Re-create a brand new stdin pipe once the container exited
	if container.Config.OpenStdin {
		container.StreamConfig.NewInputPipes()
	}

	if container.LogDriver != nil {
		if container.LogCopier != nil {
			exit := make(chan struct{})
			go func() {
				container.LogCopier.Wait()
				close(exit)
			}()
			select {
			case <-time.After(loggerCloseTimeout):
				logrus.Warn("Logger didn't exit in time: logs may be truncated")
			case <-exit:
			}
		}
		container.LogDriver.Close()
		container.LogCopier = nil
		container.LogDriver = nil
	}
}
