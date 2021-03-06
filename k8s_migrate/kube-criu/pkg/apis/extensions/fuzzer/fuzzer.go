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

package fuzzer

import (
	"fmt"

	fuzz "github.com/google/gofuzz"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

// Funcs returns the fuzzer functions for the extensions api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(j *extensions.Deployment, c fuzz.Continue) {
			c.FuzzNoCustom(j)

			// match defaulting
			if j.Spec.Selector == nil {
				j.Spec.Selector = &metav1.LabelSelector{MatchLabels: j.Spec.Template.Labels}
			}
			if len(j.Labels) == 0 {
				j.Labels = j.Spec.Template.Labels
			}
		},
		func(j *extensions.DeploymentSpec, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again
			rhl := int32(c.Rand.Int31())
			pds := int32(c.Rand.Int31())
			j.RevisionHistoryLimit = &rhl
			j.ProgressDeadlineSeconds = &pds
		},
		func(j *extensions.DeploymentStrategy, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again
			// Ensure that strategyType is one of valid values.
			strategyTypes := []extensions.DeploymentStrategyType{extensions.RecreateDeploymentStrategyType, extensions.RollingUpdateDeploymentStrategyType}
			j.Type = strategyTypes[c.Rand.Intn(len(strategyTypes))]
			if j.Type != extensions.RollingUpdateDeploymentStrategyType {
				j.RollingUpdate = nil
			} else {
				rollingUpdate := extensions.RollingUpdateDeployment{}
				if c.RandBool() {
					rollingUpdate.MaxUnavailable = intstr.FromInt(int(c.Rand.Int31()))
					rollingUpdate.MaxSurge = intstr.FromInt(int(c.Rand.Int31()))
				} else {
					rollingUpdate.MaxSurge = intstr.FromString(fmt.Sprintf("%d%%", c.Rand.Int31()))
				}
				j.RollingUpdate = &rollingUpdate
			}
		},
		func(j *extensions.DaemonSet, c fuzz.Continue) {
			c.FuzzNoCustom(j)

			// match defaulter
			j.Spec.Template.Generation = 0
			if len(j.ObjectMeta.Labels) == 0 {
				j.ObjectMeta.Labels = j.Spec.Template.ObjectMeta.Labels
			}
		},
		func(j *extensions.DaemonSetSpec, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again
			rhl := int32(c.Rand.Int31())
			j.RevisionHistoryLimit = &rhl
		},
		func(j *extensions.DaemonSetUpdateStrategy, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again
			// Ensure that strategyType is one of valid values.
			strategyTypes := []extensions.DaemonSetUpdateStrategyType{extensions.RollingUpdateDaemonSetStrategyType, extensions.OnDeleteDaemonSetStrategyType}
			j.Type = strategyTypes[c.Rand.Intn(len(strategyTypes))]
			if j.Type != extensions.RollingUpdateDaemonSetStrategyType {
				j.RollingUpdate = nil
			} else {
				rollingUpdate := extensions.RollingUpdateDaemonSet{}
				if c.RandBool() {
					if c.RandBool() {
						rollingUpdate.MaxUnavailable = intstr.FromInt(1 + int(c.Rand.Int31()))
					} else {
						rollingUpdate.MaxUnavailable = intstr.FromString(fmt.Sprintf("%d%%", 1+c.Rand.Int31()))
					}
				}
				j.RollingUpdate = &rollingUpdate
			}
		},
		func(j *extensions.ReplicaSet, c fuzz.Continue) {
			c.FuzzNoCustom(j)

			// match defaulter
			if j.Spec.Selector == nil {
				j.Spec.Selector = &metav1.LabelSelector{MatchLabels: j.Spec.Template.Labels}
			}
			if len(j.Labels) == 0 {
				j.Labels = j.Spec.Template.Labels
			}
		},
	}
}
