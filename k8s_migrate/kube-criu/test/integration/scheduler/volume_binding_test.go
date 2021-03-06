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

// This file tests the VolumeScheduling feature.

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/controller/volume/persistentvolume"
	persistentvolumeoptions "k8s.io/kubernetes/pkg/controller/volume/persistentvolume/options"
	"k8s.io/kubernetes/pkg/volume"
	volumetest "k8s.io/kubernetes/pkg/volume/testing"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

type testConfig struct {
	client   clientset.Interface
	ns       string
	stop     <-chan struct{}
	teardown func()
}

var (
	// Delete API objects immediately
	deletePeriod = int64(0)
	deleteOption = &metav1.DeleteOptions{GracePeriodSeconds: &deletePeriod}

	modeWait      = storagev1.VolumeBindingWaitForFirstConsumer
	modeImmediate = storagev1.VolumeBindingImmediate

	classWait      = "wait"
	classImmediate = "immediate"

	sharedClasses = map[storagev1.VolumeBindingMode]*storagev1.StorageClass{
		modeImmediate: makeStorageClass(classImmediate, &modeImmediate),
		modeWait:      makeStorageClass(classWait, &modeWait),
	}
)

const (
	node1                 = "node-1"
	node2                 = "node-2"
	podLimit              = 100
	volsPerPod            = 5
	nodeAffinityLabelKey  = "kubernetes.io/hostname"
	provisionerPluginName = "kubernetes.io/mock-provisioner"
)

type testPV struct {
	name        string
	scMode      storagev1.VolumeBindingMode
	preboundPVC string
	node        string
}

type testPVC struct {
	name       string
	scMode     storagev1.VolumeBindingMode
	preboundPV string
}

