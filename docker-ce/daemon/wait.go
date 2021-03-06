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
	"time"

	"golang.org/x/net/context"
)

// ContainerWait stops processing until the given container is
// stopped. If the container is not found, an error is returned. On a
// successful stop, the exit code of the container is returned. On a
// timeout, an error is returned. If you want to wait forever, supply
// a negative duration for the timeout.
func (daemon *Daemon) ContainerWait(name string, timeout time.Duration) (int, error) {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return -1, err
	}

	return container.WaitStop(timeout)
}

// ContainerWaitWithContext returns a channel where exit code is sent
// when container stops. Channel can be cancelled with a context.
func (daemon *Daemon) ContainerWaitWithContext(ctx context.Context, name string) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	return container.WaitWithContext(ctx)
}
