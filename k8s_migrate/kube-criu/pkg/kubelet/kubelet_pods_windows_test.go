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
// +build windows

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

package kubelet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/util/mount"
)

func TestMakeMountsWindows(t *testing.T) {
	container := v1.Container{
		VolumeMounts: []v1.VolumeMount{
			{
				MountPath: "c:/etc/hosts",
				Name:      "disk",
				ReadOnly:  false,
			},
			{
				MountPath: "c:/mnt/path3",
				Name:      "disk",
				ReadOnly:  true,
			},
			{
				MountPath: "c:/mnt/path4",
				Name:      "disk4",
				ReadOnly:  false,
			},
			{
				MountPath: "c:/mnt/path5",
				Name:      "disk5",
				ReadOnly:  false,
			},
		},
	}

	podVolumes := kubecontainer.VolumeMap{
		"disk":  kubecontainer.VolumeInfo{Mounter: &stubVolume{path: "c:/mnt/disk"}},
		"disk4": kubecontainer.VolumeInfo{Mounter: &stubVolume{path: "c:/mnt/host"}},
		"disk5": kubecontainer.VolumeInfo{Mounter: &stubVolume{path: "c:/var/lib/kubelet/podID/volumes/empty/disk5"}},
	}

	pod := v1.Pod{
		Spec: v1.PodSpec{
			HostNetwork: true,
		},
	}

	fm := &mount.FakeMounter{}
	mounts, _, _ := makeMounts(&pod, "/pod", &container, "fakepodname", "", "", podVolumes, fm, nil)

	expectedMounts := []kubecontainer.Mount{
		{
			Name:           "disk",
			ContainerPath:  "c:/etc/hosts",
			HostPath:       "c:/mnt/disk",
			ReadOnly:       false,
			SELinuxRelabel: false,
		},
		{
			Name:           "disk",
			ContainerPath:  "c:/mnt/path3",
			HostPath:       "c:/mnt/disk",
			ReadOnly:       true,
			SELinuxRelabel: false,
		},
		{
			Name:           "disk4",
			ContainerPath:  "c:/mnt/path4",
			HostPath:       "c:/mnt/host",
			ReadOnly:       false,
			SELinuxRelabel: false,
		},
		{
			Name:           "disk5",
			ContainerPath:  "c:/mnt/path5",
			HostPath:       "c:/var/lib/kubelet/podID/volumes/empty/disk5",
			ReadOnly:       false,
			SELinuxRelabel: false,
		},
	}
	assert.Equal(t, expectedMounts, mounts, "mounts of container %+v", container)
}
