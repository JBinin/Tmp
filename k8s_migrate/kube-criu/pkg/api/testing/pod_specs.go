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
Copyright 2015 The Kubernetes Authors.

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

package testing

import (
	"k8s.io/api/core/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// DeepEqualSafePodSpec returns a PodSpec which is ready to be used with apiequality.Semantic.DeepEqual
func DeepEqualSafePodSpec() api.PodSpec {
	grace := int64(30)
	return api.PodSpec{
		RestartPolicy:                 api.RestartPolicyAlways,
		DNSPolicy:                     api.DNSClusterFirst,
		TerminationGracePeriodSeconds: &grace,
		SecurityContext:               &api.PodSecurityContext{},
		SchedulerName:                 api.DefaultSchedulerName,
	}
}

// V1DeepEqualSafePodSpec returns a PodSpec which is ready to be used with apiequality.Semantic.DeepEqual
func V1DeepEqualSafePodSpec() v1.PodSpec {
	grace := int64(30)
	return v1.PodSpec{
		RestartPolicy:                 v1.RestartPolicyAlways,
		DNSPolicy:                     v1.DNSClusterFirst,
		TerminationGracePeriodSeconds: &grace,
		SecurityContext:               &v1.PodSecurityContext{},
	}
}
