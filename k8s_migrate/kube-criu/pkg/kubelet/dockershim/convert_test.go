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

package dockershim

import (
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"

	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

func TestConvertDockerStatusToRuntimeAPIState(t *testing.T) {
	testCases := []struct {
		input    string
		expected runtimeapi.ContainerState
	}{
		{input: "Up 5 hours", expected: runtimeapi.ContainerState_CONTAINER_RUNNING},
		{input: "Exited (0) 2 hours ago", expected: runtimeapi.ContainerState_CONTAINER_EXITED},
		{input: "Created", expected: runtimeapi.ContainerState_CONTAINER_CREATED},
		{input: "Random string", expected: runtimeapi.ContainerState_CONTAINER_UNKNOWN},
	}

	for _, test := range testCases {
		actual := toRuntimeAPIContainerState(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestConvertToPullableImageID(t *testing.T) {
	testCases := []struct {
		id       string
		image    *dockertypes.ImageInspect
		expected string
	}{
		{
			id: "image-1",
			image: &dockertypes.ImageInspect{
				RepoDigests: []string{"digest-1"},
			},
			expected: DockerPullableImageIDPrefix + "digest-1",
		},
		{
			id: "image-2",
			image: &dockertypes.ImageInspect{
				RepoDigests: []string{},
			},
			expected: DockerImageIDPrefix + "image-2",
		},
	}

	for _, test := range testCases {
		actual := toPullableImageID(test.id, test.image)
		assert.Equal(t, test.expected, actual)
	}
}
