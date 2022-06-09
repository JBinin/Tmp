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

package cpumanager

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
)

type nonePolicy struct{}

var _ Policy = &nonePolicy{}

// PolicyNone name of none policy
const PolicyNone policyName = "none"

// NewNonePolicy returns a cupset manager policy that does nothing
func NewNonePolicy() Policy {
	return &nonePolicy{}
}

func (p *nonePolicy) Name() string {
	return string(PolicyNone)
}

func (p *nonePolicy) Start(s state.State) {
	glog.Info("[cpumanager] none policy: Start")
}

func (p *nonePolicy) AddContainer(s state.State, pod *v1.Pod, container *v1.Container, containerID string) error {
	return nil
}

func (p *nonePolicy) RemoveContainer(s state.State, containerID string) error {
	return nil
}
