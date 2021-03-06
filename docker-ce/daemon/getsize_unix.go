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
// +build linux freebsd solaris

package daemon

import (
	"github.com/Sirupsen/logrus"
)

// getSize returns the real size & virtual size of the container.
func (daemon *Daemon) getSize(containerID string) (int64, int64) {
	var (
		sizeRw, sizeRootfs int64
		err                error
	)

	rwlayer, err := daemon.layerStore.GetRWLayer(containerID)
	if err != nil {
		logrus.Errorf("Failed to compute size of container rootfs %v: %v", containerID, err)
		return sizeRw, sizeRootfs
	}
	defer daemon.layerStore.ReleaseRWLayer(rwlayer)

	sizeRw, err = rwlayer.Size()
	if err != nil {
		logrus.Errorf("Driver %s couldn't return diff size of container %s: %s",
			daemon.GraphDriverName(), containerID, err)
		// FIXME: GetSize should return an error. Not changing it now in case
		// there is a side-effect.
		sizeRw = -1
	}

	if parent := rwlayer.Parent(); parent != nil {
		sizeRootfs, err = parent.Size()
		if err != nil {
			sizeRootfs = -1
		} else if sizeRw != -1 {
			sizeRootfs += sizeRw
		}
	}
	return sizeRw, sizeRootfs
}
