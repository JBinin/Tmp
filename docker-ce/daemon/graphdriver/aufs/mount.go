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

package aufs

import (
	"os/exec"
	"syscall"

	"github.com/Sirupsen/logrus"
)

// Unmount the target specified.
func Unmount(target string) error {
	if err := exec.Command("auplink", target, "flush").Run(); err != nil {
		logrus.Warnf("Couldn't run auplink before unmount %s: %s", target, err)
	}
	if err := syscall.Unmount(target, 0); err != nil {
		return err
	}
	return nil
}
