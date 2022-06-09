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
	"sort"

	"github.com/docker/docker/container"
	"github.com/docker/docker/volume"
)

// setupMounts configures the mount points for a container by appending each
// of the configured mounts on the container to the OCI mount structure
// which will ultimately be passed into the oci runtime during container creation.
// It also ensures each of the mounts are lexographically sorted.

// BUGBUG TODO Windows containerd. This would be much better if it returned
// an array of runtime spec mounts, not container mounts. Then no need to
// do multiple transitions.

func (daemon *Daemon) setupMounts(c *container.Container) ([]container.Mount, error) {
	var mnts []container.Mount
	for _, mount := range c.MountPoints { // type is volume.MountPoint
		if err := daemon.lazyInitializeVolume(c.ID, mount); err != nil {
			return nil, err
		}
		s, err := mount.Setup(c.MountLabel, 0, 0)
		if err != nil {
			return nil, err
		}

		mnts = append(mnts, container.Mount{
			Source:      s,
			Destination: mount.Destination,
			Writable:    mount.RW,
		})
	}

	sort.Sort(mounts(mnts))
	return mnts, nil
}

// setBindModeIfNull is platform specific processing which is a no-op on
// Windows.
func setBindModeIfNull(bind *volume.MountPoint) {
	return
}
