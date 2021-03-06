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
// +build cgo,linux

/*
Copyright 2017 The Kubernetes Authors.

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

package cadvisor

import (
	"fmt"

	cadvisorfs "github.com/google/cadvisor/fs"
	"k8s.io/kubernetes/pkg/kubelet/types"
)

// imageFsInfoProvider knows how to translate the configured runtime
// to its file system label for images.
type imageFsInfoProvider struct {
	runtime         string
	runtimeEndpoint string
}

// ImageFsInfoLabel returns the image fs label for the configured runtime.
// For remote runtimes, it handles additional runtimes natively understood by cAdvisor.
func (i *imageFsInfoProvider) ImageFsInfoLabel() (string, error) {
	switch i.runtime {
	case types.DockerContainerRuntime:
		return cadvisorfs.LabelDockerImages, nil
	case types.RemoteContainerRuntime:
		// This is a temporary workaround to get stats for cri-o from cadvisor
		// and should be removed.
		// Related to https://github.com/kubernetes/kubernetes/issues/51798
		if i.runtimeEndpoint == CrioSocket {
			return cadvisorfs.LabelCrioImages, nil
		}
	}
	return "", fmt.Errorf("no imagefs label for configured runtime")
}

// NewImageFsInfoProvider returns a provider for the specified runtime configuration.
func NewImageFsInfoProvider(runtime, runtimeEndpoint string) ImageFsInfoProvider {
	return &imageFsInfoProvider{runtime: runtime, runtimeEndpoint: runtimeEndpoint}
}