func TestVolumeBinding(t *testing.T) {
	features := map[string]bool{
		"VolumeScheduling":       true,
		"PersistentLocalVolumes": true,
	}
	config := setupCluster(t, "volume-scheduling", 2, features, 0)
	defer config.teardown()

	cases := map[string]struct {
		pod  *v1.Pod
		pvs  []*testPV
		pvcs []*testPVC
		// Create these, but they should not be bound in the end
		unboundPvcs []*testPVC
		unboundPvs  []*testPV
		shouldFail  bool
	}{
		"immediate can bind": {
			pod:  makePod("pod-i-canbind", config.ns, []string{"pvc-i-canbind"}),
			pvs:  []*testPV{{"pv-i-canbind", modeImmediate, "", node1}},
			pvcs: []*testPVC{{"pvc-i-canbind", modeImmediate, ""}},
		},
		"immediate cannot bind": {
			pod:         makePod("pod-i-cannotbind", config.ns, []string{"pvc-i-cannotbind"}),
			unboundPvcs: []*testPVC{{"pvc-i-cannotbind", modeImmediate, ""}},
			shouldFail:  true,
		},
		"immediate pvc prebound": {
			pod:  makePod("pod-i-pvc-prebound", config.ns, []string{"pvc-i-prebound"}),
			pvs:  []*testPV{{"pv-i-pvc-prebound", modeImmediate, "", node1}},
			pvcs: []*testPVC{{"pvc-i-prebound", modeImmediate, "pv-i-pvc-prebound"}},
		},
		"immediate pv prebound": {
			pod:  makePod("pod-i-pv-prebound", config.ns, []string{"pvc-i-pv-prebound"}),
			pvs:  []*testPV{{"pv-i-prebound", modeImmediate, "pvc-i-pv-prebound", node1}},
			pvcs: []*testPVC{{"pvc-i-pv-prebound", modeImmediate, ""}},
		},
		"wait can bind": {
			pod:  makePod("pod-w-canbind", config.ns, []string{"pvc-w-canbind"}),
			pvs:  []*testPV{{"pv-w-canbind", modeWait, "", node1}},
			pvcs: []*testPVC{{"pvc-w-canbind", modeWait, ""}},
		},
		"wait cannot bind": {
			pod:         makePod("pod-w-cannotbind", config.ns, []string{"pvc-w-cannotbind"}),
			unboundPvcs: []*testPVC{{"pvc-w-cannotbind", modeWait, ""}},
			shouldFail:  true,
		},
		"wait pvc prebound": {
			pod:  makePod("pod-w-pvc-prebound", config.ns, []string{"pvc-w-prebound"}),
			pvs:  []*testPV{{"pv-w-pvc-prebound", modeWait, "", node1}},
			pvcs: []*testPVC{{"pvc-w-prebound", modeWait, "pv-w-pvc-prebound"}},
		},
		"wait pv prebound": {
			pod:  makePod("pod-w-pv-prebound", config.ns, []string{"pvc-w-pv-prebound"}),
			pvs:  []*testPV{{"pv-w-prebound", modeWait, "pvc-w-pv-prebound", node1}},
			pvcs: []*testPVC{{"pvc-w-pv-prebound", modeWait, ""}},
		},
		"wait can bind two": {
			pod: makePod("pod-w-canbind-2", config.ns, []string{"pvc-w-canbind-2", "pvc-w-canbind-3"}),
			pvs: []*testPV{
				{"pv-w-canbind-2", modeWait, "", node2},
				{"pv-w-canbind-3", modeWait, "", node2},
			},
			pvcs: []*testPVC{
				{"pvc-w-canbind-2", modeWait, ""},
				{"pvc-w-canbind-3", modeWait, ""},
			},
			unboundPvs: []*testPV{
				{"pv-w-canbind-5", modeWait, "", node1},
			},
		},
		"wait cannot bind two": {
			pod: makePod("pod-w-cannotbind-2", config.ns, []string{"pvc-w-cannotbind-1", "pvc-w-cannotbind-2"}),
			unboundPvcs: []*testPVC{
				{"pvc-w-cannotbind-1", modeWait, ""},
				{"pvc-w-cannotbind-2", modeWait, ""},
			},
			unboundPvs: []*testPV{
				{"pv-w-cannotbind-1", modeWait, "", node2},
				{"pv-w-cannotbind-2", modeWait, "", node1},
			},
			shouldFail: true,
		},
		"mix immediate and wait": {
			pod: makePod("pod-mix-bound", config.ns, []string{"pvc-w-canbind-4", "pvc-i-canbind-2"}),
			pvs: []*testPV{
				{"pv-w-canbind-4", modeWait, "", node1},
				{"pv-i-canbind-2", modeImmediate, "", node1},
			},
			pvcs: []*testPVC{
				{"pvc-w-canbind-4", modeWait, ""},
				{"pvc-i-canbind-2", modeImmediate, ""},
			},
		},
	}

	for name, test := range cases {
		glog.Infof("Running test %v", name)

		// Create two StorageClasses
		suffix := rand.String(4)
		classes := map[storagev1.VolumeBindingMode]*storagev1.StorageClass{}
		classes[modeImmediate] = makeStorageClass(fmt.Sprintf("immediate-%v", suffix), &modeImmediate)
		classes[modeWait] = makeStorageClass(fmt.Sprintf("wait-%v", suffix), &modeWait)
		for _, sc := range classes {
			if _, err := config.client.StorageV1().StorageClasses().Create(sc); err != nil {
				t.Fatalf("Failed to create StorageClass %q: %v", sc.Name, err)
			}
		}

		// Create PVs
		for _, pvConfig := range test.pvs {
			pv := makePV(pvConfig.name, classes[pvConfig.scMode].Name, pvConfig.preboundPVC, config.ns, pvConfig.node)
			if _, err := config.client.CoreV1().PersistentVolumes().Create(pv); err != nil {
				t.Fatalf("Failed to create PersistentVolume %q: %v", pv.Name, err)
			}
		}

		for _, pvConfig := range test.unboundPvs {
			pv := makePV(pvConfig.name, classes[pvConfig.scMode].Name, pvConfig.preboundPVC, config.ns, pvConfig.node)
			if _, err := config.client.CoreV1().PersistentVolumes().Create(pv); err != nil {
				t.Fatalf("Failed to create PersistentVolume %q: %v", pv.Name, err)
			}
		}

		// Create PVCs
		for _, pvcConfig := range test.pvcs {
			pvc := makePVC(pvcConfig.name, config.ns, &classes[pvcConfig.scMode].Name, pvcConfig.preboundPV)
			if _, err := config.client.CoreV1().PersistentVolumeClaims(config.ns).Create(pvc); err != nil {
				t.Fatalf("Failed to create PersistentVolumeClaim %q: %v", pvc.Name, err)
			}
		}
		for _, pvcConfig := range test.unboundPvcs {
			pvc := makePVC(pvcConfig.name, config.ns, &classes[pvcConfig.scMode].Name, pvcConfig.preboundPV)
			if _, err := config.client.CoreV1().PersistentVolumeClaims(config.ns).Create(pvc); err != nil {
				t.Fatalf("Failed to create PersistentVolumeClaim %q: %v", pvc.Name, err)
			}
		}

		// Create Pod
		if _, err := config.client.CoreV1().Pods(config.ns).Create(test.pod); err != nil {
			t.Fatalf("Failed to create Pod %q: %v", test.pod.Name, err)
		}
		if test.shouldFail {
			if err := waitForPodUnschedulable(config.client, test.pod); err != nil {
				t.Errorf("Pod %q was not unschedulable: %v", test.pod.Name, err)
			}
		} else {
			if err := waitForPodToSchedule(config.client, test.pod); err != nil {
				t.Errorf("Failed to schedule Pod %q: %v", test.pod.Name, err)
			}
		}

		// Validate PVC/PV binding
		for _, pvc := range test.pvcs {
			validatePVCPhase(t, config.client, pvc.name, config.ns, v1.ClaimBound)
		}
		for _, pvc := range test.unboundPvcs {
			validatePVCPhase(t, config.client, pvc.name, config.ns, v1.ClaimPending)
		}
		for _, pv := range test.pvs {
			validatePVPhase(t, config.client, pv.name, v1.VolumeBound)
		}
		for _, pv := range test.unboundPvs {
			validatePVPhase(t, config.client, pv.name, v1.VolumeAvailable)
		}

		// Force delete objects, but they still may not be immediately removed
		deleteTestObjects(config.client, config.ns, deleteOption)
	}
}

