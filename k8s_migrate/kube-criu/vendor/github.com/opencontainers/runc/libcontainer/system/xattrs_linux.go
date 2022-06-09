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

// Returns a []byte slice if the xattr is set and nil otherwise
// Requires path and its attribute as arguments
func Lgetxattr(path string, attr string) ([]byte, error) {
	var sz int
	// Start with a 128 length byte array
	dest := make([]byte, 128)
	sz, errno := unix.Lgetxattr(path, attr, dest)

	switch {
	case errno == unix.ENODATA:
		return nil, errno
	case errno == unix.ENOTSUP:
		return nil, errno
	case errno == unix.ERANGE:
		// 128 byte array might just not be good enough,
		// A dummy buffer is used to get the real size
		// of the xattrs on disk
		sz, errno = unix.Lgetxattr(path, attr, []byte{})
		if errno != nil {
			return nil, errno
		}
		dest = make([]byte, sz)
		sz, errno = unix.Lgetxattr(path, attr, dest)
		if errno != nil {
			return nil, errno
		}
	case errno != nil:
		return nil, errno
	}
	return dest[:sz], nil
}
