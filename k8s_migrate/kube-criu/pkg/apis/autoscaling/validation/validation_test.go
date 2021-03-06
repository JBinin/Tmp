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

package validation

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	api "k8s.io/kubernetes/pkg/apis/core"
	utilpointer "k8s.io/utils/pointer"
)

func TestValidateScale(t *testing.T) {
	successCases := []autoscaling.Scale{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "frontend",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.ScaleSpec{
				Replicas: 1,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "frontend",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.ScaleSpec{
				Replicas: 10,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "frontend",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.ScaleSpec{
				Replicas: 0,
			},
		},
	}

	for _, successCase := range successCases {
		if errs := ValidateScale(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := []struct {
		scale autoscaling.Scale
		msg   string
	}{
		{
			scale: autoscaling.Scale{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "frontend",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: autoscaling.ScaleSpec{
					Replicas: -1,
				},
			},
			msg: "must be greater than or equal to 0",
		},
	}

	for _, c := range errorCases {
		if errs := ValidateScale(&c.scale); len(errs) == 0 {
			t.Errorf("expected failure for %s", c.msg)
		} else if !strings.Contains(errs[0].Error(), c.msg) {
			t.Errorf("unexpected error: %v, expected: %s", errs[0], c.msg)
		}
	}
}

func TestValidateHorizontalPodAutoscaler(t *testing.T) {
	successCases := []autoscaling.HorizontalPodAutoscaler{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
				Metrics: []autoscaling.MetricSpec{
					{
						Type: autoscaling.ResourceMetricSourceType,
						Resource: &autoscaling.ResourceMetricSource{
							Name: api.ResourceCPU,
							TargetAverageUtilization: utilpointer.Int32Ptr(70),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
				Metrics: []autoscaling.MetricSpec{
					{
						Type: autoscaling.ResourceMetricSourceType,
						Resource: &autoscaling.ResourceMetricSource{
							Name:               api.ResourceCPU,
							TargetAverageValue: resource.NewMilliQuantity(300, resource.DecimalSI),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
				Metrics: []autoscaling.MetricSpec{
					{
						Type: autoscaling.PodsMetricSourceType,
						Pods: &autoscaling.PodsMetricSource{
							MetricName:         "somemetric",
							TargetAverageValue: *resource.NewMilliQuantity(300, resource.DecimalSI),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
				Metrics: []autoscaling.MetricSpec{
					{
						Type: autoscaling.ObjectMetricSourceType,
						Object: &autoscaling.ObjectMetricSource{
							Target: autoscaling.CrossVersionObjectReference{
								Kind: "ReplicationController",
								Name: "myrc",
							},
							MetricName:  "somemetric",
							TargetValue: *resource.NewMilliQuantity(300, resource.DecimalSI),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
				Metrics: []autoscaling.MetricSpec{
					{
						Type: autoscaling.ExternalMetricSourceType,
						External: &autoscaling.ExternalMetricSource{
							MetricName: "somemetric",
							MetricSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"label": "value",
								},
							},
							TargetValue: resource.NewMilliQuantity(300, resource.DecimalSI),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: autoscaling.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MinReplicas: utilpointer.Int32Ptr(1),
				MaxReplicas: 5,
				Metrics: []autoscaling.MetricSpec{
					{
						Type: autoscaling.ExternalMetricSourceType,
						External: &autoscaling.ExternalMetricSource{
							MetricName: "somemetric",
							MetricSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"label": "value",
								},
							},
							TargetAverageValue: resource.NewMilliQuantity(300, resource.DecimalSI),
						},
					},
				},
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateHorizontalPodAutoscaler(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := []struct {
		horizontalPodAutoscaler autoscaling.HorizontalPodAutoscaler
		msg                     string
	}{
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: utilpointer.Int32Ptr(70),
							},
						},
					},
				},
			},
			msg: "scaleTargetRef.kind: Required",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Kind: "..", Name: "myrc"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: utilpointer.Int32Ptr(70),
							},
						},
					},
				},
			},
			msg: "scaleTargetRef.kind: Invalid",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: utilpointer.Int32Ptr(70),
							},
						},
					},
				},
			},
			msg: "scaleTargetRef.name: Required",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Kind: "ReplicationController", Name: ".."},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: utilpointer.Int32Ptr(70),
							},
						},
					},
				},
			},
			msg: "scaleTargetRef.name: Invalid",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myautoscaler",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{},
					MinReplicas:    utilpointer.Int32Ptr(-1),
					MaxReplicas:    5,
				},
			},
			msg: "must be greater than 0",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myautoscaler",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{},
					MinReplicas:    utilpointer.Int32Ptr(7),
					MaxReplicas:    5,
				},
			},
			msg: "must be greater than or equal to `minReplicas`",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myautoscaler",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: utilpointer.Int32Ptr(70),
								TargetAverageValue:       resource.NewMilliQuantity(300, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "may not set both a target raw value and a target utilization",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								TargetAverageUtilization: utilpointer.Int32Ptr(70),
							},
						},
					},
				},
			},
			msg: "must specify a resource name",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: utilpointer.Int32Ptr(-10),
							},
						},
					},
				},
			},
			msg: "must be greater than 0",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
							},
						},
					},
				},
			},
			msg: "must set either a target raw value or a target utilization",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricSource{
								TargetAverageValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "must specify a metric name",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricSource{
								MetricName: "somemetric",
							},
						},
					},
				},
			},
			msg: "must specify a positive target value",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ObjectMetricSourceType,
							Object: &autoscaling.ObjectMetricSource{
								Target: autoscaling.CrossVersionObjectReference{
									Name: "myrc",
								},
								MetricName:  "somemetric",
								TargetValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "target.kind: Required",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ObjectMetricSourceType,
							Object: &autoscaling.ObjectMetricSource{
								Target: autoscaling.CrossVersionObjectReference{
									Kind: "ReplicationController",
									Name: "myrc",
								},
								TargetValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "must specify a metric name",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								TargetValue: resource.NewMilliQuantity(300, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "must specify a metric name",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricName: "foo/../",
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								TargetValue: resource.NewMilliQuantity(300, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "'/'",
		},

		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricName: "somemetric",
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
							},
						},
					},
				},
			},
			msg: "must set either a target value for metric or a per-pod target",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricName: "somemetric",
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								TargetValue: resource.NewMilliQuantity(-300, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "must be positive",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricName: "somemetric",
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								TargetAverageValue: resource.NewMilliQuantity(-300, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "must be positive",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricName: "somemetric",
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								TargetValue:        resource.NewMilliQuantity(300, resource.DecimalSI),
								TargetAverageValue: resource.NewMilliQuantity(300, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "may not set both a target value for metric and a per-pod target",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{},
					},
				},
			},
			msg: "must specify a metric source type",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.MetricSourceType("InvalidType"),
						},
					},
				},
			},
			msg: "type: Unsupported value",
		},
		{
			horizontalPodAutoscaler: autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name:               api.ResourceCPU,
								TargetAverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
							Pods: &autoscaling.PodsMetricSource{
								MetricName:         "somemetric",
								TargetAverageValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
			},
			msg: "must populate the given metric source only",
		},
	}

	for _, c := range errorCases {
		errs := ValidateHorizontalPodAutoscaler(&c.horizontalPodAutoscaler)
		if len(errs) == 0 {
			t.Errorf("expected failure for %q", c.msg)
		} else if !strings.Contains(errs[0].Error(), c.msg) {
			t.Errorf("unexpected error: %q, expected: %q", errs[0], c.msg)
		}
	}

	sourceTypes := map[autoscaling.MetricSourceType]autoscaling.MetricSpec{
		autoscaling.ResourceMetricSourceType: {
			Resource: &autoscaling.ResourceMetricSource{
				Name:               api.ResourceCPU,
				TargetAverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
			},
		},
		autoscaling.PodsMetricSourceType: {
			Pods: &autoscaling.PodsMetricSource{
				MetricName:         "somemetric",
				TargetAverageValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
			},
		},
		autoscaling.ObjectMetricSourceType: {
			Object: &autoscaling.ObjectMetricSource{
				Target: autoscaling.CrossVersionObjectReference{
					Kind: "ReplicationController",
					Name: "myrc",
				},
				MetricName:  "somemetric",
				TargetValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
			},
		},
	}

	for correctType, spec := range sourceTypes {
		for incorrectType := range sourceTypes {
			if correctType == incorrectType {
				continue
			}

			spec.Type = incorrectType

			errs := ValidateHorizontalPodAutoscaler(&autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "myautoscaler", Namespace: metav1.NamespaceDefault},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{Name: "myrc", Kind: "ReplicationController"},
					MinReplicas:    utilpointer.Int32Ptr(1),
					MaxReplicas:    5, Metrics: []autoscaling.MetricSpec{spec},
				},
			})

			expectedMsg := "must populate information for the given metric source"

			if len(errs) == 0 {
				t.Errorf("expected failure with type of %v and spec for %v", incorrectType, correctType)
			} else if !strings.Contains(errs[0].Error(), expectedMsg) {
				t.Errorf("unexpected error: %q, expected %q", errs[0], expectedMsg)
			}
		}
	}
}
