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

package scheduling

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/test/e2e/framework"
	testutils "k8s.io/kubernetes/test/utils"
	imageutils "k8s.io/kubernetes/test/utils/image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	defaultTimeout = 3 * time.Minute
)

// This test requires Rescheduler to be enabled.
var _ = SIGDescribe("Rescheduler [Serial]", func() {
	f := framework.NewDefaultFramework("rescheduler")
	var ns string
	var totalMillicores int

	BeforeEach(func() {
		framework.Skipf("Rescheduler is being deprecated soon in Kubernetes 1.10")
		ns = f.Namespace.Name
		nodes := framework.GetReadySchedulableNodesOrDie(f.ClientSet)
		nodeCount := len(nodes.Items)
		Expect(nodeCount).NotTo(BeZero())

		cpu := nodes.Items[0].Status.Capacity[v1.ResourceCPU]
		totalMillicores = int((&cpu).MilliValue()) * nodeCount
	})

	It("should ensure that critical pod is scheduled in case there is no resources available", func() {
		By("reserving all available cpu")
		err := reserveAllCpu(f, "reserve-all-cpu", totalMillicores)
		defer framework.DeleteRCAndWaitForGC(f.ClientSet, ns, "reserve-all-cpu")
		framework.ExpectNoError(err)

		By("creating a new instance of Dashboard and waiting for Dashboard to be scheduled")
		label := labels.SelectorFromSet(labels.Set(map[string]string{"k8s-app": "kubernetes-dashboard"}))
		listOpts := metav1.ListOptions{LabelSelector: label.String()}
		deployments, err := f.ClientSet.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).List(listOpts)
		framework.ExpectNoError(err)
		Expect(len(deployments.Items)).Should(Equal(1))

		deployment := deployments.Items[0]
		replicas := uint(*(deployment.Spec.Replicas))

		err = framework.ScaleDeployment(f.ClientSet, f.ScalesGetter, metav1.NamespaceSystem, deployment.Name, replicas+1, true)
		defer framework.ExpectNoError(framework.ScaleDeployment(f.ClientSet, f.ScalesGetter, metav1.NamespaceSystem, deployment.Name, replicas, true))
		framework.ExpectNoError(err)

	})
})

func reserveAllCpu(f *framework.Framework, id string, millicores int) error {
	timeout := 5 * time.Minute
	replicas := millicores / 100

	reserveCpu(f, id, 1, 100)
	framework.ExpectNoError(framework.ScaleRC(f.ClientSet, f.ScalesGetter, f.Namespace.Name, id, uint(replicas), false))

	for start := time.Now(); time.Since(start) < timeout; time.Sleep(10 * time.Second) {
		pods, err := framework.GetPodsInNamespace(f.ClientSet, f.Namespace.Name, framework.ImagePullerLabels)
		if err != nil {
			return err
		}

		if len(pods) != replicas {
			continue
		}

		allRunningOrUnschedulable := true
		for _, pod := range pods {
			if !podRunningOrUnschedulable(pod) {
				allRunningOrUnschedulable = false
				break
			}
		}
		if allRunningOrUnschedulable {
			return nil
		}
	}
	return fmt.Errorf("Pod name %s: Gave up waiting %v for %d pods to come up", id, timeout, replicas)
}

func podRunningOrUnschedulable(pod *v1.Pod) bool {
	_, cond := podutil.GetPodCondition(&pod.Status, v1.PodScheduled)
	if cond != nil && cond.Status == v1.ConditionFalse && cond.Reason == "Unschedulable" {
		return true
	}
	running, _ := testutils.PodRunningReady(pod)
	return running
}

func reserveCpu(f *framework.Framework, id string, replicas, millicores int) {
	By(fmt.Sprintf("Running RC which reserves %v millicores", millicores))
	request := int64(millicores / replicas)
	config := &testutils.RCConfig{
		Client:         f.ClientSet,
		InternalClient: f.InternalClientset,
		Name:           id,
		Namespace:      f.Namespace.Name,
		Timeout:        defaultTimeout,
		Image:          imageutils.GetPauseImageName(),
		Replicas:       replicas,
		CpuRequest:     request,
	}
	framework.ExpectNoError(framework.RunRC(*config))
}
