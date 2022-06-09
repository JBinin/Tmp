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
// +build !linux,!windows,!freebsd,!solaris

package osl

import "errors"

var (
	// ErrNotImplemented is for platforms which don't implement sandbox
	ErrNotImplemented = errors.New("not implemented")
)

// NewSandbox provides a new sandbox instance created in an os specific way
// provided a key which uniquely identifies the sandbox
func NewSandbox(key string, osCreate, isRestore bool) (Sandbox, error) {
	return nil, ErrNotImplemented
}

// GenerateKey generates a sandbox key based on the passed
// container id.
func GenerateKey(containerID string) string {
	return ""
}
