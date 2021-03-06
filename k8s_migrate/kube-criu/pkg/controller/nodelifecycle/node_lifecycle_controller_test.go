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

package nodelifecycle

import (
	"context"
	"strings"
	"testing"
	"time"

	"k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	extensionsinformers "k8s.io/client-go/informers/extensions/v1beta1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	testcore "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/cloudprovider"
	fakecloud "k8s.io/kubernetes/pkg/cloudprovider/providers/fake"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/controller/nodelifecycle/scheduler"
	"k8s.io/kubernetes/pkg/controller/testutil"
	nodeutil "k8s.io/kubernetes/pkg/controller/util/node"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
	"k8s.io/kubernetes/pkg/util/node"
	taintutils "k8s.io/kubernetes/pkg/util/taints"
)

const (
	testNodeMonitorGracePeriod = 40 * time.Second
	testNodeStartupGracePeriod = 60 * time.Second
	testNodeMonitorPeriod      = 5 * time.Second
	testRateLimiterQPS         = float32(10000)
	testLargeClusterThreshold  = 20
	testUnhealthyThreshold     = float32(0.55)
)

func alwaysReady() bool { return true }

type nodeLifecycleController struct {
	*Controller
	nodeInformer      coreinformers.NodeInformer
	daemonSetInformer extensionsinformers.DaemonSetInformer
}

// doEviction does the fake eviction and returns the status of eviction operation.
func (nc *nodeLifecycleController) doEviction(fakeNodeHandler *testutil.FakeNodeHandler) bool {
	var podEvicted bool
	zones := testutil.GetZones(fakeNodeHandler)
	for _, zone := range zones {
		nc.zonePodEvictor[zone].Try(func(value scheduler.TimedValue) (bool, time.Duration) {
			uid, _ := value.UID.(string)
			nodeutil.DeletePods(fakeNodeHandler, nc.recorder, value.Value, uid, nc.daemonSetStore)
			return true, 0
		})
	}

	for _, action := range fakeNodeHandler.Actions() {
		if action.GetVerb() == "delete" && action.GetResource().Resource == "pods" {
			podEvicted = true
			return podEvicted
		}
	}
	return podEvicted
}

func (nc *nodeLifecycleController) syncNodeStore(fakeNodeHandler *testutil.FakeNodeHandler) error {
	nodes, err := fakeNodeHandler.List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	newElems := make([]interface{}, 0, len(nodes.Items))
	for i := range nodes.Items {
		newElems = append(newElems, &nodes.Items[i])
	}
	return nc.nodeInformer.Informer().GetStore().Replace(newElems, "newRV")
}

func newNodeLifecycleControllerFromClient(
	cloud cloudprovider.Interface,
	kubeClient clientset.Interface,
	podEvictionTimeout time.Duration,
	evictionLimiterQPS float32,
	secondaryEvictionLimiterQPS float32,
	largeClusterThreshold int32,
	unhealthyZoneThreshold float32,
	nodeMonitorGracePeriod time.Duration,
	nodeStartupGracePeriod time.Duration,
	nodeMonitorPeriod time.Duration,
	useTaints bool,
) (*nodeLifecycleController, error) {

	factory := informers.NewSharedInformerFactory(kubeClient, controller.NoResyncPeriodFunc())

	nodeInformer := factory.Core().V1().Nodes()
	daemonSetInformer := factory.Extensions().V1beta1().DaemonSets()

	nc, err := NewNodeLifecycleController(
		factory.Core().V1().Pods(),
		nodeInformer,
		daemonSetInformer,
		cloud,
		kubeClient,
		nodeMonitorPeriod,
		nodeStartupGracePeriod,
		nodeMonitorGracePeriod,
		podEvictionTimeout,
		evictionLimiterQPS,
		secondaryEvictionLimiterQPS,
		largeClusterThreshold,
		unhealthyZoneThreshold,
		useTaints,
		useTaints,
		useTaints,
	)
	if err != nil {
		return nil, err
	}

	nc.podInformerSynced = alwaysReady
	nc.nodeInformerSynced = alwaysReady
	nc.daemonSetInformerSynced = alwaysReady

	return &nodeLifecycleController{nc, nodeInformer, daemonSetInformer}, nil
}

