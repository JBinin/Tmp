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
// +build !linux,!windows,!freebsd

package osl

// GC triggers garbage collection of namespace path right away
// and waits for it.
func GC() {
}

func GetSandboxForExternalKey(path string, key string) (Sandbox, error) {
	return nil, nil
}

// SetBasePath sets the base url prefix for the ns path
func SetBasePath(path string) {
}