// TestVolumeBindingRescheduling tests scheduler will retry scheduling when needed.
func TestVolumeBindingRescheduling(t *testing.T) {
	features := map[string]bool{
		"VolumeScheduling":              true,
		"PersistentLocalVolumes":        true,
		"DynamicProvisioningScheduling": true,
	}
	config := setupCluster(t, "volume-scheduling", 2, features, 0)
	defer config.teardown()

	storageClassName := "local-storage"

	cases := map[string]struct {
		pod        *v1.Pod
		pvcs       []*testPVC
		pvs        []*testPV
		trigger    func(config *testConfig)
		shouldFail bool
	}{
		"reschedule on WaitForFirstConsumer dynamic storage class add": {
			pod: makePod("pod-reschedule-onclassadd-dynamic", config.ns, []string{"pvc-reschedule-onclassadd-dynamic"}),
			pvcs: []*testPVC{
				{"pvc-reschedule-onclassadd-dynamic", "", ""},
			},
			trigger: func(config *testConfig) {
				sc := makeDynamicProvisionerStorageClass(storageClassName, &modeWait)
				if _, err := config.client.StorageV1().StorageClasses().Create(sc); err != nil {
					t.Fatalf("Failed to create StorageClass %q: %v", sc.Name, err)
				}
			},
			shouldFail: false,
		},
		"reschedule on WaitForFirstConsumer static storage class add": {
			pod: makePod("pod-reschedule-onclassadd-static", config.ns, []string{"pvc-reschedule-onclassadd-static"}),
			pvcs: []*testPVC{
				{"pvc-reschedule-onclassadd-static", "", ""},
			},
			trigger: func(config *testConfig) {
				sc := makeStorageClass(storageClassName, &modeWait)
				if _, err := config.client.StorageV1().StorageClasses().Create(sc); err != nil {
					t.Fatalf("Failed to create StorageClass %q: %v", sc.Name, err)
				}
				// Create pv for this class to mock static provisioner behavior.
				pv := makePV("pv-reschedule-onclassadd-static", storageClassName, "", "", node1)
				if pv, err := config.client.CoreV1().PersistentVolumes().Create(pv); err != nil {
					t.Fatalf("Failed to create PersistentVolume %q: %v", pv.Name, err)
				}
			},
			shouldFail: false,
		},
		"reschedule on delay binding PVC add": {
			pod: makePod("pod-reschedule-onpvcadd", config.ns, []string{"pvc-reschedule-onpvcadd"}),
			pvs: []*testPV{
				{
					name:   "pv-reschedule-onpvcadd",
					scMode: modeWait,
					node:   node1,
				},
			},
			trigger: func(config *testConfig) {
				pvc := makePVC("pvc-reschedule-onpvcadd", config.ns, &classWait, "")
				if _, err := config.client.CoreV1().PersistentVolumeClaims(config.ns).Create(pvc); err != nil {
					t.Fatalf("Failed to create PersistentVolumeClaim %q: %v", pvc.Name, err)
				}
			},
			shouldFail: false,
		},
	}

	for name, test := range cases {
		glog.Infof("Running test %v", name)

		if test.pod == nil {
			t.Fatal("pod is required for this test")
		}

		// Create unbound pvc
		for _, pvcConfig := range test.pvcs {
			pvc := makePVC(pvcConfig.name, config.ns, &storageClassName, "")
			if _, err := config.client.CoreV1().PersistentVolumeClaims(config.ns).Create(pvc); err != nil {
				t.Fatalf("Failed to create PersistentVolumeClaim %q: %v", pvc.Name, err)
			}
		}

		// Create PVs
		for _, pvConfig := range test.pvs {
			pv := makePV(pvConfig.name, sharedClasses[pvConfig.scMode].Name, pvConfig.preboundPVC, config.ns, pvConfig.node)
			if _, err := config.client.CoreV1().PersistentVolumes().Create(pv); err != nil {
				t.Fatalf("Failed to create PersistentVolume %q: %v", pv.Name, err)
			}
		}

		// Create pod
		if _, err := config.client.CoreV1().Pods(config.ns).Create(test.pod); err != nil {
			t.Fatalf("Failed to create Pod %q: %v", test.pod.Name, err)
		}

		// Wait for pod is unschedulable.
		glog.Infof("Waiting for pod is unschedulable")
		if err := waitForPodUnschedulable(config.client, test.pod); err != nil {
			t.Errorf("Failed as Pod %s was not unschedulable: %v", test.pod.Name, err)
		}

		// Trigger
		test.trigger(config)

		// Wait for pod is scheduled or unscheduable.
		if !test.shouldFail {
			glog.Infof("Waiting for pod is scheduled")
			if err := waitForPodToSchedule(config.client, test.pod); err != nil {
				t.Errorf("Failed to schedule Pod %q: %v", test.pod.Name, err)
			}
		} else {
			glog.Infof("Waiting for pod is unschedulable")
			if err := waitForPodUnschedulable(config.client, test.pod); err != nil {
				t.Errorf("Failed as Pod %s was not unschedulable: %v", test.pod.Name, err)
			}
		}

		// Force delete objects, but they still may not be immediately removed
		deleteTestObjects(config.client, config.ns, deleteOption)
	}
}

