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

package scheduler

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"k8s.io/kubernetes/pkg/scheduler/core"
	"k8s.io/kubernetes/pkg/scheduler/util"
)

// FakeConfigurator is an implementation for test.
type FakeConfigurator struct {
	Config *Config
}

// GetHardPodAffinitySymmetricWeight is not implemented yet.
func (fc *FakeConfigurator) GetHardPodAffinitySymmetricWeight() int32 {
	panic("not implemented")
}

// MakeDefaultErrorFunc is not implemented yet.
func (fc *FakeConfigurator) MakeDefaultErrorFunc(backoff *util.PodBackoff, podQueue core.SchedulingQueue) func(pod *v1.Pod, err error) {
	return nil
}

// GetNodeLister is not implemented yet.
func (fc *FakeConfigurator) GetNodeLister() corelisters.NodeLister {
	return nil
}

// GetClient is not implemented yet.
func (fc *FakeConfigurator) GetClient() clientset.Interface {
	return nil
}

// GetScheduledPodLister is not implemented yet.
func (fc *FakeConfigurator) GetScheduledPodLister() corelisters.PodLister {
	return nil
}

// Create returns FakeConfigurator.Config
func (fc *FakeConfigurator) Create() (*Config, error) {
	return fc.Config, nil
}

// CreateFromProvider returns FakeConfigurator.Config
func (fc *FakeConfigurator) CreateFromProvider(providerName string) (*Config, error) {
	return fc.Config, nil
}

// CreateFromConfig returns FakeConfigurator.Config
func (fc *FakeConfigurator) CreateFromConfig(policy schedulerapi.Policy) (*Config, error) {
	return fc.Config, nil
}

// CreateFromKeys returns FakeConfigurator.Config
func (fc *FakeConfigurator) CreateFromKeys(predicateKeys, priorityKeys sets.String, extenders []algorithm.SchedulerExtender) (*Config, error) {
	return fc.Config, nil
}
