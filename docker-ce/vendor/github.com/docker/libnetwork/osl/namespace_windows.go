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
package osl

import "testing"

// GenerateKey generates a sandbox key based on the passed
// container id.
func GenerateKey(containerID string) string {
	maxLen := 12
	if len(containerID) < maxLen {
		maxLen = len(containerID)
	}

	return containerID[:maxLen]
}

// NewSandbox provides a new sandbox instance created in an os specific way
// provided a key which uniquely identifies the sandbox
func NewSandbox(key string, osCreate, isRestore bool) (Sandbox, error) {
	return nil, nil
}

func GetSandboxForExternalKey(path string, key string) (Sandbox, error) {
	return nil, nil
}

// GC triggers garbage collection of namespace path right away
// and waits for it.
func GC() {
}

// InitOSContext initializes OS context while configuring network resources
func InitOSContext() func() {
	return func() {}
}

// SetupTestOSContext sets up a separate test  OS context in which tests will be executed.
func SetupTestOSContext(t *testing.T) func() {
	return func() {}
}

// SetBasePath sets the base url prefix for the ns path
func SetBasePath(path string) {
}
