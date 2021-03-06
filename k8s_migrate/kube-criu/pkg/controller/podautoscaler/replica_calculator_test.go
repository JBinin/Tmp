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

package podautoscaler

import (
	"fmt"
	"math"
	"testing"
	"time"

	autoscalingv2 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
	cmapi "k8s.io/metrics/pkg/apis/custom_metrics/v1beta1"
	emapi "k8s.io/metrics/pkg/apis/external_metrics/v1beta1"
	metricsapi "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
	cmfake "k8s.io/metrics/pkg/client/custom_metrics/fake"
	emfake "k8s.io/metrics/pkg/client/external_metrics/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type resourceInfo struct {
	name     v1.ResourceName
	requests []resource.Quantity
	levels   []int64
	// only applies to pod names returned from "heapster"
	podNames []string

	targetUtilization   int32
	expectedUtilization int32
	expectedValue       int64
}

type metricType int

const (
	objectMetric metricType = iota
	externalMetric
	externalPerPodMetric
	podMetric
)

type metricInfo struct {
	name         string
	levels       []int64
	singleObject *autoscalingv2.CrossVersionObjectReference
	selector     *metav1.LabelSelector
	metricType   metricType

	targetUtilization       int64
	perPodTargetUtilization int64
	expectedUtilization     int64
}

type replicaCalcTestCase struct {
	currentReplicas  int32
	expectedReplicas int32
	expectedError    error

	timestamp time.Time

	resource *resourceInfo
	metric   *metricInfo

	podReadiness []v1.ConditionStatus
	podPhase     []v1.PodPhase
}

const (
	testNamespace       = "test-namespace"
	podNamePrefix       = "test-pod"
	numContainersPerPod = 2
)

func (tc *replicaCalcTestCase) prepareTestClientSet() *fake.Clientset {
	fakeClient := &fake.Clientset{}
	fakeClient.AddReactor("list", "pods", func(action core.Action) (handled bool, ret runtime.Object, err error) {
		obj := &v1.PodList{}
		podsCount := int(tc.currentReplicas)
		// Failed pods are not included in tc.currentReplicas
		if tc.podPhase != nil && len(tc.podPhase) > podsCount {
			podsCount = len(tc.podPhase)
		}
		for i := 0; i < podsCount; i++ {
			podReadiness := v1.ConditionTrue
			if tc.podReadiness != nil && i < len(tc.podReadiness) {
				podReadiness = tc.podReadiness[i]
			}
			podPhase := v1.PodRunning
			if tc.podPhase != nil {
				podPhase = tc.podPhase[i]
			}
			podName := fmt.Sprintf("%s-%d", podNamePrefix, i)
			pod := v1.Pod{
				Status: v1.PodStatus{
					Phase: podPhase,
					Conditions: []v1.PodCondition{
						{
							Type:   v1.PodReady,
							Status: podReadiness,
						},
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: testNamespace,
					Labels: map[string]string{
						"name": podNamePrefix,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{}, {}},
				},
			}

			if tc.resource != nil && i < len(tc.resource.requests) {
				pod.Spec.Containers[0].Resources = v1.ResourceRequirements{
					Requests: v1.ResourceList{
						tc.resource.name: tc.resource.requests[i],
					},
				}
				pod.Spec.Containers[1].Resources = v1.ResourceRequirements{
					Requests: v1.ResourceList{
						tc.resource.name: tc.resource.requests[i],
					},
				}
			}
			obj.Items = append(obj.Items, pod)
		}
		return true, obj, nil
	})
	return fakeClient
}

func (tc *replicaCalcTestCase) prepareTestMetricsClient() *metricsfake.Clientset {
	fakeMetricsClient := &metricsfake.Clientset{}
	// NB: we have to sound like Gollum due to gengo's inability to handle already-plural resource names
	fakeMetricsClient.AddReactor("list", "pods", func(action core.Action) (handled bool, ret runtime.Object, err error) {
		if tc.resource != nil {
			metrics := &metricsapi.PodMetricsList{}
			for i, resValue := range tc.resource.levels {
				podName := fmt.Sprintf("%s-%d", podNamePrefix, i)
				if len(tc.resource.podNames) > i {
					podName = tc.resource.podNames[i]
				}
				// NB: the list reactor actually does label selector filtering for us,
				// so we have to make sure our results match the label selector
				podMetric := metricsapi.PodMetrics{
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: testNamespace,
						Labels:    map[string]string{"name": podNamePrefix},
					},
					Timestamp:  metav1.Time{Time: tc.timestamp},
					Containers: make([]metricsapi.ContainerMetrics, numContainersPerPod),
				}

				for i := 0; i < numContainersPerPod; i++ {
					podMetric.Containers[i] = metricsapi.ContainerMetrics{
						Name: fmt.Sprintf("container%v", i),
						Usage: v1.ResourceList{
							v1.ResourceName(tc.resource.name): *resource.NewMilliQuantity(
								int64(resValue),
								resource.DecimalSI),
						},
					}
				}
				metrics.Items = append(metrics.Items, podMetric)
			}
			return true, metrics, nil
		}

		return true, nil, fmt.Errorf("no pod resource metrics specified in test client")
	})
	return fakeMetricsClient
}

