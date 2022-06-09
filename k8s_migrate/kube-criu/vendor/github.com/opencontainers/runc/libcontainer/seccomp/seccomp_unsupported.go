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
// +build !linux !cgo !seccomp

package seccomp

import (
	"errors"

	"github.com/opencontainers/runc/libcontainer/configs"
)

var ErrSeccompNotEnabled = errors.New("seccomp: config provided but seccomp not supported")

// InitSeccomp does nothing because seccomp is not supported.
func InitSeccomp(config *configs.Seccomp) error {
	if config != nil {
		return ErrSeccompNotEnabled
	}
	return nil
}

// IsEnabled returns false, because it is not supported.
func IsEnabled() bool {
	return false
}
