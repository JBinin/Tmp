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
package specs

// State holds information about the runtime state of the container.
type State struct {
	// Version is the version of the specification that is supported.
	Version string `json:"version"`
	// ID is the container ID
	ID string `json:"id"`
	// Status is the runtime state of the container.
	Status string `json:"status"`
	// Pid is the process ID for the container process.
	Pid int `json:"pid"`
	// BundlePath is the path to the container's bundle directory.
	BundlePath string `json:"bundlePath"`
	// Annotations are the annotations associated with the container.
	Annotations map[string]string `json:"annotations"`
}
