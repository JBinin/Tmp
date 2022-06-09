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
Copyright 2018 The Kubernetes Authors.

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

package types

import (
	"k8s.io/api/core/v1"
	"testing"
)

func TestPodConditionByKubelet(t *testing.T) {
	trueCases := []v1.PodConditionType{
		v1.PodScheduled,
		v1.PodReady,
		v1.PodInitialized,
		v1.PodReasonUnschedulable,
	}

	for _, tc := range trueCases {
		if !PodConditionByKubelet(tc) {
			t.Errorf("Expect %q to be condition owned by kubelet.", tc)
		}
	}

	falseCases := []v1.PodConditionType{
		v1.PodConditionType("abcd"),
	}

	for _, tc := range falseCases {
		if PodConditionByKubelet(tc) {
			t.Errorf("Expect %q NOT to be condition owned by kubelet.", tc)
		}
	}
}
