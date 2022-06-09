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
package swarm

// RuntimeType is the type of runtime used for the TaskSpec
type RuntimeType string

// RuntimeURL is the proto type url
type RuntimeURL string

const (
	// RuntimeContainer is the container based runtime
	RuntimeContainer RuntimeType = "container"
	// RuntimePlugin is the plugin based runtime
	RuntimePlugin RuntimeType = "plugin"

	// RuntimeURLContainer is the proto url for the container type
	RuntimeURLContainer RuntimeURL = "types.docker.com/RuntimeContainer"
	// RuntimeURLPlugin is the proto url for the plugin type
	RuntimeURLPlugin RuntimeURL = "types.docker.com/RuntimePlugin"
)
