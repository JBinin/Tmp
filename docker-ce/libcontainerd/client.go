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
package libcontainerd

import (
	"fmt"
	"sync"

	"github.com/docker/docker/pkg/locker"
)

// clientCommon contains the platform agnostic fields used in the client structure
type clientCommon struct {
	backend    Backend
	containers map[string]*container
	locker     *locker.Locker
	mapMutex   sync.RWMutex // protects read/write operations from containers map
}

func (clnt *client) lock(containerID string) {
	clnt.locker.Lock(containerID)
}

func (clnt *client) unlock(containerID string) {
	clnt.locker.Unlock(containerID)
}

// must hold a lock for cont.containerID
func (clnt *client) appendContainer(cont *container) {
	clnt.mapMutex.Lock()
	clnt.containers[cont.containerID] = cont
	clnt.mapMutex.Unlock()
}
func (clnt *client) deleteContainer(containerID string) {
	clnt.mapMutex.Lock()
	delete(clnt.containers, containerID)
	clnt.mapMutex.Unlock()
}

func (clnt *client) getContainer(containerID string) (*container, error) {
	clnt.mapMutex.RLock()
	container, ok := clnt.containers[containerID]
	defer clnt.mapMutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("invalid container: %s", containerID) // fixme: typed error
	}
	return container, nil
}
