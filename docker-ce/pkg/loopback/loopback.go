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
// +build linux

package loopback

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func getLoopbackBackingFile(file *os.File) (uint64, uint64, error) {
	loopInfo, err := ioctlLoopGetStatus64(file.Fd())
	if err != nil {
		logrus.Errorf("Error get loopback backing file: %s", err)
		return 0, 0, ErrGetLoopbackBackingFile
	}
	return loopInfo.loDevice, loopInfo.loInode, nil
}

// SetCapacity reloads the size for the loopback device.
func SetCapacity(file *os.File) error {
	if err := ioctlLoopSetCapacity(file.Fd(), 0); err != nil {
		logrus.Errorf("Error loopbackSetCapacity: %s", err)
		return ErrSetCapacity
	}
	return nil
}

// FindLoopDeviceFor returns a loopback device file for the specified file which
// is backing file of a loop back device.
func FindLoopDeviceFor(file *os.File) *os.File {
	stat, err := file.Stat()
	if err != nil {
		return nil
	}
	targetInode := stat.Sys().(*syscall.Stat_t).Ino
	targetDevice := stat.Sys().(*syscall.Stat_t).Dev

	for i := 0; true; i++ {
		path := fmt.Sprintf("/dev/loop%d", i)

		file, err := os.OpenFile(path, os.O_RDWR, 0)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}

			// Ignore all errors until the first not-exist
			// we want to continue looking for the file
			continue
		}

		dev, inode, err := getLoopbackBackingFile(file)
		if err == nil && dev == targetDevice && inode == targetInode {
			return file
		}
		file.Close()
	}

	return nil
}
