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
package graphdriver

import "syscall"

var (
	// Slice of drivers that should be used in an order
	priority = []string{
		"zfs",
	}
)

// Mounted checks if the given path is mounted as the fs type
func Mounted(fsType FsMagic, mountPath string) (bool, error) {
	var buf syscall.Statfs_t
	if err := syscall.Statfs(mountPath, &buf); err != nil {
		return false, err
	}
	return FsMagic(buf.Type) == fsType, nil
}