// TestVolumeBindingStress creates <podLimit> pods, each with <volsPerPod> unbound PVCs.
func TestVolumeBindingStress(t *testing.T) {
	testVolumeBindingStress(t, 0)
}

// Like TestVolumeBindingStress but with scheduler resync. In real cluster,
// scheduler will schedule failed pod frequently due to various events, e.g.
// service/node update events.
// This is useful to detect possible race conditions.
func TestVolumeBindingStressWithSchedulerResync(t *testing.T) {
	testVolumeBindingStress(t, time.Second)
}

func testVolumeBindingStress(t *testing.T, schedulerResyncPeriod time.Duration) {
	features := map[string]bool{
		"VolumeScheduling":       true,
		"PersistentLocalVolumes": true,
	}
	config := setupCluster(t, "volume-binding-stress", 1, features, schedulerResyncPeriod)
	defer config.teardown()

	// Create enough PVs and PVCs for all the pods
	pvs := []*v1.PersistentVolume{}
	pvcs := []*v1.PersistentVolumeClaim{}
	for i := 0; i < podLimit*volsPerPod; i++ {
		pv := makePV(fmt.Sprintf("pv-stress-%v", i), classWait, "", "", node1)
		pvc := makePVC(fmt.Sprintf("pvc-stress-%v", i), config.ns, &classWait, "")

		if pv, err := config.client.CoreV1().PersistentVolumes().Create(pv); err != nil {
			t.Fatalf("Failed to create PersistentVolume %q: %v", pv.Name, err)
		}
		if pvc, err := config.client.CoreV1().PersistentVolumeClaims(config.ns).Create(pvc); err != nil {
			t.Fatalf("Failed to create PersistentVolumeClaim %q: %v", pvc.Name, err)
		}

		pvs = append(pvs, pv)
		pvcs = append(pvcs, pvc)
	}

	pods := []*v1.Pod{}
	for i := 0; i < podLimit; i++ {
		// Generate string of all the PVCs for the pod
		podPvcs := []string{}
		for j := i * volsPerPod; j < (i+1)*volsPerPod; j++ {
			podPvcs = append(podPvcs, pvcs[j].Name)
		}

		pod := makePod(fmt.Sprintf("pod%v", i), config.ns, podPvcs)
		if pod, err := config.client.CoreV1().Pods(config.ns).Create(pod); err != nil {
			t.Fatalf("Failed to create Pod %q: %v", pod.Name, err)
		}
		pods = append(pods, pod)
	}

	// Validate Pods scheduled
	for _, pod := range pods {
		// Use increased timeout for stress test because there is a higher chance of
		// PV sync error
		if err := waitForPodToScheduleWithTimeout(config.client, pod, 60*time.Second); err != nil {
			t.Errorf("Failed to schedule Pod %q: %v", pod.Name, err)
		}
	}

	// Validate PVC/PV binding
	for _, pvc := range pvcs {
		validatePVCPhase(t, config.client, pvc.Name, config.ns, v1.ClaimBound)
	}
	for _, pv := range pvs {
		validatePVPhase(t, config.client, pv.Name, v1.VolumeBound)
	}
}

