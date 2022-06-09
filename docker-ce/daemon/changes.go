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
	"errors"
	"runtime"
	"time"

	"github.com/docker/docker/pkg/archive"
)

// ContainerChanges returns a list of container fs changes
func (daemon *Daemon) ContainerChanges(name string) ([]archive.Change, error) {
	start := time.Now()
	container, err := daemon.GetContainer(name)
	if err != nil {
		return nil, err
	}

	if runtime.GOOS == "windows" && container.IsRunning() {
		return nil, errors.New("Windows does not support diff of a running container")
	}

	container.Lock()
	defer container.Unlock()
	c, err := container.RWLayer.Changes()
	if err != nil {
		return nil, err
	}
	containerActions.WithValues("changes").UpdateSince(start)
	return c, nil
}