func (tc *replicaCalcTestCase) prepareTestCMClient(t *testing.T) *cmfake.FakeCustomMetricsClient {
	fakeCMClient := &cmfake.FakeCustomMetricsClient{}
	fakeCMClient.AddReactor("get", "*", func(action core.Action) (handled bool, ret runtime.Object, err error) {
		getForAction, wasGetFor := action.(cmfake.GetForAction)
		if !wasGetFor {
			return true, nil, fmt.Errorf("expected a get-for action, got %v instead", action)
		}

		if tc.metric == nil {
			return true, nil, fmt.Errorf("no custom metrics specified in test client")
		}

		assert.Equal(t, tc.metric.name, getForAction.GetMetricName(), "the metric requested should have matched the one specified")

		if getForAction.GetName() == "*" {
			metrics := cmapi.MetricValueList{}

			// multiple objects
			assert.Equal(t, "pods", getForAction.GetResource().Resource, "the type of object that we requested multiple metrics for should have been pods")

			for i, level := range tc.metric.levels {
				podMetric := cmapi.MetricValue{
					DescribedObject: v1.ObjectReference{
						Kind:      "Pod",
						Name:      fmt.Sprintf("%s-%d", podNamePrefix, i),
						Namespace: testNamespace,
					},
					Timestamp:  metav1.Time{Time: tc.timestamp},
					MetricName: tc.metric.name,
					Value:      *resource.NewMilliQuantity(level, resource.DecimalSI),
				}
				metrics.Items = append(metrics.Items, podMetric)
			}

			return true, &metrics, nil
		}
		name := getForAction.GetName()
		mapper := testrestmapper.TestOnlyStaticRESTMapper(legacyscheme.Scheme)
		metrics := &cmapi.MetricValueList{}
		assert.NotNil(t, tc.metric.singleObject, "should have only requested a single-object metric when calling GetObjectMetricReplicas")
		gk := schema.FromAPIVersionAndKind(tc.metric.singleObject.APIVersion, tc.metric.singleObject.Kind).GroupKind()
		mapping, err := mapper.RESTMapping(gk)
		if err != nil {
			return true, nil, fmt.Errorf("unable to get mapping for %s: %v", gk.String(), err)
		}
		groupResource := mapping.Resource.GroupResource()

		assert.Equal(t, groupResource.String(), getForAction.GetResource().Resource, "should have requested metrics for the resource matching the GroupKind passed in")
		assert.Equal(t, tc.metric.singleObject.Name, name, "should have requested metrics for the object matching the name passed in")

		metrics.Items = []cmapi.MetricValue{
			{
				DescribedObject: v1.ObjectReference{
					Kind:       tc.metric.singleObject.Kind,
					APIVersion: tc.metric.singleObject.APIVersion,
					Name:       name,
				},
				Timestamp:  metav1.Time{Time: tc.timestamp},
				MetricName: tc.metric.name,
				Value:      *resource.NewMilliQuantity(int64(tc.metric.levels[0]), resource.DecimalSI),
			},
		}

		return true, metrics, nil
	})
	return fakeCMClient
}