func TestMonitorNodeStatusEvictPods(t *testing.T) {
	fakeNow := metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
	evictionTimeout := 10 * time.Minute
	labels := map[string]string{
		kubeletapis.LabelZoneRegion:        "region1",
		kubeletapis.LabelZoneFailureDomain: "zone1",
	}

	// Because of the logic that prevents NC from evicting anything when all Nodes are NotReady
	// we need second healthy node in tests. Because of how the tests are written we need to update
	// the status of this Node.
	healthyNodeNewStatus := v1.NodeStatus{
		Conditions: []v1.NodeCondition{
			{
				Type:   v1.NodeReady,
				Status: v1.ConditionTrue,
				// Node status has just been updated, and is NotReady for 10min.
				LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 9, 0, 0, time.UTC),
				LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	table := []struct {
		fakeNodeHandler     *testutil.FakeNodeHandler
		daemonSets          []extensions.DaemonSet
		timeToPass          time.Duration
		newNodeStatus       v1.NodeStatus
		secondNodeNewStatus v1.NodeStatus
		expectedEvictPods   bool
		description         string
	}{
		// Node created recently, with no status (happens only at cluster startup).
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: fakeNow,
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			daemonSets:          nil,
			timeToPass:          0,
			newNodeStatus:       v1.NodeStatus{},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   false,
			description:         "Node created recently, with no status.",
		},
		// Node created recently without FailureDomain labels which is added back later, with no status (happens only at cluster startup).
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: fakeNow,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			daemonSets:          nil,
			timeToPass:          0,
			newNodeStatus:       v1.NodeStatus{},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   false,
			description:         "Node created recently without FailureDomain labels which is added back later, with no status (happens only at cluster startup).",
		},
		// Node created long time ago, and kubelet posted NotReady for a short period of time.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionFalse,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			daemonSets: nil,
			timeToPass: evictionTimeout,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionFalse,
						// Node status has just been updated, and is NotReady for 10min.
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 9, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   false,
			description:         "Node created long time ago, and kubelet posted NotReady for a short period of time.",
		},
		// Pod is ds-managed, and kubelet posted NotReady for a long period of time.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionFalse,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(
					&v1.PodList{
						Items: []v1.Pod{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "pod0",
									Namespace: "default",
									Labels:    map[string]string{"daemon": "yes"},
								},
								Spec: v1.PodSpec{
									NodeName: "node0",
								},
							},
						},
					},
				),
			},
			daemonSets: []extensions.DaemonSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ds0",
						Namespace: "default",
					},
					Spec: extensions.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"daemon": "yes"},
						},
					},
				},
			},
			timeToPass: time.Hour,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionFalse,
						// Node status has just been updated, and is NotReady for 1hr.
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 59, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   false,
			description:         "Pod is ds-managed, and kubelet posted NotReady for a long period of time.",
		},
		// Node created long time ago, and kubelet posted NotReady for a long period of time.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionFalse,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			daemonSets: nil,
			timeToPass: time.Hour,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionFalse,
						// Node status has just been updated, and is NotReady for 1hr.
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 59, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   true,
			description:         "Node created long time ago, and kubelet posted NotReady for a long period of time.",
		},
		// Node created long time ago, node controller posted Unknown for a short period of time.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionUnknown,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			daemonSets: nil,
			timeToPass: evictionTimeout - testNodeMonitorGracePeriod,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionUnknown,
						// Node status was updated by nodecontroller 10min ago
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   false,
			description:         "Node created long time ago, node controller posted Unknown for a short period of time.",
		},
		// Node created long time ago, node controller posted Unknown for a long period of time.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionUnknown,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			daemonSets: nil,
			timeToPass: 60 * time.Minute,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionUnknown,
						// Node status was updated by nodecontroller 1hr ago
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedEvictPods:   true,
			description:         "Node created long time ago, node controller posted Unknown for a long period of time.",
		},
	}

	for _, item := range table {
		nodeController, _ := newNodeLifecycleControllerFromClient(
			nil,
			item.fakeNodeHandler,
			evictionTimeout,
			testRateLimiterQPS,
			testRateLimiterQPS,
			testLargeClusterThreshold,
			testUnhealthyThreshold,
			testNodeMonitorGracePeriod,
			testNodeStartupGracePeriod,
			testNodeMonitorPeriod,
			false)
		nodeController.now = func() metav1.Time { return fakeNow }
		nodeController.recorder = testutil.NewFakeRecorder()
		for _, ds := range item.daemonSets {
			nodeController.daemonSetInformer.Informer().GetStore().Add(&ds)
		}
		if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if item.timeToPass > 0 {
			nodeController.now = func() metav1.Time { return metav1.Time{Time: fakeNow.Add(item.timeToPass)} }
			item.fakeNodeHandler.Existing[0].Status = item.newNodeStatus
			item.fakeNodeHandler.Existing[1].Status = item.secondNodeNewStatus
		}
		if len(item.fakeNodeHandler.Existing[0].Labels) == 0 && len(item.fakeNodeHandler.Existing[1].Labels) == 0 {
			item.fakeNodeHandler.Existing[0].Labels = labels
			item.fakeNodeHandler.Existing[1].Labels = labels
		}
		if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		zones := testutil.GetZones(item.fakeNodeHandler)
		for _, zone := range zones {
			if _, ok := nodeController.zonePodEvictor[zone]; ok {
				nodeController.zonePodEvictor[zone].Try(func(value scheduler.TimedValue) (bool, time.Duration) {
					nodeUID, _ := value.UID.(string)
					nodeutil.DeletePods(item.fakeNodeHandler, nodeController.recorder, value.Value, nodeUID, nodeController.daemonSetInformer.Lister())
					return true, 0
				})
			} else {
				t.Fatalf("Zone %v was unitialized!", zone)
			}
		}

		podEvicted := false
		for _, action := range item.fakeNodeHandler.Actions() {
			if action.GetVerb() == "delete" && action.GetResource().Resource == "pods" {
				podEvicted = true
			}
		}

		if item.expectedEvictPods != podEvicted {
			t.Errorf("expected pod eviction: %+v, got %+v for %+v", item.expectedEvictPods,
				podEvicted, item.description)
		}
	}
}

