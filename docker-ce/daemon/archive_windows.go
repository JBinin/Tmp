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

import "github.com/docker/docker/container"

// checkIfPathIsInAVolume checks if the path is in a volume. If it is, it
// cannot be in a read-only volume. If it  is not in a volume, the container
// cannot be configured with a read-only rootfs.
//
// This is a no-op on Windows which does not support read-only volumes, or
// extracting to a mount point inside a volume. TODO Windows: FIXME Post-TP5
func checkIfPathIsInAVolume(container *container.Container, absPath string) (bool, error) {
	return false, nil
}

func fixPermissions(source, destination string, uid, gid int, destExisted bool) error {
	// chown is not supported on Windows
	return nil
}