func (tc *replicaCalcTestCase) prepareTestEMClient(t *testing.T) *emfake.FakeExternalMetricsClient {
	fakeEMClient := &emfake.FakeExternalMetricsClient{}
	fakeEMClient.AddReactor("list", "*", func(action core.Action) (handled bool, ret runtime.Object, err error) {
		listAction, wasList := action.(core.ListAction)
		if !wasList {
			return true, nil, fmt.Errorf("expected a list-for action, got %v instead", action)
		}

		if tc.metric == nil {
			return true, nil, fmt.Errorf("no external metrics specified in test client")
		}

		assert.Equal(t, tc.metric.name, listAction.GetResource().Resource, "the metric requested should have matched the one specified")

		selector, err := metav1.LabelSelectorAsSelector(tc.metric.selector)
		if err != nil {
			return true, nil, fmt.Errorf("failed to convert label selector specified in test client")
		}
		assert.Equal(t, selector, listAction.GetListRestrictions().Labels, "the metric selector should have matched the one specified")

		metrics := emapi.ExternalMetricValueList{}

		for _, level := range tc.metric.levels {
			metric := emapi.ExternalMetricValue{
				Timestamp:  metav1.Time{Time: tc.timestamp},
				MetricName: tc.metric.name,
				Value:      *resource.NewMilliQuantity(level, resource.DecimalSI),
			}
			metrics.Items = append(metrics.Items, metric)
		}

		return true, &metrics, nil
	})
	return fakeEMClient
}

func (tc *replicaCalcTestCase) prepareTestClient(t *testing.T) (*fake.Clientset, *metricsfake.Clientset, *cmfake.FakeCustomMetricsClient, *emfake.FakeExternalMetricsClient) {
	fakeClient := tc.prepareTestClientSet()
	fakeMetricsClient := tc.prepareTestMetricsClient()
	fakeCMClient := tc.prepareTestCMClient(t)
	fakeEMClient := tc.prepareTestEMClient(t)
	return fakeClient, fakeMetricsClient, fakeCMClient, fakeEMClient
}

