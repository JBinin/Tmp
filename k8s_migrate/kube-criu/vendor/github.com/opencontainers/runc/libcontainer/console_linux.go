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
package libcontainer

import (
	"os"

	"golang.org/x/sys/unix"
)

// mount initializes the console inside the rootfs mounting with the specified mount label
// and applying the correct ownership of the console.
func mountConsole(slavePath string) error {
	oldMask := unix.Umask(0000)
	defer unix.Umask(oldMask)
	f, err := os.Create("/dev/console")
	if err != nil && !os.IsExist(err) {
		return err
	}
	if f != nil {
		f.Close()
	}
	return unix.Mount(slavePath, "/dev/console", "bind", unix.MS_BIND, "")
}

// dupStdio opens the slavePath for the console and dups the fds to the current
// processes stdio, fd 0,1,2.
func dupStdio(slavePath string) error {
	fd, err := unix.Open(slavePath, unix.O_RDWR, 0)
	if err != nil {
		return &os.PathError{
			Op:   "open",
			Path: slavePath,
			Err:  err,
		}
	}
	for _, i := range []int{0, 1, 2} {
		if err := unix.Dup3(fd, i, 0); err != nil {
			return err
		}
	}
	return nil
}
