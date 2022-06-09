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
package system

import "golang.org/x/sys/unix"

// Lgetxattr retrieves the value of the extended attribute identified by attr
// and associated with the given path in the file system.
// It will returns a nil slice and nil error if the xattr is not set.
func Lgetxattr(path string, attr string) ([]byte, error) {
	dest := make([]byte, 128)
	sz, errno := unix.Lgetxattr(path, attr, dest)
	if errno == unix.ENODATA {
		return nil, nil
	}
	if errno == unix.ERANGE {
		dest = make([]byte, sz)
		sz, errno = unix.Lgetxattr(path, attr, dest)
	}
	if errno != nil {
		return nil, errno
	}

	return dest[:sz], nil
}

// Lsetxattr sets the value of the extended attribute identified by attr
// and associated with the given path in the file system.
func Lsetxattr(path string, attr string, data []byte, flags int) error {
	return unix.Lsetxattr(path, attr, data, flags)
}