func (tc *replicaCalcTestCase) runTest(t *testing.T) {
	testClient, testMetricsClient, testCMClient, testEMClient := tc.prepareTestClient(t)
	metricsClient := metrics.NewRESTMetricsClient(testMetricsClient.MetricsV1beta1(), testCMClient, testEMClient)

	replicaCalc := &ReplicaCalculator{
		metricsClient: metricsClient,
		podsGetter:    testClient.Core(),
		tolerance:     defaultTestingTolerance,
	}

	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{"name": podNamePrefix},
	})
	if err != nil {
		require.Nil(t, err, "something went horribly wrong...")
	}

	if tc.resource != nil {
		outReplicas, outUtilization, outRawValue, outTimestamp, err := replicaCalc.GetResourceReplicas(tc.currentReplicas, tc.resource.targetUtilization, tc.resource.name, testNamespace, selector)

		if tc.expectedError != nil {
			require.Error(t, err, "there should be an error calculating the replica count")
			assert.Contains(t, err.Error(), tc.expectedError.Error(), "the error message should have contained the expected error message")
			return
		}
		require.NoError(t, err, "there should not have been an error calculating the replica count")
		assert.Equal(t, tc.expectedReplicas, outReplicas, "replicas should be as expected")
		assert.Equal(t, tc.resource.expectedUtilization, outUtilization, "utilization should be as expected")
		assert.Equal(t, tc.resource.expectedValue, outRawValue, "raw value should be as expected")
		assert.True(t, tc.timestamp.Equal(outTimestamp), "timestamp should be as expected")
		return
	}

	var outReplicas int32
	var outUtilization int64
	var outTimestamp time.Time
	switch tc.metric.metricType {
	case objectMetric:
		if tc.metric.singleObject == nil {
			t.Fatal("Metric specified as objectMetric but metric.singleObject is nil.")
		}
		outReplicas, outUtilization, outTimestamp, err = replicaCalc.GetObjectMetricReplicas(tc.currentReplicas, tc.metric.targetUtilization, tc.metric.name, testNamespace, tc.metric.singleObject, selector)
	case externalMetric:
		if tc.metric.selector == nil {
			t.Fatal("Metric specified as externalMetric but metric.selector is nil.")
		}
		if tc.metric.targetUtilization <= 0 {
			t.Fatalf("Metric specified as externalMetric but metric.targetUtilization is %d which is <=0.", tc.metric.targetUtilization)
		}
		outReplicas, outUtilization, outTimestamp, err = replicaCalc.GetExternalMetricReplicas(tc.currentReplicas, tc.metric.targetUtilization, tc.metric.name, testNamespace, tc.metric.selector, selector)
	case externalPerPodMetric:
		if tc.metric.selector == nil {
			t.Fatal("Metric specified as externalPerPodMetric but metric.selector is nil.")
		}
		if tc.metric.perPodTargetUtilization <= 0 {
			t.Fatalf("Metric specified as externalPerPodMetric but metric.perPodTargetUtilization is %d which is <=0.", tc.metric.perPodTargetUtilization)
		}

		outReplicas, outUtilization, outTimestamp, err = replicaCalc.GetExternalPerPodMetricReplicas(tc.currentReplicas, tc.metric.perPodTargetUtilization, tc.metric.name, testNamespace, tc.metric.selector)
	case podMetric:
		outReplicas, outUtilization, outTimestamp, err = replicaCalc.GetMetricReplicas(tc.currentReplicas, tc.metric.targetUtilization, tc.metric.name, testNamespace, selector)
	default:
		t.Fatalf("Unknown metric type: %d", tc.metric.metricType)
	}

	if tc.expectedError != nil {
		require.Error(t, err, "there should be an error calculating the replica count")
		assert.Contains(t, err.Error(), tc.expectedError.Error(), "the error message should have contained the expected error message")
		return
	}
	require.NoError(t, err, "there should not have been an error calculating the replica count")
	assert.Equal(t, tc.expectedReplicas, outReplicas, "replicas should be as expected")
	assert.Equal(t, tc.metric.expectedUtilization, outUtilization, "utilization should be as expected")
	assert.True(t, tc.timestamp.Equal(outTimestamp), "timestamp should be as expected")
}