func TestPodStatusChange(t *testing.T) {
	fakeNow := metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
	evictionTimeout := 10 * time.Minute

	// Because of the logic that prevents NC from evicting anything when all Nodes are NotReady
	// we need second healthy node in tests. Because of how the tests are written we need to update
	// the status of this Node.
	healthyNodeNewStatus := v1.NodeStatus{
		Conditions: []v1.NodeCondition{
			{
				Type:   v1.NodeReady,
				Status: v1.ConditionTrue,
				// Node status has just been updated, and is NotReady for 10min.
				LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 9, 0, 0, time.UTC),
				LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	// Node created long time ago, node controller posted Unknown for a long period of time.
	table := []struct {
		fakeNodeHandler     *testutil.FakeNodeHandler
		daemonSets          []extensions.DaemonSet
		timeToPass          time.Duration
		newNodeStatus       v1.NodeStatus
		secondNodeNewStatus v1.NodeStatus
		expectedPodUpdate   bool
		expectedReason      string
		description         string
	}{
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionUnknown,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node1",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
							Labels: map[string]string{
								kubeletapis.LabelZoneRegion:        "region1",
								kubeletapis.LabelZoneFailureDomain: "zone1",
							},
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:               v1.NodeReady,
									Status:             v1.ConditionTrue,
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			timeToPass: 60 * time.Minute,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionUnknown,
						// Node status was updated by nodecontroller 1hr ago
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			secondNodeNewStatus: healthyNodeNewStatus,
			expectedPodUpdate:   true,
			expectedReason:      node.NodeUnreachablePodReason,
			description: "Node created long time ago, node controller posted Unknown for a " +
				"long period of time, the pod status must include reason for termination.",
		},
	}

	for _, item := range table {
		nodeController, _ := newNodeLifecycleControllerFromClient(
			nil,
			item.fakeNodeHandler,
			evictionTimeout,
			testRateLimiterQPS,
			testRateLimiterQPS,
			testLargeClusterThreshold,
			testUnhealthyThreshold,
			testNodeMonitorGracePeriod,
			testNodeStartupGracePeriod,
			testNodeMonitorPeriod,
			false)
		nodeController.now = func() metav1.Time { return fakeNow }
		nodeController.recorder = testutil.NewFakeRecorder()
		if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if item.timeToPass > 0 {
			nodeController.now = func() metav1.Time { return metav1.Time{Time: fakeNow.Add(item.timeToPass)} }
			item.fakeNodeHandler.Existing[0].Status = item.newNodeStatus
			item.fakeNodeHandler.Existing[1].Status = item.secondNodeNewStatus
		}
		if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		zones := testutil.GetZones(item.fakeNodeHandler)
		for _, zone := range zones {
			nodeController.zonePodEvictor[zone].Try(func(value scheduler.TimedValue) (bool, time.Duration) {
				nodeUID, _ := value.UID.(string)
				nodeutil.DeletePods(item.fakeNodeHandler, nodeController.recorder, value.Value, nodeUID, nodeController.daemonSetStore)
				return true, 0
			})
		}

		podReasonUpdate := false
		for _, action := range item.fakeNodeHandler.Actions() {
			if action.GetVerb() == "update" && action.GetResource().Resource == "pods" {
				updateReason := action.(testcore.UpdateActionImpl).GetObject().(*v1.Pod).Status.Reason
				podReasonUpdate = true
				if updateReason != item.expectedReason {
					t.Errorf("expected pod status reason: %+v, got %+v for %+v", item.expectedReason, updateReason, item.description)
				}
			}
		}

		if podReasonUpdate != item.expectedPodUpdate {
			t.Errorf("expected pod update: %+v, got %+v for %+v", podReasonUpdate, item.expectedPodUpdate, item.description)
		}
	}
}

func TestMonitorNodeStatusEvictPodsWithDisruption(t *testing.T) {
	fakeNow := metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
	evictionTimeout := 10 * time.Minute
	timeToPass := 60 * time.Minute

	// Because of the logic that prevents NC from evicting anything when all Nodes are NotReady
	// we need second healthy node in tests. Because of how the tests are written we need to update
	// the status of this Node.
	healthyNodeNewStatus := v1.NodeStatus{
		Conditions: []v1.NodeCondition{
			{
				Type:               v1.NodeReady,
				Status:             v1.ConditionTrue,
				LastHeartbeatTime:  metav1.Date(2015, 1, 1, 13, 0, 0, 0, time.UTC),
				LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	unhealthyNodeNewStatus := v1.NodeStatus{
		Conditions: []v1.NodeCondition{
			{
				Type:   v1.NodeReady,
				Status: v1.ConditionUnknown,
				// Node status was updated by nodecontroller 1hr ago
				LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
				LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	table := []struct {
		nodeList                []*v1.Node
		podList                 []v1.Pod
		updatedNodeStatuses     []v1.NodeStatus
		expectedInitialStates   map[string]ZoneState
		expectedFollowingStates map[string]ZoneState
		expectedEvictPods       bool
		description             string
	}{
		// NetworkDisruption: Node created long time ago, node controller posted Unknown for a long period of time on both Nodes.
		// Only zone is down - eviction shouldn't take place
		{
			nodeList: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node1",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			podList: []v1.Pod{*testutil.NewPod("pod0", "node0")},
			updatedNodeStatuses: []v1.NodeStatus{
				unhealthyNodeNewStatus,
				unhealthyNodeNewStatus,
			},
			expectedInitialStates:   map[string]ZoneState{testutil.CreateZoneID("region1", "zone1"): stateFullDisruption},
			expectedFollowingStates: map[string]ZoneState{testutil.CreateZoneID("region1", "zone1"): stateFullDisruption},
			expectedEvictPods:       false,
			description:             "Network Disruption: Only zone is down - eviction shouldn't take place.",
		},
		// NetworkDisruption: Node created long time ago, node controller posted Unknown for a long period of time on both Nodes.
		// Both zones down - eviction shouldn't take place
		{
			nodeList: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node1",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region2",
							kubeletapis.LabelZoneFailureDomain: "zone2",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},

			podList: []v1.Pod{*testutil.NewPod("pod0", "node0")},
			updatedNodeStatuses: []v1.NodeStatus{
				unhealthyNodeNewStatus,
				unhealthyNodeNewStatus,
			},
			expectedInitialStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
				testutil.CreateZoneID("region2", "zone2"): stateFullDisruption,
			},
			expectedFollowingStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
				testutil.CreateZoneID("region2", "zone2"): stateFullDisruption,
			},
			expectedEvictPods: false,
			description:       "Network Disruption: Both zones down - eviction shouldn't take place.",
		},
		// NetworkDisruption: Node created long time ago, node controller posted Unknown for a long period of time on both Nodes.
		// One zone is down - eviction should take place
		{
			nodeList: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node1",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone2",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionTrue,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			podList: []v1.Pod{*testutil.NewPod("pod0", "node0")},
			updatedNodeStatuses: []v1.NodeStatus{
				unhealthyNodeNewStatus,
				healthyNodeNewStatus,
			},
			expectedInitialStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
				testutil.CreateZoneID("region1", "zone2"): stateNormal,
			},
			expectedFollowingStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
				testutil.CreateZoneID("region1", "zone2"): stateNormal,
			},
			expectedEvictPods: true,
			description:       "Network Disruption: One zone is down - eviction should take place.",
		},
		// NetworkDisruption: Node created long time ago, node controller posted Unknown for a long period
		// of on first Node, eviction should stop even though -master Node is healthy.
		{
			nodeList: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node-master",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionTrue,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			podList: []v1.Pod{*testutil.NewPod("pod0", "node0")},
			updatedNodeStatuses: []v1.NodeStatus{
				unhealthyNodeNewStatus,
				healthyNodeNewStatus,
			},
			expectedInitialStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
			},
			expectedFollowingStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
			},
			expectedEvictPods: false,
			description:       "NetworkDisruption: eviction should stop, only -master Node is healthy",
		},
		// NetworkDisruption: Node created long time ago, node controller posted Unknown for a long period of time on both Nodes.
		// Initially both zones down, one comes back - eviction should take place
		{
			nodeList: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node1",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone2",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},

			podList: []v1.Pod{*testutil.NewPod("pod0", "node0")},
			updatedNodeStatuses: []v1.NodeStatus{
				unhealthyNodeNewStatus,
				healthyNodeNewStatus,
			},
			expectedInitialStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
				testutil.CreateZoneID("region1", "zone2"): stateFullDisruption,
			},
			expectedFollowingStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): stateFullDisruption,
				testutil.CreateZoneID("region1", "zone2"): stateNormal,
			},
			expectedEvictPods: true,
			description:       "Initially both zones down, one comes back - eviction should take place",
		},
		// NetworkDisruption: Node created long time ago, node controller posted Unknown for a long period of time on both Nodes.
		// Zone is partially disrupted - eviction should take place
		{
			nodeList: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node1",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node2",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node3",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionTrue,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node4",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						Labels: map[string]string{
							kubeletapis.LabelZoneRegion:        "region1",
							kubeletapis.LabelZoneFailureDomain: "zone1",
						},
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionTrue,
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},

			podList: []v1.Pod{*testutil.NewPod("pod0", "node0")},
			updatedNodeStatuses: []v1.NodeStatus{
				unhealthyNodeNewStatus,
				unhealthyNodeNewStatus,
				unhealthyNodeNewStatus,
				healthyNodeNewStatus,
				healthyNodeNewStatus,
			},
			expectedInitialStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): statePartialDisruption,
			},
			expectedFollowingStates: map[string]ZoneState{
				testutil.CreateZoneID("region1", "zone1"): statePartialDisruption,
			},
			expectedEvictPods: true,
			description:       "Zone is partially disrupted - eviction should take place.",
		},
	}

	for _, item := range table {
		fakeNodeHandler := &testutil.FakeNodeHandler{
			Existing:  item.nodeList,
			Clientset: fake.NewSimpleClientset(&v1.PodList{Items: item.podList}),
		}
		nodeController, _ := newNodeLifecycleControllerFromClient(
			nil,
			fakeNodeHandler,
			evictionTimeout,
			testRateLimiterQPS,
			testRateLimiterQPS,
			testLargeClusterThreshold,
			testUnhealthyThreshold,
			testNodeMonitorGracePeriod,
			testNodeStartupGracePeriod,
			testNodeMonitorPeriod,
			false)
		nodeController.now = func() metav1.Time { return fakeNow }
		nodeController.enterPartialDisruptionFunc = func(nodeNum int) float32 {
			return testRateLimiterQPS
		}
		nodeController.recorder = testutil.NewFakeRecorder()
		nodeController.enterFullDisruptionFunc = func(nodeNum int) float32 {
			return testRateLimiterQPS
		}
		if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("%v: unexpected error: %v", item.description, err)
		}

		for zone, state := range item.expectedInitialStates {
			if state != nodeController.zoneStates[zone] {
				t.Errorf("%v: Unexpected zone state: %v: %v instead %v", item.description, zone, nodeController.zoneStates[zone], state)
			}
		}

		nodeController.now = func() metav1.Time { return metav1.Time{Time: fakeNow.Add(timeToPass)} }
		for i := range item.updatedNodeStatuses {
			fakeNodeHandler.Existing[i].Status = item.updatedNodeStatuses[i]
		}

		if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("%v: unexpected error: %v", item.description, err)
		}
		for zone, state := range item.expectedFollowingStates {
			if state != nodeController.zoneStates[zone] {
				t.Errorf("%v: Unexpected zone state: %v: %v instead %v", item.description, zone, nodeController.zoneStates[zone], state)
			}
		}
		var podEvicted bool
		start := time.Now()
		// Infinite loop, used for retrying in case ratelimiter fails to reload for Try function.
		// this breaks when we have the status that we need for test case or when we don't see the
		// intended result after 1 minute.
		for {
			podEvicted = nodeController.doEviction(fakeNodeHandler)
			if podEvicted == item.expectedEvictPods || time.Since(start) > 1*time.Minute {
				break
			}
		}
		if item.expectedEvictPods != podEvicted {
			t.Errorf("%v: expected pod eviction: %+v, got %+v", item.description, item.expectedEvictPods, podEvicted)
		}
	}
}

