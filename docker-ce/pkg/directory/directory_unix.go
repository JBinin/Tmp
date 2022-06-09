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

package directory

import (
	"os"
	"path/filepath"
	"syscall"
)

// Size walks a directory tree and returns its total size in bytes.
func Size(dir string) (size int64, err error) {
	data := make(map[uint64]struct{})
	err = filepath.Walk(dir, func(d string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			// if dir does not exist, Size() returns the error.
			// if dir/x disappeared while walking, Size() ignores dir/x.
			if os.IsNotExist(err) && d != dir {
				return nil
			}
			return err
		}

		// Ignore directory sizes
		if fileInfo == nil {
			return nil
		}

		s := fileInfo.Size()
		if fileInfo.IsDir() || s == 0 {
			return nil
		}

		// Check inode to handle hard links correctly
		inode := fileInfo.Sys().(*syscall.Stat_t).Ino
		// inode is not a uint64 on all platforms. Cast it to avoid issues.
		if _, exists := data[uint64(inode)]; exists {
			return nil
		}
		// inode is not a uint64 on all platforms. Cast it to avoid issues.
		data[uint64(inode)] = struct{}{}

		size += s

		return nil
	})
	return
}
