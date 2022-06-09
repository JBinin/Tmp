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

package idtools

import (
	"os"

	"github.com/docker/docker/pkg/system"
)

// Platforms such as Windows do not support the UID/GID concept. So make this
// just a wrapper around system.MkdirAll.
func mkdirAs(path string, mode os.FileMode, ownerUID, ownerGID int, mkAll, chownExisting bool) error {
	if err := system.MkdirAll(path, mode); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// CanAccess takes a valid (existing) directory and a uid, gid pair and determines
// if that uid, gid pair has access (execute bit) to the directory
// Windows does not require/support this function, so always return true
func CanAccess(path string, uid, gid int) bool {
	return true
}