func TestPVAffinityConflict(t *testing.T) {
	features := map[string]bool{
		"VolumeScheduling":       true,
		"PersistentLocalVolumes": true,
	}
	config := setupCluster(t, "volume-scheduling", 3, features, 0)
	defer config.teardown()

	pv := makePV("local-pv", classImmediate, "", "", node1)
	pvc := makePVC("local-pvc", config.ns, &classImmediate, "")

	// Create PV
	if _, err := config.client.CoreV1().PersistentVolumes().Create(pv); err != nil {
		t.Fatalf("Failed to create PersistentVolume %q: %v", pv.Name, err)
	}

	// Create PVC
	if _, err := config.client.CoreV1().PersistentVolumeClaims(config.ns).Create(pvc); err != nil {
		t.Fatalf("Failed to create PersistentVolumeClaim %q: %v", pvc.Name, err)
	}

	// Wait for PVC bound
	if err := waitForPVCBound(config.client, pvc); err != nil {
		t.Fatalf("PVC %q failed to bind: %v", pvc.Name, err)
	}

	nodeMarkers := []interface{}{
		markNodeAffinity,
		markNodeSelector,
	}
	for i := 0; i < len(nodeMarkers); i++ {
		podName := "local-pod-" + strconv.Itoa(i+1)
		pod := makePod(podName, config.ns, []string{"local-pvc"})
		nodeMarkers[i].(func(*v1.Pod, string))(pod, "node-2")
		// Create Pod
		if _, err := config.client.CoreV1().Pods(config.ns).Create(pod); err != nil {
			t.Fatalf("Failed to create Pod %q: %v", pod.Name, err)
		}
		// Give time to shceduler to attempt to schedule pod
		if err := waitForPodUnschedulable(config.client, pod); err != nil {
			t.Errorf("Failed as Pod %s was not unschedulable: %v", pod.Name, err)
		}
		// Check pod conditions
		p, err := config.client.CoreV1().Pods(config.ns).Get(podName, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to access Pod %s status: %v", podName, err)
		}
		if strings.Compare(string(p.Status.Phase), "Pending") != 0 {
			t.Fatalf("Failed as Pod %s was in: %s state and not in expected: Pending state", podName, p.Status.Phase)
		}
		if strings.Compare(p.Status.Conditions[0].Reason, "Unschedulable") != 0 {
			t.Fatalf("Failed as Pod %s reason was: %s but expected: Unschedulable", podName, p.Status.Conditions[0].Reason)
		}
		if !strings.Contains(p.Status.Conditions[0].Message, "node(s) didn't match node selector") || !strings.Contains(p.Status.Conditions[0].Message, "node(s) had volume node affinity conflict") {
			t.Fatalf("Failed as Pod's %s failure message does not contain expected message: node(s) didn't match node selector, node(s) had volume node affinity conflict. Got message %q", podName, p.Status.Conditions[0].Message)
		}
		// Deleting test pod
		if err := config.client.CoreV1().Pods(config.ns).Delete(podName, &metav1.DeleteOptions{}); err != nil {
			t.Fatalf("Failed to delete Pod %s: %v", podName, err)
		}
	}
}

