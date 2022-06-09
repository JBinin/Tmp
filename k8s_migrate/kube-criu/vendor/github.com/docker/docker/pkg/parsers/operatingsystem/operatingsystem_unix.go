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
// +build freebsd darwin

package operatingsystem

import (
	"errors"
	"os/exec"
)

// GetOperatingSystem gets the name of the current operating system.
func GetOperatingSystem() (string, error) {
	cmd := exec.Command("uname", "-s")
	osName, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(osName), nil
}

// IsContainerized returns true if we are running inside a container.
// No-op on FreeBSD and Darwin, always returns false.
func IsContainerized() (bool, error) {
	// TODO: Implement jail detection for freeBSD
	return false, errors.New("Cannot detect if we are in container")
}
