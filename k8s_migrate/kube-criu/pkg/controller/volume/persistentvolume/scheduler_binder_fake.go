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

package persistentvolume

import (
	"k8s.io/api/core/v1"
)

type FakeVolumeBinderConfig struct {
	AllBound              bool
	FindUnboundSatsified  bool
	FindBoundSatsified    bool
	FindErr               error
	AssumeBindingRequired bool
	AssumeErr             error
	BindErr               error
}

// NewVolumeBinder sets up all the caches needed for the scheduler to make
// topology-aware volume binding decisions.
func NewFakeVolumeBinder(config *FakeVolumeBinderConfig) *FakeVolumeBinder {
	return &FakeVolumeBinder{
		config: config,
	}
}

type FakeVolumeBinder struct {
	config       *FakeVolumeBinderConfig
	AssumeCalled bool
	BindCalled   bool
}

func (b *FakeVolumeBinder) FindPodVolumes(pod *v1.Pod, node *v1.Node) (unboundVolumesSatisfied, boundVolumesSatsified bool, err error) {
	return b.config.FindUnboundSatsified, b.config.FindBoundSatsified, b.config.FindErr
}

func (b *FakeVolumeBinder) AssumePodVolumes(assumedPod *v1.Pod, nodeName string) (bool, bool, error) {
	b.AssumeCalled = true
	return b.config.AllBound, b.config.AssumeBindingRequired, b.config.AssumeErr
}

func (b *FakeVolumeBinder) BindPodVolumes(assumedPod *v1.Pod) error {
	b.BindCalled = true
	return b.config.BindErr
}

func (b *FakeVolumeBinder) GetBindingsCache() PodBindingCache {
	return nil
}
