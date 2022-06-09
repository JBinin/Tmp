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
// +build !windows

package utils

import (
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

func CloseExecFrom(minFd int) error {
	fdList, err := ioutil.ReadDir("/proc/self/fd")
	if err != nil {
		return err
	}
	for _, fi := range fdList {
		fd, err := strconv.Atoi(fi.Name())
		if err != nil {
			// ignore non-numeric file names
			continue
		}

		if fd < minFd {
			// ignore descriptors lower than our specified minimum
			continue
		}

		// intentionally ignore errors from unix.CloseOnExec
		unix.CloseOnExec(fd)
		// the cases where this might fail are basically file descriptors that have already been closed (including and especially the one that was created when ioutil.ReadDir did the "opendir" syscall)
	}
	return nil
}

// NewSockPair returns a new unix socket pair
func NewSockPair(name string) (parent *os.File, child *os.File, err error) {
	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, nil, err
	}
	return os.NewFile(uintptr(fds[1]), name+"-p"), os.NewFile(uintptr(fds[0]), name+"-c"), nil
}