func TestCloudProviderNodeShutdown(t *testing.T) {

	testCases := []struct {
		testName string
		node     *v1.Node
		shutdown bool
	}{
		{
			testName: "node shutdowned add taint",
			shutdown: true,
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionUnknown,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		{
			testName: "node started after shutdown remove taint",
			shutdown: false,
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
					Taints: []v1.Taint{
						{
							Key:    algorithm.TaintNodeShutdown,
							Effect: v1.TaintEffectNoSchedule,
						},
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			fnh := &testutil.FakeNodeHandler{
				Existing:  []*v1.Node{tc.node},
				Clientset: fake.NewSimpleClientset(),
			}
			nodeController, _ := newNodeLifecycleControllerFromClient(
				nil,
				fnh,
				10*time.Minute,
				testRateLimiterQPS,
				testRateLimiterQPS,
				testLargeClusterThreshold,
				testUnhealthyThreshold,
				testNodeMonitorGracePeriod,
				testNodeStartupGracePeriod,
				testNodeMonitorPeriod,
				false)
			nodeController.cloud = &fakecloud.FakeCloud{}
			nodeController.now = func() metav1.Time { return metav1.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC) }
			nodeController.recorder = testutil.NewFakeRecorder()
			nodeController.nodeShutdownInCloudProvider = func(ctx context.Context, node *v1.Node) (bool, error) {
				return tc.shutdown, nil
			}

			if err := nodeController.syncNodeStore(fnh); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err := nodeController.monitorNodeStatus(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(fnh.UpdatedNodes) != 1 {
				t.Errorf("Node was not updated")
			}
			if tc.shutdown {
				if len(fnh.UpdatedNodes[0].Spec.Taints) != 1 {
					t.Errorf("Node Taint was not added")
				}
				if fnh.UpdatedNodes[0].Spec.Taints[0].Key != "node.cloudprovider.kubernetes.io/shutdown" {
					t.Errorf("Node Taint key is not correct")
				}
			} else {
				if len(fnh.UpdatedNodes[0].Spec.Taints) != 0 {
					t.Errorf("Node Taint was not removed after node is back in ready state")
				}
			}
		})
	}

}

// TestCloudProviderNoRateLimit tests that monitorNodes() immediately deletes
// pods and the node when kubelet has not reported, and the cloudprovider says
// the node is gone.
func TestCloudProviderNoRateLimit(t *testing.T) {
	fnh := &testutil.FakeNodeHandler{
		Existing: []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionUnknown,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		Clientset:      fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0"), *testutil.NewPod("pod1", "node0")}}),
		DeleteWaitChan: make(chan struct{}),
	}
	nodeController, _ := newNodeLifecycleControllerFromClient(
		nil,
		fnh,
		10*time.Minute,
		testRateLimiterQPS,
		testRateLimiterQPS,
		testLargeClusterThreshold,
		testUnhealthyThreshold,
		testNodeMonitorGracePeriod,
		testNodeStartupGracePeriod,
		testNodeMonitorPeriod,
		false)
	nodeController.cloud = &fakecloud.FakeCloud{}
	nodeController.now = func() metav1.Time { return metav1.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC) }
	nodeController.recorder = testutil.NewFakeRecorder()
	nodeController.nodeExistsInCloudProvider = func(nodeName types.NodeName) (bool, error) {
		return false, nil
	}
	nodeController.nodeShutdownInCloudProvider = func(ctx context.Context, node *v1.Node) (bool, error) {
		return false, nil
	}
	// monitorNodeStatus should allow this node to be immediately deleted
	if err := nodeController.syncNodeStore(fnh); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := nodeController.monitorNodeStatus(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	select {
	case <-fnh.DeleteWaitChan:
	case <-time.After(wait.ForeverTestTimeout):
		t.Errorf("Timed out waiting %v for node to be deleted", wait.ForeverTestTimeout)
	}
	if len(fnh.DeletedNodes) != 1 || fnh.DeletedNodes[0].Name != "node0" {
		t.Errorf("Node was not deleted")
	}
	if nodeOnQueue := nodeController.zonePodEvictor[""].Remove("node0"); nodeOnQueue {
		t.Errorf("Node was queued for eviction. Should have been immediately deleted.")
	}
}

func TestMonitorNodeStatusUpdateStatus(t *testing.T) {
	fakeNow := metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
	table := []struct {
		fakeNodeHandler      *testutil.FakeNodeHandler
		timeToPass           time.Duration
		newNodeStatus        v1.NodeStatus
		expectedEvictPods    bool
		expectedRequestCount int
		expectedNodes        []*v1.Node
	}{
		// Node created long time ago, without status:
		// Expect Unknown status posted from node controller.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			expectedRequestCount: 2, // List+Update
			expectedNodes: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusNeverUpdated",
								Message:            "Kubelet never posted node status.",
								LastHeartbeatTime:  metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
								LastTransitionTime: fakeNow,
							},
							{
								Type:               v1.NodeOutOfDisk,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusNeverUpdated",
								Message:            "Kubelet never posted node status.",
								LastHeartbeatTime:  metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
								LastTransitionTime: fakeNow,
							},
							{
								Type:               v1.NodeMemoryPressure,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusNeverUpdated",
								Message:            "Kubelet never posted node status.",
								LastHeartbeatTime:  metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
								LastTransitionTime: fakeNow,
							},
							{
								Type:               v1.NodeDiskPressure,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusNeverUpdated",
								Message:            "Kubelet never posted node status.",
								LastHeartbeatTime:  metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
								LastTransitionTime: fakeNow,
							},
						},
					},
				},
			},
		},
		// Node created recently, without status.
		// Expect no action from node controller (within startup grace period).
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: fakeNow,
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			expectedRequestCount: 1, // List
			expectedNodes:        nil,
		},
		// Node created long time ago, with status updated by kubelet exceeds grace period.
		// Expect Unknown status posted from node controller.
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:   v1.NodeReady,
									Status: v1.ConditionTrue,
									// Node status hasn't been updated for 1hr.
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
								{
									Type:   v1.NodeOutOfDisk,
									Status: v1.ConditionFalse,
									// Node status hasn't been updated for 1hr.
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
							Capacity: v1.ResourceList{
								v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
								v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			expectedRequestCount: 3, // (List+)List+Update
			timeToPass:           time.Hour,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionTrue,
						// Node status hasn't been updated for 1hr.
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
					{
						Type:   v1.NodeOutOfDisk,
						Status: v1.ConditionFalse,
						// Node status hasn't been updated for 1hr.
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
				Capacity: v1.ResourceList{
					v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
					v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
				},
			},
			expectedNodes: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "node0",
						CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					Status: v1.NodeStatus{
						Conditions: []v1.NodeCondition{
							{
								Type:               v1.NodeReady,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusUnknown",
								Message:            "Kubelet stopped posting node status.",
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Time{Time: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC).Add(time.Hour)},
							},
							{
								Type:               v1.NodeOutOfDisk,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusUnknown",
								Message:            "Kubelet stopped posting node status.",
								LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								LastTransitionTime: metav1.Time{Time: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC).Add(time.Hour)},
							},
							{
								Type:               v1.NodeMemoryPressure,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusNeverUpdated",
								Message:            "Kubelet never posted node status.",
								LastHeartbeatTime:  metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), // should default to node creation time if condition was never updated
								LastTransitionTime: metav1.Time{Time: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC).Add(time.Hour)},
							},
							{
								Type:               v1.NodeDiskPressure,
								Status:             v1.ConditionUnknown,
								Reason:             "NodeStatusNeverUpdated",
								Message:            "Kubelet never posted node status.",
								LastHeartbeatTime:  metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), // should default to node creation time if condition was never updated
								LastTransitionTime: metav1.Time{Time: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC).Add(time.Hour)},
							},
						},
						Capacity: v1.ResourceList{
							v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
							v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
						},
					},
				},
			},
		},
		// Node created long time ago, with status updated recently.
		// Expect no action from node controller (within monitor grace period).
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:   v1.NodeReady,
									Status: v1.ConditionTrue,
									// Node status has just been updated.
									LastHeartbeatTime:  fakeNow,
									LastTransitionTime: fakeNow,
								},
							},
							Capacity: v1.ResourceList{
								v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
								v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			expectedRequestCount: 1, // List
			expectedNodes:        nil,
		},
	}

	for i, item := range table {
		nodeController, _ := newNodeLifecycleControllerFromClient(
			nil,
			item.fakeNodeHandler,
			5*time.Minute,
			testRateLimiterQPS,
			testRateLimiterQPS,
			testLargeClusterThreshold,
			testUnhealthyThreshold,
			testNodeMonitorGracePeriod,
			testNodeStartupGracePeriod,
			testNodeMonitorPeriod,
			false)
		nodeController.now = func() metav1.Time { return fakeNow }
		nodeController.recorder = testutil.NewFakeRecorder()
		if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if item.timeToPass > 0 {
			nodeController.now = func() metav1.Time { return metav1.Time{Time: fakeNow.Add(item.timeToPass)} }
			item.fakeNodeHandler.Existing[0].Status = item.newNodeStatus
			if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err := nodeController.monitorNodeStatus(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
		if item.expectedRequestCount != item.fakeNodeHandler.RequestCount {
			t.Errorf("expected %v call, but got %v.", item.expectedRequestCount, item.fakeNodeHandler.RequestCount)
		}
		if len(item.fakeNodeHandler.UpdatedNodes) > 0 && !apiequality.Semantic.DeepEqual(item.expectedNodes, item.fakeNodeHandler.UpdatedNodes) {
			t.Errorf("Case[%d] unexpected nodes: %s", i, diff.ObjectDiff(item.expectedNodes[0], item.fakeNodeHandler.UpdatedNodes[0]))
		}
		if len(item.fakeNodeHandler.UpdatedNodeStatuses) > 0 && !apiequality.Semantic.DeepEqual(item.expectedNodes, item.fakeNodeHandler.UpdatedNodeStatuses) {
			t.Errorf("Case[%d] unexpected nodes: %s", i, diff.ObjectDiff(item.expectedNodes[0], item.fakeNodeHandler.UpdatedNodeStatuses[0]))
		}
	}
}

func TestMonitorNodeStatusMarkPodsNotReady(t *testing.T) {
	fakeNow := metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
	table := []struct {
		fakeNodeHandler         *testutil.FakeNodeHandler
		timeToPass              time.Duration
		newNodeStatus           v1.NodeStatus
		expectedPodStatusUpdate bool
	}{
		// Node created recently, without status.
		// Expect no action from node controller (within startup grace period).
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: fakeNow,
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			expectedPodStatusUpdate: false,
		},
		// Node created long time ago, with status updated recently.
		// Expect no action from node controller (within monitor grace period).
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:   v1.NodeReady,
									Status: v1.ConditionTrue,
									// Node status has just been updated.
									LastHeartbeatTime:  fakeNow,
									LastTransitionTime: fakeNow,
								},
							},
							Capacity: v1.ResourceList{
								v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
								v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			expectedPodStatusUpdate: false,
		},
		// Node created long time ago, with status updated by kubelet exceeds grace period.
		// Expect pods status updated and Unknown node status posted from node controller
		{
			fakeNodeHandler: &testutil.FakeNodeHandler{
				Existing: []*v1.Node{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "node0",
							CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						Status: v1.NodeStatus{
							Conditions: []v1.NodeCondition{
								{
									Type:   v1.NodeReady,
									Status: v1.ConditionTrue,
									// Node status hasn't been updated for 1hr.
									LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
									LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
								},
							},
							Capacity: v1.ResourceList{
								v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
								v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
							},
						},
					},
				},
				Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
			},
			timeToPass: 1 * time.Minute,
			newNodeStatus: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					{
						Type:   v1.NodeReady,
						Status: v1.ConditionTrue,
						// Node status hasn't been updated for 1hr.
						LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				},
				Capacity: v1.ResourceList{
					v1.ResourceName(v1.ResourceCPU):    resource.MustParse("10"),
					v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
				},
			},
			expectedPodStatusUpdate: true,
		},
	}

	for i, item := range table {
		nodeController, _ := newNodeLifecycleControllerFromClient(
			nil,
			item.fakeNodeHandler,
			5*time.Minute,
			testRateLimiterQPS,
			testRateLimiterQPS,
			testLargeClusterThreshold,
			testUnhealthyThreshold,
			testNodeMonitorGracePeriod,
			testNodeStartupGracePeriod,
			testNodeMonitorPeriod,
			false)
		nodeController.now = func() metav1.Time { return fakeNow }
		nodeController.recorder = testutil.NewFakeRecorder()
		if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := nodeController.monitorNodeStatus(); err != nil {
			t.Errorf("Case[%d] unexpected error: %v", i, err)
		}
		if item.timeToPass > 0 {
			nodeController.now = func() metav1.Time { return metav1.Time{Time: fakeNow.Add(item.timeToPass)} }
			item.fakeNodeHandler.Existing[0].Status = item.newNodeStatus
			if err := nodeController.syncNodeStore(item.fakeNodeHandler); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err := nodeController.monitorNodeStatus(); err != nil {
				t.Errorf("Case[%d] unexpected error: %v", i, err)
			}
		}

		podStatusUpdated := false
		for _, action := range item.fakeNodeHandler.Actions() {
			if action.GetVerb() == "update" && action.GetResource().Resource == "pods" && action.GetSubresource() == "status" {
				podStatusUpdated = true
			}
		}
		if podStatusUpdated != item.expectedPodStatusUpdate {
			t.Errorf("Case[%d] expect pod status updated to be %v, but got %v", i, item.expectedPodStatusUpdate, podStatusUpdated)
		}
	}
}