func setupCluster(t *testing.T, nsName string, numberOfNodes int, features map[string]bool, resyncPeriod time.Duration) *testConfig {
	oldFeatures := make(map[string]bool, len(features))
	for feature := range features {
		oldFeatures[feature] = utilfeature.DefaultFeatureGate.Enabled(utilfeature.Feature(feature))
	}
	// Set feature gates
	utilfeature.DefaultFeatureGate.SetFromMap(features)

	controllerCh := make(chan struct{})

	context := initTestSchedulerWithOptions(t, initTestMaster(t, nsName, nil), controllerCh, false, nil, false, resyncPeriod)

	clientset := context.clientSet
	ns := context.ns.Name
	// Informers factory for controllers, we disable resync period for testing.
	informerFactory := informers.NewSharedInformerFactory(context.clientSet, 0)

	// Start PV controller for volume binding.
	host := volumetest.NewFakeVolumeHost("/tmp/fake", nil, nil)
	plugin := &volumetest.FakeVolumePlugin{
		PluginName:             provisionerPluginName,
		Host:                   host,
		Config:                 volume.VolumeConfig{},
		LastProvisionerOptions: volume.VolumeOptions{},
		NewAttacherCallCount:   0,
		NewDetacherCallCount:   0,
		Mounters:               nil,
		Unmounters:             nil,
		Attachers:              nil,
		Detachers:              nil,
	}
	plugins := []volume.VolumePlugin{plugin}

	controllerOptions := persistentvolumeoptions.NewPersistentVolumeControllerOptions()
	params := persistentvolume.ControllerParameters{
		KubeClient:                clientset,
		SyncPeriod:                controllerOptions.PVClaimBinderSyncPeriod,
		VolumePlugins:             plugins,
		Cloud:                     nil,
		ClusterName:               "volume-test-cluster",
		VolumeInformer:            informerFactory.Core().V1().PersistentVolumes(),
		ClaimInformer:             informerFactory.Core().V1().PersistentVolumeClaims(),
		ClassInformer:             informerFactory.Storage().V1().StorageClasses(),
		PodInformer:               informerFactory.Core().V1().Pods(),
		NodeInformer:              informerFactory.Core().V1().Nodes(),
		EnableDynamicProvisioning: true,
	}
	ctrl, err := persistentvolume.NewController(params)
	if err != nil {
		t.Fatalf("Failed to create PV controller: %v", err)
	}
	go ctrl.Run(controllerCh)
	// Start informer factory after all controllers are configured and running.
	informerFactory.Start(controllerCh)
	informerFactory.WaitForCacheSync(controllerCh)

	// Create shared objects
	// Create nodes
	for i := 0; i < numberOfNodes; i++ {
		testNode := &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("node-%d", i+1),
				Labels: map[string]string{nodeAffinityLabelKey: fmt.Sprintf("node-%d", i+1)},
			},
			Spec: v1.NodeSpec{Unschedulable: false},
			Status: v1.NodeStatus{
				Capacity: v1.ResourceList{
					v1.ResourcePods: *resource.NewQuantity(podLimit, resource.DecimalSI),
				},
				Conditions: []v1.NodeCondition{
					{
						Type:              v1.NodeReady,
						Status:            v1.ConditionTrue,
						Reason:            fmt.Sprintf("schedulable condition"),
						LastHeartbeatTime: metav1.Time{Time: time.Now()},
					},
				},
			},
		}
		if _, err := clientset.CoreV1().Nodes().Create(testNode); err != nil {
			t.Fatalf("Failed to create Node %q: %v", testNode.Name, err)
		}
	}

	// Create SCs
	for _, sc := range sharedClasses {
		if _, err := clientset.StorageV1().StorageClasses().Create(sc); err != nil {
			t.Fatalf("Failed to create StorageClass %q: %v", sc.Name, err)
		}
	}

	return &testConfig{
		client: clientset,
		ns:     ns,
		stop:   controllerCh,
		teardown: func() {
			deleteTestObjects(clientset, ns, nil)
			cleanupTest(t, context)
			// Restore feature gates
			utilfeature.DefaultFeatureGate.SetFromMap(oldFeatures)
		},
	}
}

