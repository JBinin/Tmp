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
// +build !linux

/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package procfs

import (
	"fmt"
	"syscall"
)

type ProcFS struct{}

func NewProcFS() ProcFSInterface {
	return &ProcFS{}
}

// GetFullContainerName gets the container name given the root process id of the container.
func (pfs *ProcFS) GetFullContainerName(pid int) (string, error) {
	return "", fmt.Errorf("GetFullContainerName is unsupported in this build")
}

// Find process(es) using a regular expression and send a specified
// signal to each process
func PKill(name string, sig syscall.Signal) error {
	return fmt.Errorf("PKill is unsupported in this build")
}

// Find process(es) with a specified name (exact match)
// and return their pid(s)
func PidOf(name string) ([]int, error) {
	return []int{}, fmt.Errorf("PidOf is unsupported in this build")
}