func TestReplicaCalcDisjointResourcesMetrics(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas: 1,
		expectedError:   fmt.Errorf("no metrics returned matched known pods"),
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0")},
			levels:   []int64{100},
			podNames: []string{"an-older-pod-name"},

			targetUtilization: 100,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUp(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 5,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{300, 500, 700},

			targetUtilization:   30,
			expectedUtilization: 50,
			expectedValue:       numContainersPerPod * 500,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpUnreadyLessScale(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionTrue},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{300, 500, 700},

			targetUtilization:   30,
			expectedUtilization: 60,
			expectedValue:       numContainersPerPod * 600,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpUnreadyNoScale(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		podReadiness:     []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionFalse, v1.ConditionFalse},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{400, 500, 700},

			targetUtilization:   30,
			expectedUtilization: 40,
			expectedValue:       numContainersPerPod * 400,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpIgnoresFailedPods(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  2,
		expectedReplicas: 4,
		podReadiness:     []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionTrue, v1.ConditionFalse, v1.ConditionFalse},
		podPhase:         []v1.PodPhase{v1.PodRunning, v1.PodRunning, v1.PodFailed, v1.PodFailed},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{500, 700},

			targetUtilization:   30,
			expectedUtilization: 60,
			expectedValue:       numContainersPerPod * 600,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCM(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{20000, 10000, 30000},
			targetUtilization:   15000,
			expectedUtilization: 20000,
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMUnreadyLessScale(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		podReadiness:     []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionTrue, v1.ConditionFalse},
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{50000, 10000, 30000},
			targetUtilization:   15000,
			expectedUtilization: 30000,
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMUnreadyNoScaleWouldScaleDown(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionFalse},
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{50000, 15000, 30000},
			targetUtilization:   15000,
			expectedUtilization: 15000,
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMObject(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{20000},
			targetUtilization:   15000,
			expectedUtilization: 20000,
			singleObject: &autoscalingv2.CrossVersionObjectReference{
				Kind:       "Deployment",
				APIVersion: "extensions/v1beta1",
				Name:       "some-deployment",
			},
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMObjectIgnoresUnreadyPods(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 5, // If we did not ignore unready pods, we'd expect 15 replicas.
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionFalse},
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{50000},
			targetUtilization:   10000,
			expectedUtilization: 50000,
			singleObject: &autoscalingv2.CrossVersionObjectReference{
				Kind:       "Deployment",
				APIVersion: "extensions/v1beta1",
				Name:       "some-deployment",
			},
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMExternal(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  1,
		expectedReplicas: 2,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{8600},
			targetUtilization:   4400,
			expectedUtilization: 8600,
			selector:            &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMExternalIgnoresUnreadyPods(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 2, // Would expect 6 if we didn't ignore unready pods
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionFalse},
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{8600},
			targetUtilization:   4400,
			expectedUtilization: 8600,
			selector:            &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:          externalMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpCMExternalNoLabels(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  1,
		expectedReplicas: 2,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{8600},
			targetUtilization:   4400,
			expectedUtilization: 8600,
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleUpPerPodCMExternal(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		metric: &metricInfo{
			name:                    "qps",
			levels:                  []int64{8600},
			perPodTargetUtilization: 2150,
			expectedUtilization:     2867,
			selector:                &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:              externalPerPodMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDown(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{100, 300, 500, 250, 250},

			targetUtilization:   50,
			expectedUtilization: 28,
			expectedValue:       numContainersPerPod * 280,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDownCM(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{12000, 12000, 12000, 12000, 12000},
			targetUtilization:   20000,
			expectedUtilization: 12000,
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDownCMObject(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{12000},
			targetUtilization:   20000,
			expectedUtilization: 12000,
			singleObject: &autoscalingv2.CrossVersionObjectReference{
				Kind:       "Deployment",
				APIVersion: "extensions/v1beta1",
				Name:       "some-deployment",
			},
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDownCMExternal(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{8600},
			targetUtilization:   14334,
			expectedUtilization: 8600,
			selector:            &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:          externalMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDownPerPodCMExternal(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                    "qps",
			levels:                  []int64{8600},
			perPodTargetUtilization: 2867,
			expectedUtilization:     1720,
			selector:                &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:              externalPerPodMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDownIgnoresUnreadyPods(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 2,
		podReadiness:     []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionTrue, v1.ConditionTrue, v1.ConditionFalse, v1.ConditionFalse},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{100, 300, 500, 250, 250},

			targetUtilization:   50,
			expectedUtilization: 30,
			expectedValue:       numContainersPerPod * 300,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcScaleDownIgnoresFailedPods(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		podReadiness:     []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionTrue, v1.ConditionTrue, v1.ConditionTrue, v1.ConditionTrue, v1.ConditionFalse, v1.ConditionFalse},
		podPhase:         []v1.PodPhase{v1.PodRunning, v1.PodRunning, v1.PodRunning, v1.PodRunning, v1.PodRunning, v1.PodFailed, v1.PodFailed},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{100, 300, 500, 250, 250},

			targetUtilization:   50,
			expectedUtilization: 28,
			expectedValue:       numContainersPerPod * 280,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcTolerance(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("0.9"), resource.MustParse("1.0"), resource.MustParse("1.1")},
			levels:   []int64{1010, 1030, 1020},

			targetUtilization:   100,
			expectedUtilization: 102,
			expectedValue:       numContainersPerPod * 1020,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcToleranceCM(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{20000, 21000, 21000},
			targetUtilization:   20000,
			expectedUtilization: 20666,
			metricType:          podMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcToleranceCMObject(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{20666},
			targetUtilization:   20000,
			expectedUtilization: 20666,
			singleObject: &autoscalingv2.CrossVersionObjectReference{
				Kind:       "Deployment",
				APIVersion: "extensions/v1beta1",
				Name:       "some-deployment",
			},
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcToleranceCMExternal(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{8600},
			targetUtilization:   8888,
			expectedUtilization: 8600,
			selector:            &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:          externalMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcTolerancePerPodCMExternal(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                    "qps",
			levels:                  []int64{8600},
			perPodTargetUtilization: 2900,
			expectedUtilization:     2867,
			selector:                &metav1.LabelSelector{MatchLabels: map[string]string{"label": "value"}},
			metricType:              externalPerPodMetric,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcSuperfluousMetrics(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  4,
		expectedReplicas: 24,
		resource: &resourceInfo{
			name:                v1.ResourceCPU,
			requests:            []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:              []int64{4000, 9500, 3000, 7000, 3200, 2000},
			targetUtilization:   100,
			expectedUtilization: 587,
			expectedValue:       numContainersPerPod * 5875,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetrics(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  4,
		expectedReplicas: 3,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{400, 95},

			targetUtilization:   100,
			expectedUtilization: 24,
			expectedValue:       495, // numContainersPerPod * 247, for sufficiently large values of 247
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcEmptyMetrics(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas: 4,
		expectedError:   fmt.Errorf("unable to get metrics for resource cpu: no metrics returned from resource metrics API"),
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{},

			targetUtilization: 100,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcEmptyCPURequest(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas: 1,
		expectedError:   fmt.Errorf("missing request for"),
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{},
			levels:   []int64{200},

			targetUtilization: 100,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetricsNoChangeEq(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  2,
		expectedReplicas: 2,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{1000},

			targetUtilization:   100,
			expectedUtilization: 100,
			expectedValue:       numContainersPerPod * 1000,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetricsNoChangeGt(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  2,
		expectedReplicas: 2,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{1900},

			targetUtilization:   100,
			expectedUtilization: 190,
			expectedValue:       numContainersPerPod * 1900,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetricsNoChangeLt(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  2,
		expectedReplicas: 2,
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{600},

			targetUtilization:   100,
			expectedUtilization: 60,
			expectedValue:       numContainersPerPod * 600,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetricsUnreadyNoChange(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionTrue},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{100, 450},

			targetUtilization:   50,
			expectedUtilization: 45,
			expectedValue:       numContainersPerPod * 450,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetricsUnreadyScaleUp(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionTrue},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{100, 2000},

			targetUtilization:   50,
			expectedUtilization: 200,
			expectedValue:       numContainersPerPod * 2000,
		},
	}
	tc.runTest(t)
}

func TestReplicaCalcMissingMetricsUnreadyScaleDown(t *testing.T) {
	tc := replicaCalcTestCase{
		currentReplicas:  4,
		expectedReplicas: 3,
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionTrue, v1.ConditionTrue},
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{100, 100, 100},

			targetUtilization:   50,
			expectedUtilization: 10,
			expectedValue:       numContainersPerPod * 100,
		},
	}
	tc.runTest(t)
}

// TestComputedToleranceAlgImplementation is a regression test which
// back-calculates a minimal percentage for downscaling based on a small percentage
// increase in pod utilization which is calibrated against the tolerance value.
func TestReplicaCalcComputedToleranceAlgImplementation(t *testing.T) {

	startPods := int32(10)
	// 150 mCPU per pod.
	totalUsedCPUOfAllPods := int64(startPods * 150)
	// Each pod starts out asking for 2X what is really needed.
	// This means we will have a 50% ratio of used/requested
	totalRequestedCPUOfAllPods := int32(2 * totalUsedCPUOfAllPods)
	requestedToUsed := float64(totalRequestedCPUOfAllPods / int32(totalUsedCPUOfAllPods))
	// Spread the amount we ask over 10 pods.  We can add some jitter later in reportedLevels.
	perPodRequested := totalRequestedCPUOfAllPods / startPods

	// Force a minimal scaling event by satisfying  (tolerance < 1 - resourcesUsedRatio).
	target := math.Abs(1/(requestedToUsed*(1-defaultTestingTolerance))) + .01
	finalCPUPercentTarget := int32(target * 100)
	resourcesUsedRatio := float64(totalUsedCPUOfAllPods) / float64(float64(totalRequestedCPUOfAllPods)*target)

	// i.e. .60 * 20 -> scaled down expectation.
	finalPods := int32(math.Ceil(resourcesUsedRatio * float64(startPods)))

	// To breach tolerance we will create a utilization ratio difference of tolerance to usageRatioToleranceValue)
	tc := replicaCalcTestCase{
		currentReplicas:  startPods,
		expectedReplicas: finalPods,
		resource: &resourceInfo{
			name: v1.ResourceCPU,
			levels: []int64{
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
				totalUsedCPUOfAllPods / 10,
			},
			requests: []resource.Quantity{
				resource.MustParse(fmt.Sprint(perPodRequested+100) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested-100) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested+10) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested-10) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested+2) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested-2) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested+1) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested-1) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested) + "m"),
				resource.MustParse(fmt.Sprint(perPodRequested) + "m"),
			},

			targetUtilization:   finalCPUPercentTarget,
			expectedUtilization: int32(totalUsedCPUOfAllPods*100) / totalRequestedCPUOfAllPods,
			expectedValue:       numContainersPerPod * totalUsedCPUOfAllPods / 10,
		},
	}

	tc.runTest(t)

	// Reuse the data structure above, now testing "unscaling".
	// Now, we test that no scaling happens if we are in a very close margin to the tolerance
	target = math.Abs(1/(requestedToUsed*(1-defaultTestingTolerance))) + .004
	finalCPUPercentTarget = int32(target * 100)
	tc.resource.targetUtilization = finalCPUPercentTarget
	tc.currentReplicas = startPods
	tc.expectedReplicas = startPods
	tc.runTest(t)
}

func TestHasPodBeenReadyBefore(t *testing.T) {
	tests := []struct {
		name       string
		conditions []v1.PodCondition
		started    time.Time
		expected   bool
	}{
		{
			"initially unready",
			[]v1.PodCondition{
				{
					Type: v1.PodReady,
					LastTransitionTime: metav1.Time{
						Time: metav1.Date(2018, 7, 25, 17, 10, 0, 0, time.UTC).Time,
					},
					Status: v1.ConditionFalse,
				},
			},
			metav1.Date(2018, 7, 25, 17, 10, 0, 0, time.UTC).Time,
			false,
		},
		{
			"currently unready",
			[]v1.PodCondition{
				{
					Type: v1.PodReady,
					LastTransitionTime: metav1.Time{
						Time: metav1.Date(2018, 7, 25, 17, 10, 0, 0, time.UTC).Time,
					},
					Status: v1.ConditionFalse,
				},
			},
			metav1.Date(2018, 7, 25, 17, 0, 0, 0, time.UTC).Time,
			true,
		},
		{
			"currently ready",
			[]v1.PodCondition{
				{
					Type: v1.PodReady,
					LastTransitionTime: metav1.Time{
						Time: metav1.Date(2018, 7, 25, 17, 10, 0, 0, time.UTC).Time,
					},
					Status: v1.ConditionTrue,
				},
			},
			metav1.Date(2018, 7, 25, 17, 10, 0, 0, time.UTC).Time,
			true,
		},
		{
			"no ready status",
			[]v1.PodCondition{},
			metav1.Date(2018, 7, 25, 17, 10, 0, 0, time.UTC).Time,
			false,
		},
	}
	for _, tc := range tests {
		pod := &v1.Pod{
			Status: v1.PodStatus{
				Conditions: tc.conditions,
				StartTime: &metav1.Time{
					Time: tc.started,
				},
			},
		}
		got := hasPodBeenReadyBefore(pod)
		if got != tc.expected {
			t.Errorf("[TestHasPodBeenReadyBefore.%s] got %v, want %v", tc.name, got, tc.expected)
		}
	}
}

// TODO: add more tests
