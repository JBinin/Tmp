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
	fuzz "github.com/google/gofuzz"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// Funcs returns the fuzzer functions for the autoscaling api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *autoscaling.ScaleStatus, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again

			// ensure we have a valid selector
			metaSelector := &metav1.LabelSelector{}
			c.Fuzz(metaSelector)
			labelSelector, _ := metav1.LabelSelectorAsSelector(metaSelector)
			s.Selector = labelSelector.String()
		},
		func(s *autoscaling.HorizontalPodAutoscalerSpec, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again
			minReplicas := int32(c.Rand.Int31())
			s.MinReplicas = &minReplicas

			randomQuantity := func() resource.Quantity {
				var q resource.Quantity
				c.Fuzz(&q)
				// precalc the string for benchmarking purposes
				_ = q.String()
				return q
			}

			targetUtilization := int32(c.RandUint64())
			s.Metrics = []autoscaling.MetricSpec{
				{
					Type: autoscaling.PodsMetricSourceType,
					Pods: &autoscaling.PodsMetricSource{
						MetricName:         c.RandString(),
						TargetAverageValue: randomQuantity(),
					},
				},
				{
					Type: autoscaling.ResourceMetricSourceType,
					Resource: &autoscaling.ResourceMetricSource{
						Name: api.ResourceCPU,
						TargetAverageUtilization: &targetUtilization,
					},
				},
			}
		},
		func(s *autoscaling.HorizontalPodAutoscalerStatus, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again
			randomQuantity := func() resource.Quantity {
				var q resource.Quantity
				c.Fuzz(&q)
				// precalc the string for benchmarking purposes
				_ = q.String()
				return q
			}
			currentUtilization := int32(c.RandUint64())
			s.CurrentMetrics = []autoscaling.MetricStatus{
				{
					Type: autoscaling.PodsMetricSourceType,
					Pods: &autoscaling.PodsMetricStatus{
						MetricName:          c.RandString(),
						CurrentAverageValue: randomQuantity(),
					},
				},
				{
					Type: autoscaling.ResourceMetricSourceType,
					Resource: &autoscaling.ResourceMetricStatus{
						Name: api.ResourceCPU,
						CurrentAverageUtilization: &currentUtilization,
					},
				},
			}
		},
	}
}
