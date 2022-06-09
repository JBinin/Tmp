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
// +build solaris freebsd

package container

import (
	"golang.org/x/sys/unix"
)

func detachMounted(path string) error {
	//Solaris and FreeBSD do not support the lazy unmount or MNT_DETACH feature.
	// Therefore there are separate definitions for this.
	return unix.Unmount(path, 0)
}

// SecretMount returns the mount for the secret path
func (container *Container) SecretMount() *Mount {
	return nil
}

// UnmountSecrets unmounts the fs for secrets
func (container *Container) UnmountSecrets() error {
	return nil
}
