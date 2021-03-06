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

package common

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/watch"
	watchtools "k8s.io/client-go/tools/watch"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = framework.KubeDescribe("InitContainer [NodeConformance]", func() {
	f := framework.NewDefaultFramework("init-container")
	var podClient *framework.PodClient
	BeforeEach(func() {
		podClient = f.PodClient()
	})

	/*
		Release: v1.12
		Testname: init-container-starts-app-restartnever-pod
		Description: Ensure that all InitContainers are started
		and all containers in pod are voluntarily terminated with exit status 0,
		and the system is not going to restart any of these containers
		when Pod has restart policy as RestartNever.
	*/
	framework.ConformanceIt("should invoke init containers on a RestartNever pod", func() {
		By("creating the pod")
		name := "pod-init-" + string(uuid.NewUUID())
		value := strconv.Itoa(time.Now().Nanosecond())
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name": "foo",
					"time": value,
				},
			},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				InitContainers: []v1.Container{
					{
						Name:    "init1",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
					{
						Name:    "init2",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
				},
				Containers: []v1.Container{
					{
						Name:    "run1",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
				},
			},
		}
		framework.Logf("PodSpec: initContainers in spec.initContainers")
		startedPod := podClient.Create(pod)
		w, err := podClient.Watch(metav1.SingleObject(startedPod.ObjectMeta))
		Expect(err).NotTo(HaveOccurred(), "error watching a pod")
		wr := watch.NewRecorder(w)
		ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), framework.PodStartTimeout)
		defer cancel()
		event, err := watchtools.UntilWithoutRetry(ctx, wr, conditions.PodCompleted)
		Expect(err).To(BeNil())
		framework.CheckInvariants(wr.Events(), framework.ContainerInitInvariant)
		endPod := event.Object.(*v1.Pod)
		Expect(endPod.Status.Phase).To(Equal(v1.PodSucceeded))
		_, init := podutil.GetPodCondition(&endPod.Status, v1.PodInitialized)
		Expect(init).NotTo(BeNil())
		Expect(init.Status).To(Equal(v1.ConditionTrue))

		Expect(len(endPod.Status.InitContainerStatuses)).To(Equal(2))
		for _, status := range endPod.Status.InitContainerStatuses {
			Expect(status.Ready).To(BeTrue())
			Expect(status.State.Terminated).NotTo(BeNil())
			Expect(status.State.Terminated.ExitCode).To(BeZero())
		}
	})

	/*
		Release: v1.12
		Testname: init-container-starts-app-restartalways-pod
		Description: Ensure that all InitContainers are started
		and all containers in pod started
		and at least one container is still running or is in the process of being restarted
		when Pod has restart policy as RestartAlways.
	*/
	framework.ConformanceIt("should invoke init containers on a RestartAlways pod", func() {
		By("creating the pod")
		name := "pod-init-" + string(uuid.NewUUID())
		value := strconv.Itoa(time.Now().Nanosecond())
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name": "foo",
					"time": value,
				},
			},
			Spec: v1.PodSpec{
				InitContainers: []v1.Container{
					{
						Name:    "init1",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
					{
						Name:    "init2",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
				},
				Containers: []v1.Container{
					{
						Name:  "run1",
						Image: imageutils.GetPauseImageName(),
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
								v1.ResourceMemory: *resource.NewQuantity(50*1024*1024, resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		framework.Logf("PodSpec: initContainers in spec.initContainers")
		startedPod := podClient.Create(pod)
		w, err := podClient.Watch(metav1.SingleObject(startedPod.ObjectMeta))
		Expect(err).NotTo(HaveOccurred(), "error watching a pod")
		wr := watch.NewRecorder(w)
		ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), framework.PodStartTimeout)
		defer cancel()
		event, err := watchtools.UntilWithoutRetry(ctx, wr, conditions.PodRunning)
		Expect(err).To(BeNil())
		framework.CheckInvariants(wr.Events(), framework.ContainerInitInvariant)
		endPod := event.Object.(*v1.Pod)
		Expect(endPod.Status.Phase).To(Equal(v1.PodRunning))
		_, init := podutil.GetPodCondition(&endPod.Status, v1.PodInitialized)
		Expect(init).NotTo(BeNil())
		Expect(init.Status).To(Equal(v1.ConditionTrue))

		Expect(len(endPod.Status.InitContainerStatuses)).To(Equal(2))
		for _, status := range endPod.Status.InitContainerStatuses {
			Expect(status.Ready).To(BeTrue())
			Expect(status.State.Terminated).NotTo(BeNil())
			Expect(status.State.Terminated.ExitCode).To(BeZero())
		}
	})

	/*
		Release: v1.12
		Testname: init-container-fails-stops-app-restartalways-pod
		Description: Ensure that app container is not started
		when all InitContainers failed to start
		and Pod has restarted for few occurrences
		and pod has restart policy as RestartAlways.
	*/
	framework.ConformanceIt("should not start app containers if init containers fail on a RestartAlways pod", func() {
		By("creating the pod")
		name := "pod-init-" + string(uuid.NewUUID())
		value := strconv.Itoa(time.Now().Nanosecond())

		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name": "foo",
					"time": value,
				},
			},
			Spec: v1.PodSpec{
				InitContainers: []v1.Container{
					{
						Name:    "init1",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/false"},
					},
					{
						Name:    "init2",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
				},
				Containers: []v1.Container{
					{
						Name:  "run1",
						Image: imageutils.GetPauseImageName(),
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
								v1.ResourceMemory: *resource.NewQuantity(50*1024*1024, resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		framework.Logf("PodSpec: initContainers in spec.initContainers")
		startedPod := podClient.Create(pod)
		w, err := podClient.Watch(metav1.SingleObject(startedPod.ObjectMeta))
		Expect(err).NotTo(HaveOccurred(), "error watching a pod")

		wr := watch.NewRecorder(w)
		ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), framework.PodStartTimeout)
		defer cancel()
		event, err := watchtools.UntilWithoutRetry(
			ctx, wr,
			// check for the first container to fail at least once
			func(evt watch.Event) (bool, error) {
				switch t := evt.Object.(type) {
				case *v1.Pod:
					for _, status := range t.Status.ContainerStatuses {
						if status.State.Waiting == nil {
							return false, fmt.Errorf("container %q should not be out of waiting: %#v", status.Name, status)
						}
						if status.State.Waiting.Reason != "PodInitializing" {
							return false, fmt.Errorf("container %q should have reason PodInitializing: %#v", status.Name, status)
						}
					}
					if len(t.Status.InitContainerStatuses) != 2 {
						return false, nil
					}
					status := t.Status.InitContainerStatuses[1]
					if status.State.Waiting == nil {
						return false, fmt.Errorf("second init container should not be out of waiting: %#v", status)
					}
					if status.State.Waiting.Reason != "PodInitializing" {
						return false, fmt.Errorf("second init container should have reason PodInitializing: %#v", status)
					}
					status = t.Status.InitContainerStatuses[0]
					if status.State.Terminated != nil && status.State.Terminated.ExitCode == 0 {
						return false, fmt.Errorf("first init container should have exitCode != 0: %#v", status)
					}
					// continue until we see an attempt to restart the pod
					return status.LastTerminationState.Terminated != nil, nil
				default:
					return false, fmt.Errorf("unexpected object: %#v", t)
				}
			},
			// verify we get two restarts
			func(evt watch.Event) (bool, error) {
				switch t := evt.Object.(type) {
				case *v1.Pod:
					status := t.Status.InitContainerStatuses[0]
					if status.RestartCount < 3 {
						return false, nil
					}
					framework.Logf("init container has failed twice: %#v", t)
					// TODO: more conditions
					return true, nil
				default:
					return false, fmt.Errorf("unexpected object: %#v", t)
				}
			},
		)
		Expect(err).To(BeNil())
		framework.CheckInvariants(wr.Events(), framework.ContainerInitInvariant)
		endPod := event.Object.(*v1.Pod)
		Expect(endPod.Status.Phase).To(Equal(v1.PodPending))
		_, init := podutil.GetPodCondition(&endPod.Status, v1.PodInitialized)
		Expect(init).NotTo(BeNil())
		Expect(init.Status).To(Equal(v1.ConditionFalse))
		Expect(init.Reason).To(Equal("ContainersNotInitialized"))
		Expect(init.Message).To(Equal("containers with incomplete status: [init1 init2]"))
		Expect(len(endPod.Status.InitContainerStatuses)).To(Equal(2))
	})

	/*
		Release: v1.12
		Testname: init-container-fails-stops-app-restartnever-pod
		Description: Ensure that app container is not started
		when atleast one InitContainer fails to start and Pod has restart policy as RestartNever.
	*/
	framework.ConformanceIt("should not start app containers and fail the pod if init containers fail on a RestartNever pod", func() {
		By("creating the pod")
		name := "pod-init-" + string(uuid.NewUUID())
		value := strconv.Itoa(time.Now().Nanosecond())
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name": "foo",
					"time": value,
				},
			},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				InitContainers: []v1.Container{
					{
						Name:    "init1",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
					},
					{
						Name:    "init2",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/false"},
					},
				},
				Containers: []v1.Container{
					{
						Name:    "run1",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"/bin/true"},
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
								v1.ResourceMemory: *resource.NewQuantity(50*1024*1024, resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		framework.Logf("PodSpec: initContainers in spec.initContainers")
		startedPod := podClient.Create(pod)

		w, err := podClient.Watch(metav1.SingleObject(startedPod.ObjectMeta))
		Expect(err).NotTo(HaveOccurred(), "error watching a pod")

		wr := watch.NewRecorder(w)
		ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), framework.PodStartTimeout)
		defer cancel()
		event, err := watchtools.UntilWithoutRetry(
			ctx, wr,
			// check for the second container to fail at least once
			func(evt watch.Event) (bool, error) {
				switch t := evt.Object.(type) {
				case *v1.Pod:
					for _, status := range t.Status.ContainerStatuses {
						if status.State.Waiting == nil {
							return false, fmt.Errorf("container %q should not be out of waiting: %#v", status.Name, status)
						}
						if status.State.Waiting.Reason != "PodInitializing" {
							return false, fmt.Errorf("container %q should have reason PodInitializing: %#v", status.Name, status)
						}
					}
					if len(t.Status.InitContainerStatuses) != 2 {
						return false, nil
					}
					status := t.Status.InitContainerStatuses[0]
					if status.State.Terminated == nil {
						if status.State.Waiting != nil && status.State.Waiting.Reason != "PodInitializing" {
							return false, fmt.Errorf("second init container should have reason PodInitializing: %#v", status)
						}
						return false, nil
					}
					if status.State.Terminated != nil && status.State.Terminated.ExitCode != 0 {
						return false, fmt.Errorf("first init container should have exitCode != 0: %#v", status)
					}
					status = t.Status.InitContainerStatuses[1]
					if status.State.Terminated == nil {
						return false, nil
					}
					if status.State.Terminated.ExitCode == 0 {
						return false, fmt.Errorf("second init container should have failed: %#v", status)
					}
					return true, nil
				default:
					return false, fmt.Errorf("unexpected object: %#v", t)
				}
			},
			conditions.PodCompleted,
		)
		Expect(err).To(BeNil())
		framework.CheckInvariants(wr.Events(), framework.ContainerInitInvariant)
		endPod := event.Object.(*v1.Pod)

		Expect(endPod.Status.Phase).To(Equal(v1.PodFailed))
		_, init := podutil.GetPodCondition(&endPod.Status, v1.PodInitialized)
		Expect(init).NotTo(BeNil())
		Expect(init.Status).To(Equal(v1.ConditionFalse))
		Expect(init.Reason).To(Equal("ContainersNotInitialized"))
		Expect(init.Message).To(Equal("containers with incomplete status: [init2]"))
		Expect(len(endPod.Status.InitContainerStatuses)).To(Equal(2))
		Expect(endPod.Status.ContainerStatuses[0].State.Waiting).ToNot(BeNil())
	})
})