func deleteTestObjects(client clientset.Interface, ns string, option *metav1.DeleteOptions) {
	client.CoreV1().Pods(ns).DeleteCollection(option, metav1.ListOptions{})
	client.CoreV1().PersistentVolumeClaims(ns).DeleteCollection(option, metav1.ListOptions{})
	client.CoreV1().PersistentVolumes().DeleteCollection(option, metav1.ListOptions{})
	client.StorageV1().StorageClasses().DeleteCollection(option, metav1.ListOptions{})
}

func makeStorageClass(name string, mode *storagev1.VolumeBindingMode) *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner:       "kubernetes.io/no-provisioner",
		VolumeBindingMode: mode,
	}
}

func makeDynamicProvisionerStorageClass(name string, mode *storagev1.VolumeBindingMode) *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner:       provisionerPluginName,
		VolumeBindingMode: mode,
	}
}

func makePV(name, scName, pvcName, ns, node string) *v1.PersistentVolume {
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{},
		},
		Spec: v1.PersistentVolumeSpec{
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): resource.MustParse("5Gi"),
			},
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			StorageClassName: scName,
			PersistentVolumeSource: v1.PersistentVolumeSource{
				Local: &v1.LocalVolumeSource{
					Path: "/test-path",
				},
			},
			NodeAffinity: &v1.VolumeNodeAffinity{
				Required: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{
						{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      nodeAffinityLabelKey,
									Operator: v1.NodeSelectorOpIn,
									Values:   []string{node},
								},
							},
						},
					},
				},
			},
		},
	}

	if pvcName != "" {
		pv.Spec.ClaimRef = &v1.ObjectReference{Name: pvcName, Namespace: ns}
	}

	return pv
}

func makePVC(name, ns string, scName *string, volumeName string) *v1.PersistentVolumeClaim {
	return &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("5Gi"),
				},
			},
			StorageClassName: scName,
			VolumeName:       volumeName,
		},
	}
}

func makePod(name, ns string, pvcs []string) *v1.Pod {
	volumes := []v1.Volume{}
	for i, pvc := range pvcs {
		volumes = append(volumes, v1.Volume{
			Name: fmt.Sprintf("vol%v", i),
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvc,
				},
			},
		})
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "write-pod",
					Image:   imageutils.GetE2EImage(imageutils.BusyBox),
					Command: []string{"/bin/sh"},
					Args:    []string{"-c", "while true; do sleep 1; done"},
				},
			},
			Volumes: volumes,
		},
	}
}

func validatePVCPhase(t *testing.T, client clientset.Interface, pvcName string, ns string, phase v1.PersistentVolumeClaimPhase) {
	claim, err := client.CoreV1().PersistentVolumeClaims(ns).Get(pvcName, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Failed to get PVC %v/%v: %v", ns, pvcName, err)
	}

	if claim.Status.Phase != phase {
		t.Errorf("PVC %v/%v phase not %v, got %v", ns, pvcName, phase, claim.Status.Phase)
	}
}

func validatePVPhase(t *testing.T, client clientset.Interface, pvName string, phase v1.PersistentVolumePhase) {
	pv, err := client.CoreV1().PersistentVolumes().Get(pvName, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Failed to get PV %v: %v", pvName, err)
	}

	if pv.Status.Phase != phase {
		t.Errorf("PV %v phase not %v, got %v", pvName, phase, pv.Status.Phase)
	}
}

func waitForPVCBound(client clientset.Interface, pvc *v1.PersistentVolumeClaim) error {
	return wait.Poll(time.Second, 30*time.Second, func() (bool, error) {
		claim, err := client.CoreV1().PersistentVolumeClaims(pvc.Namespace).Get(pvc.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if claim.Status.Phase == v1.ClaimBound {
			return true, nil
		}
		return false, nil
	})
}

func markNodeAffinity(pod *v1.Pod, node string) {
	affinity := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
				NodeSelectorTerms: []v1.NodeSelectorTerm{
					{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      nodeAffinityLabelKey,
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{node},
							},
						},
					},
				},
			},
		},
	}
	pod.Spec.Affinity = affinity
}

func markNodeSelector(pod *v1.Pod, node string) {
	ns := map[string]string{
		nodeAffinityLabelKey: node,
	}
	pod.Spec.NodeSelector = ns
}
