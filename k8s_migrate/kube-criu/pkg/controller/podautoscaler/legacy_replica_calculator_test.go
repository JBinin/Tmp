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
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"
	core "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"

	heapster "k8s.io/heapster/metrics/api/v1/types"
	metricsapi "k8s.io/metrics/pkg/apis/metrics/v1alpha1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type legacyReplicaCalcTestCase struct {
	currentReplicas  int32
	expectedReplicas int32
	expectedError    error

	timestamp time.Time

	resource *resourceInfo
	metric   *metricInfo

	podReadiness []v1.ConditionStatus
}

func (tc *legacyReplicaCalcTestCase) prepareTestClient(t *testing.T) *fake.Clientset {

	fakeClient := &fake.Clientset{}
	fakeClient.AddReactor("list", "pods", func(action core.Action) (handled bool, ret runtime.Object, err error) {
		obj := &v1.PodList{}
		for i := 0; i < int(tc.currentReplicas); i++ {
			podReadiness := v1.ConditionTrue
			if tc.podReadiness != nil {
				podReadiness = tc.podReadiness[i]
			}
			podName := fmt.Sprintf("%s-%d", podNamePrefix, i)
			pod := v1.Pod{
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
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

	fakeClient.AddProxyReactor("services", func(action core.Action) (handled bool, ret restclient.ResponseWrapper, err error) {
		var heapsterRawMemResponse []byte

		if tc.resource != nil {
			metrics := metricsapi.PodMetricsList{}
			for i, resValue := range tc.resource.levels {
				podName := fmt.Sprintf("%s-%d", podNamePrefix, i)
				if len(tc.resource.podNames) > i {
					podName = tc.resource.podNames[i]
				}
				podMetric := metricsapi.PodMetrics{
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: testNamespace,
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
			heapsterRawMemResponse, _ = json.Marshal(&metrics)
		} else {
			// only return the pods that we actually asked for
			proxyAction := action.(core.ProxyGetAction)
			pathParts := strings.Split(proxyAction.GetPath(), "/")
			// pathParts should look like [ api, v1, model, namespaces, $NS, pod-list, $PODS, metrics, $METRIC... ]
			if len(pathParts) < 9 {
				return true, nil, fmt.Errorf("invalid heapster path %q", proxyAction.GetPath())
			}

			podNames := strings.Split(pathParts[7], ",")
			podPresent := make([]bool, len(tc.metric.levels))
			for _, name := range podNames {
				if len(name) <= len(podNamePrefix)+1 {
					return true, nil, fmt.Errorf("unknown pod %q", name)
				}
				num, err := strconv.Atoi(name[len(podNamePrefix)+1:])
				if err != nil {
					return true, nil, fmt.Errorf("unknown pod %q", name)
				}
				podPresent[num] = true
			}

			timestamp := tc.timestamp
			metrics := heapster.MetricResultList{}
			for i, level := range tc.metric.levels {
				if !podPresent[i] {
					continue
				}

				floatVal := float64(tc.metric.levels[i]) / 1000.0
				metric := heapster.MetricResult{
					Metrics:         []heapster.MetricPoint{{Timestamp: timestamp, Value: uint64(level), FloatValue: &floatVal}},
					LatestTimestamp: timestamp,
				}
				metrics.Items = append(metrics.Items, metric)
			}
			heapsterRawMemResponse, _ = json.Marshal(&metrics)
		}

		return true, newFakeResponseWrapper(heapsterRawMemResponse), nil
	})

	return fakeClient
}

func (tc *legacyReplicaCalcTestCase) runTest(t *testing.T) {
	testClient := tc.prepareTestClient(t)
	metricsClient := metrics.NewHeapsterMetricsClient(testClient, metrics.DefaultHeapsterNamespace, metrics.DefaultHeapsterScheme, metrics.DefaultHeapsterService, metrics.DefaultHeapsterPort)

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

	} else {
		outReplicas, outUtilization, outTimestamp, err := replicaCalc.GetMetricReplicas(tc.currentReplicas, tc.metric.targetUtilization, tc.metric.name, testNamespace, selector)

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
}

func TestLegacyReplicaCalcDisjointResourcesMetrics(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcScaleUp(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcScaleUpUnreadyLessScale(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcScaleUpUnreadyNoScale(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcScaleUpCM(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{20000, 10000, 30000},
			targetUtilization:   15000,
			expectedUtilization: 20000,
		},
	}
	tc.runTest(t)
}

func TestLegacyReplicaCalcScaleUpCMUnreadyLessScale(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 4,
		podReadiness:     []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionTrue, v1.ConditionFalse},
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{50000, 10000, 30000},
			targetUtilization:   15000,
			expectedUtilization: 30000,
		},
	}
	tc.runTest(t)
}

func TestLegacyReplicaCalcScaleUpCMUnreadyNoScaleWouldScaleDown(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		podReadiness:     []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionTrue, v1.ConditionFalse},
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{50000, 15000, 30000},
			targetUtilization:   15000,
			expectedUtilization: 15000,
		},
	}
	tc.runTest(t)
}

func TestLegacyReplicaCalcScaleDown(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcScaleDownCM(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
		currentReplicas:  5,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{12000, 12000, 12000, 12000, 12000},
			targetUtilization:   20000,
			expectedUtilization: 12000,
		},
	}
	tc.runTest(t)
}

func TestLegacyReplicaCalcScaleDownIgnoresUnreadyPods(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcTolerance(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcToleranceCM(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
		currentReplicas:  3,
		expectedReplicas: 3,
		metric: &metricInfo{
			name:                "qps",
			levels:              []int64{20000, 21000, 21000},
			targetUtilization:   20000,
			expectedUtilization: 20666,
		},
	}
	tc.runTest(t)
}

func TestLegacyReplicaCalcSuperfluousMetrics(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetrics(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcEmptyMetrics(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
		currentReplicas: 4,
		expectedError:   fmt.Errorf("unable to get metrics for resource cpu: no metrics returned from heapster"),
		resource: &resourceInfo{
			name:     v1.ResourceCPU,
			requests: []resource.Quantity{resource.MustParse("1.0"), resource.MustParse("1.0"), resource.MustParse("1.0")},
			levels:   []int64{},

			targetUtilization: 100,
		},
	}
	tc.runTest(t)
}

func TestLegacyReplicaCalcEmptyCPURequest(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetricsNoChangeEq(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetricsNoChangeGt(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetricsNoChangeLt(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetricsUnreadyNoChange(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetricsUnreadyScaleUp(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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

func TestLegacyReplicaCalcMissingMetricsUnreadyScaleDown(t *testing.T) {
	tc := legacyReplicaCalcTestCase{
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
func TestLegacyReplicaCalcComputedToleranceAlgImplementation(t *testing.T) {

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
	tc := legacyReplicaCalcTestCase{
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

// TODO: add more tests
