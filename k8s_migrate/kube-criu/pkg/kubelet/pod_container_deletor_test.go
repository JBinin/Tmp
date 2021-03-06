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

package kubelet

import (
	"reflect"
	"testing"
	"time"

	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

func TestGetContainersToDeleteInPodWithFilter(t *testing.T) {
	pod := kubecontainer.PodStatus{
		ContainerStatuses: []*kubecontainer.ContainerStatus{
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "1"},
				Name:      "foo",
				CreatedAt: time.Now(),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "2"},
				Name:      "bar",
				CreatedAt: time.Now().Add(time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "3"},
				Name:      "bar",
				CreatedAt: time.Now().Add(2 * time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "4"},
				Name:      "bar",
				CreatedAt: time.Now().Add(3 * time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "5"},
				Name:      "bar",
				CreatedAt: time.Now().Add(4 * time.Second),
				State:     kubecontainer.ContainerStateRunning,
			},
		},
	}

	testCases := []struct {
		containersToKeep           int
		expectedContainersToDelete containerStatusbyCreatedList
	}{
		{
			0,
			[]*kubecontainer.ContainerStatus{pod.ContainerStatuses[3], pod.ContainerStatuses[2], pod.ContainerStatuses[1]},
		},
		{
			1,
			[]*kubecontainer.ContainerStatus{pod.ContainerStatuses[2], pod.ContainerStatuses[1]},
		},
		{
			2,
			[]*kubecontainer.ContainerStatus{pod.ContainerStatuses[1]},
		},
	}

	for _, test := range testCases {
		candidates := getContainersToDeleteInPod("4", &pod, test.containersToKeep)
		if !reflect.DeepEqual(candidates, test.expectedContainersToDelete) {
			t.Errorf("expected %v got %v", test.expectedContainersToDelete, candidates)
		}
	}
}

func TestGetContainersToDeleteInPod(t *testing.T) {
	pod := kubecontainer.PodStatus{
		ContainerStatuses: []*kubecontainer.ContainerStatus{
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "1"},
				Name:      "foo",
				CreatedAt: time.Now(),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "2"},
				Name:      "bar",
				CreatedAt: time.Now().Add(time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "3"},
				Name:      "bar",
				CreatedAt: time.Now().Add(2 * time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "4"},
				Name:      "bar",
				CreatedAt: time.Now().Add(3 * time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "5"},
				Name:      "bar",
				CreatedAt: time.Now().Add(4 * time.Second),
				State:     kubecontainer.ContainerStateRunning,
			},
		},
	}

	testCases := []struct {
		containersToKeep           int
		expectedContainersToDelete containerStatusbyCreatedList
	}{
		{
			0,
			[]*kubecontainer.ContainerStatus{pod.ContainerStatuses[3], pod.ContainerStatuses[2], pod.ContainerStatuses[1], pod.ContainerStatuses[0]},
		},
		{
			1,
			[]*kubecontainer.ContainerStatus{pod.ContainerStatuses[2], pod.ContainerStatuses[1], pod.ContainerStatuses[0]},
		},
		{
			2,
			[]*kubecontainer.ContainerStatus{pod.ContainerStatuses[1], pod.ContainerStatuses[0]},
		},
	}

	for _, test := range testCases {
		candidates := getContainersToDeleteInPod("", &pod, test.containersToKeep)
		if !reflect.DeepEqual(candidates, test.expectedContainersToDelete) {
			t.Errorf("expected %v got %v", test.expectedContainersToDelete, candidates)
		}
	}
}

func TestGetContainersToDeleteInPodWithNoMatch(t *testing.T) {
	pod := kubecontainer.PodStatus{
		ContainerStatuses: []*kubecontainer.ContainerStatus{
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "1"},
				Name:      "foo",
				CreatedAt: time.Now(),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "2"},
				Name:      "bar",
				CreatedAt: time.Now().Add(time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "3"},
				Name:      "bar",
				CreatedAt: time.Now().Add(2 * time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "4"},
				Name:      "bar",
				CreatedAt: time.Now().Add(3 * time.Second),
				State:     kubecontainer.ContainerStateExited,
			},
			{
				ID:        kubecontainer.ContainerID{Type: "test", ID: "5"},
				Name:      "bar",
				CreatedAt: time.Now().Add(4 * time.Second),
				State:     kubecontainer.ContainerStateRunning,
			},
		},
	}

	testCases := []struct {
		filterID                   string
		expectedContainersToDelete containerStatusbyCreatedList
	}{
		{
			"abc",
			[]*kubecontainer.ContainerStatus{},
		},
	}

	for _, test := range testCases {
		candidates := getContainersToDeleteInPod(test.filterID, &pod, len(pod.ContainerStatuses))
		if !reflect.DeepEqual(candidates, test.expectedContainersToDelete) {
			t.Errorf("expected %v got %v", test.expectedContainersToDelete, candidates)
		}
	}
}
