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
// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Types used for docker containers.
package v1

type DockerStatus struct {
	Version       string            `json:"version"`
	APIVersion    string            `json:"api_version"`
	KernelVersion string            `json:"kernel_version"`
	OS            string            `json:"os"`
	Hostname      string            `json:"hostname"`
	RootDir       string            `json:"root_dir"`
	Driver        string            `json:"driver"`
	DriverStatus  map[string]string `json:"driver_status"`
	ExecDriver    string            `json:"exec_driver"`
	NumImages     int               `json:"num_images"`
	NumContainers int               `json:"num_containers"`
}

type DockerImage struct {
	ID          string   `json:"id"`
	RepoTags    []string `json:"repo_tags"` // repository name and tags.
	Created     int64    `json:"created"`   // unix time since creation.
	VirtualSize int64    `json:"virtual_size"`
	Size        int64    `json:"size"`
}