func TestSwapUnreachableNotReadyTaints(t *testing.T) {
	fakeNow := metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC)
	evictionTimeout := 10 * time.Minute

	fakeNodeHandler := &testutil.FakeNodeHandler{
		Existing: []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionUnknown,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			// Because of the logic that prevents NC from evicting anything when all Nodes are NotReady
			// we need second healthy node in tests. Because of how the tests are written we need to update
			// the status of this Node.
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node1",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
	}
	timeToPass := evictionTimeout
	newNodeStatus := v1.NodeStatus{
		Conditions: []v1.NodeCondition{
			{
				Type:   v1.NodeReady,
				Status: v1.ConditionFalse,
				// Node status has just been updated, and is NotReady for 10min.
				LastHeartbeatTime:  metav1.Date(2017, 1, 1, 12, 9, 0, 0, time.UTC),
				LastTransitionTime: metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	healthyNodeNewStatus := v1.NodeStatus{
		Conditions: []v1.NodeCondition{
			{
				Type:               v1.NodeReady,
				Status:             v1.ConditionTrue,
				LastHeartbeatTime:  metav1.Date(2017, 1, 1, 12, 10, 0, 0, time.UTC),
				LastTransitionTime: metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	originalTaint := UnreachableTaintTemplate
	updatedTaint := NotReadyTaintTemplate

	nodeController, _ := newNodeLifecycleControllerFromClient(
		nil,
		fakeNodeHandler,
		evictionTimeout,
		testRateLimiterQPS,
		testRateLimiterQPS,
		testLargeClusterThreshold,
		testUnhealthyThreshold,
		testNodeMonitorGracePeriod,
		testNodeStartupGracePeriod,
		testNodeMonitorPeriod,
		true)
	nodeController.now = func() metav1.Time { return fakeNow }
	nodeController.recorder = testutil.NewFakeRecorder()
	if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := nodeController.monitorNodeStatus(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	nodeController.doNoExecuteTaintingPass()

	node0, err := fakeNodeHandler.Get("node0", metav1.GetOptions{})
	if err != nil {
		t.Errorf("Can't get current node0...")
		return
	}
	node1, err := fakeNodeHandler.Get("node1", metav1.GetOptions{})
	if err != nil {
		t.Errorf("Can't get current node1...")
		return
	}

	if originalTaint != nil && !taintutils.TaintExists(node0.Spec.Taints, originalTaint) {
		t.Errorf("Can't find taint %v in %v", originalTaint, node0.Spec.Taints)
	}

	nodeController.now = func() metav1.Time { return metav1.Time{Time: fakeNow.Add(timeToPass)} }

	node0.Status = newNodeStatus
	node1.Status = healthyNodeNewStatus
	_, err = fakeNodeHandler.UpdateStatus(node0)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = fakeNodeHandler.UpdateStatus(node1)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := nodeController.monitorNodeStatus(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	nodeController.doNoExecuteTaintingPass()

	node0, err = fakeNodeHandler.Get("node0", metav1.GetOptions{})
	if err != nil {
		t.Errorf("Can't get current node0...")
		return
	}
	if updatedTaint != nil {
		if !taintutils.TaintExists(node0.Spec.Taints, updatedTaint) {
			t.Errorf("Can't find taint %v in %v", updatedTaint, node0.Spec.Taints)
		}
	}
}

func TestTaintsNodeByCondition(t *testing.T) {
	fakeNow := metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC)
	evictionTimeout := 10 * time.Minute

	fakeNodeHandler := &testutil.FakeNodeHandler{
		Existing: []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
	}

	nodeController, _ := newNodeLifecycleControllerFromClient(
		nil,
		fakeNodeHandler,
		evictionTimeout,
		testRateLimiterQPS,
		testRateLimiterQPS,
		testLargeClusterThreshold,
		testUnhealthyThreshold,
		testNodeMonitorGracePeriod,
		testNodeStartupGracePeriod,
		testNodeMonitorPeriod,
		true)
	nodeController.now = func() metav1.Time { return fakeNow }
	nodeController.recorder = testutil.NewFakeRecorder()

	outOfDiskTaint := &v1.Taint{
		Key:    algorithm.TaintNodeOutOfDisk,
		Effect: v1.TaintEffectNoSchedule,
	}
	networkUnavailableTaint := &v1.Taint{
		Key:    algorithm.TaintNodeNetworkUnavailable,
		Effect: v1.TaintEffectNoSchedule,
	}
	notReadyTaint := &v1.Taint{
		Key:    algorithm.TaintNodeNotReady,
		Effect: v1.TaintEffectNoSchedule,
	}

	tests := []struct {
		Name           string
		Node           *v1.Node
		ExpectedTaints []*v1.Taint
	}{
		{
			Name: "NetworkUnavailable is true",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
						{
							Type:               v1.NodeNetworkUnavailable,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			ExpectedTaints: []*v1.Taint{networkUnavailableTaint},
		},
		{
			Name: "NetworkUnavailable and OutOfDisk are true",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
						{
							Type:               v1.NodeNetworkUnavailable,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
						{
							Type:               v1.NodeOutOfDisk,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			ExpectedTaints: []*v1.Taint{networkUnavailableTaint, outOfDiskTaint},
		},
		{
			Name: "NetworkUnavailable is true, OutOfDisk is unknown",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
						{
							Type:               v1.NodeNetworkUnavailable,
							Status:             v1.ConditionTrue,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
						{
							Type:               v1.NodeOutOfDisk,
							Status:             v1.ConditionUnknown,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			ExpectedTaints: []*v1.Taint{networkUnavailableTaint},
		},
		{
			Name: "Ready is false",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionFalse,
							LastHeartbeatTime:  metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			ExpectedTaints: []*v1.Taint{notReadyTaint},
		},
	}

	for _, test := range tests {
		fakeNodeHandler.Update(test.Node)
		if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		nodeController.doNoScheduleTaintingPass(test.Node)
		if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		node0, err := nodeController.nodeLister.Get("node0")
		if err != nil {
			t.Errorf("Can't get current node0...")
			return
		}
		if len(node0.Spec.Taints) != len(test.ExpectedTaints) {
			t.Errorf("%s: Unexpected number of taints: expected %d, got %d",
				test.Name, len(test.ExpectedTaints), len(node0.Spec.Taints))
		}
		for _, taint := range test.ExpectedTaints {
			if !taintutils.TaintExists(node0.Spec.Taints, taint) {
				t.Errorf("%s: Can't find taint %v in %v", test.Name, taint, node0.Spec.Taints)
			}
		}
	}
}

func TestNodeEventGeneration(t *testing.T) {
	fakeNow := metav1.Date(2016, 9, 10, 12, 0, 0, 0, time.UTC)
	fakeNodeHandler := &testutil.FakeNodeHandler{
		Existing: []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					UID:               "1234567890",
					CreationTimestamp: metav1.Date(2015, 8, 10, 0, 0, 0, 0, time.UTC),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:               v1.NodeReady,
							Status:             v1.ConditionUnknown,
							LastHeartbeatTime:  metav1.Date(2015, 8, 10, 0, 0, 0, 0, time.UTC),
							LastTransitionTime: metav1.Date(2015, 8, 10, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
	}

	nodeController, _ := newNodeLifecycleControllerFromClient(
		nil,
		fakeNodeHandler,
		5*time.Minute,
		testRateLimiterQPS,
		testRateLimiterQPS,
		testLargeClusterThreshold,
		testUnhealthyThreshold,
		testNodeMonitorGracePeriod,
		testNodeStartupGracePeriod,
		testNodeMonitorPeriod,
		false)
	nodeController.cloud = &fakecloud.FakeCloud{}
	nodeController.nodeExistsInCloudProvider = func(nodeName types.NodeName) (bool, error) {
		return false, nil
	}
	nodeController.nodeShutdownInCloudProvider = func(ctx context.Context, node *v1.Node) (bool, error) {
		return false, nil
	}
	nodeController.now = func() metav1.Time { return fakeNow }
	fakeRecorder := testutil.NewFakeRecorder()
	nodeController.recorder = fakeRecorder
	if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := nodeController.monitorNodeStatus(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(fakeRecorder.Events) != 2 {
		t.Fatalf("unexpected events, got %v, expected %v: %+v", len(fakeRecorder.Events), 2, fakeRecorder.Events)
	}
	if fakeRecorder.Events[0].Reason != "RegisteredNode" || fakeRecorder.Events[1].Reason != "DeletingNode" {
		var reasons []string
		for _, event := range fakeRecorder.Events {
			reasons = append(reasons, event.Reason)
		}
		t.Fatalf("unexpected events generation: %v", strings.Join(reasons, ","))
	}
	for _, event := range fakeRecorder.Events {
		involvedObject := event.InvolvedObject
		actualUID := string(involvedObject.UID)
		if actualUID != "1234567890" {
			t.Fatalf("unexpected event uid: %v", actualUID)
		}
	}
}

// TestFixDeprecatedTaintKey verifies we have backwards compatibility after upgraded alpha taint key to GA taint key.
// TODO(resouer): this is introduced in 1.9 and should be removed in the future.
func TestFixDeprecatedTaintKey(t *testing.T) {
	fakeNow := metav1.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC)
	evictionTimeout := 10 * time.Minute

	fakeNodeHandler := &testutil.FakeNodeHandler{
		Existing: []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
			},
		},
		Clientset: fake.NewSimpleClientset(&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod0", "node0")}}),
	}

	nodeController, _ := newNodeLifecycleControllerFromClient(
		nil,
		fakeNodeHandler,
		evictionTimeout,
		testRateLimiterQPS,
		testRateLimiterQPS,
		testLargeClusterThreshold,
		testUnhealthyThreshold,
		testNodeMonitorGracePeriod,
		testNodeStartupGracePeriod,
		testNodeMonitorPeriod,
		true)
	nodeController.now = func() metav1.Time { return fakeNow }
	nodeController.recorder = testutil.NewFakeRecorder()

	deprecatedNotReadyTaint := &v1.Taint{
		Key:    algorithm.DeprecatedTaintNodeNotReady,
		Effect: v1.TaintEffectNoExecute,
	}

	nodeNotReadyTaint := &v1.Taint{
		Key:    algorithm.TaintNodeNotReady,
		Effect: v1.TaintEffectNoExecute,
	}

	deprecatedUnreachableTaint := &v1.Taint{
		Key:    algorithm.DeprecatedTaintNodeUnreachable,
		Effect: v1.TaintEffectNoExecute,
	}

	nodeUnreachableTaint := &v1.Taint{
		Key:    algorithm.TaintNodeUnreachable,
		Effect: v1.TaintEffectNoExecute,
	}

	tests := []struct {
		Name           string
		Node           *v1.Node
		ExpectedTaints []*v1.Taint
	}{
		{
			Name: "Node with deprecated not-ready taint key",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						*deprecatedNotReadyTaint,
					},
				},
			},
			ExpectedTaints: []*v1.Taint{nodeNotReadyTaint},
		},
		{
			Name: "Node with deprecated unreachable taint key",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						*deprecatedUnreachableTaint,
					},
				},
			},
			ExpectedTaints: []*v1.Taint{nodeUnreachableTaint},
		},
		{
			Name: "Node with not-ready taint key",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						*nodeNotReadyTaint,
					},
				},
			},
			ExpectedTaints: []*v1.Taint{nodeNotReadyTaint},
		},
		{
			Name: "Node with unreachable taint key",
			Node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "node0",
					CreationTimestamp: metav1.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
					Labels: map[string]string{
						kubeletapis.LabelZoneRegion:        "region1",
						kubeletapis.LabelZoneFailureDomain: "zone1",
					},
				},
				Spec: v1.NodeSpec{
					Taints: []v1.Taint{
						*nodeUnreachableTaint,
					},
				},
			},
			ExpectedTaints: []*v1.Taint{nodeUnreachableTaint},
		},
	}

	for _, test := range tests {
		fakeNodeHandler.Update(test.Node)
		if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		nodeController.doFixDeprecatedTaintKeyPass(test.Node)
		if err := nodeController.syncNodeStore(fakeNodeHandler); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		node, err := nodeController.nodeLister.Get(test.Node.GetName())
		if err != nil {
			t.Errorf("Can't get current node...")
			return
		}
		if len(node.Spec.Taints) != len(test.ExpectedTaints) {
			t.Errorf("%s: Unexpected number of taints: expected %d, got %d",
				test.Name, len(test.ExpectedTaints), len(node.Spec.Taints))
		}
		for _, taint := range test.ExpectedTaints {
			if !taintutils.TaintExists(node.Spec.Taints, taint) {
				t.Errorf("%s: Can't find taint %v in %v", test.Name, taint, node.Spec.Taints)
			}
		}
	}
}
