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

package apimachinery

import (
	"fmt"
	"sync/atomic"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionstestserver "k8s.io/apiextensions-apiserver/test/integration/fixtures"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage/names"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

// estimateMaximumPods estimates how many pods the cluster can handle
// with some wiggle room, to prevent pods being unable to schedule due
// to max pod constraints.
func estimateMaximumPods(c clientset.Interface, min, max int32) int32 {
	availablePods := int32(0)
	for _, node := range framework.GetReadySchedulableNodesOrDie(c).Items {
		if q, ok := node.Status.Allocatable["pods"]; ok {
			if num, ok := q.AsInt64(); ok {
				availablePods += int32(num)
				continue
			}
		}
		// best guess per node, since default maxPerCore is 10 and most nodes have at least
		// one core.
		availablePods += 10
	}
	//avoid creating exactly max pods
	availablePods *= 8 / 10
	// bound the top and bottom
	if availablePods > max {
		availablePods = max
	}
	if availablePods < min {
		availablePods = min
	}
	return availablePods
}

func getForegroundOptions() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationForeground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func getBackgroundOptions() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationBackground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func getOrphanOptions() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationOrphan
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

var (
	zero                        = int64(0)
	lablecount                  = int64(0)
	CronJobGroupVersionResource = schema.GroupVersionResource{Group: batchv1beta1.GroupName, Version: "v1beta1", Resource: "cronjobs"}
)

func getPodTemplateSpec(labels map[string]string) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: &zero,
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: imageutils.GetE2EImage(imageutils.Nginx),
				},
			},
		},
	}
}

func newOwnerDeployment(f *framework.Framework, deploymentName string, labels map[string]string) *v1beta1.Deployment {
	replicas := int32(2)
	return &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
			},
			Template: getPodTemplateSpec(labels),
		},
	}
}

func getSelector() map[string]string {
	return map[string]string{"app": "gc-test"}
}

func newOwnerRC(f *framework.Framework, name string, replicas int32, labels map[string]string) *v1.ReplicationController {
	template := getPodTemplateSpec(labels)
	return &v1.ReplicationController{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.Namespace.Name,
			Name:      name,
		},
		Spec: v1.ReplicationControllerSpec{
			Replicas: &replicas,
			Selector: labels,
			Template: &template,
		},
	}
}

// verifyRemainingDeploymentsReplicaSetsPods verifies if the number
// of the remaining deployments, replica set and pods are deploymentNum,
// rsNum and podNum. It returns error if the communication with the API
// server fails.
func verifyRemainingDeploymentsReplicaSetsPods(
	f *framework.Framework,
	clientSet clientset.Interface,
	deployment *v1beta1.Deployment,
	deploymentNum, rsNum, podNum int,
) (bool, error) {
	var ret = true
	rs, err := clientSet.ExtensionsV1beta1().ReplicaSets(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list rs: %v", err)
	}
	if len(rs.Items) != rsNum {
		ret = false
		By(fmt.Sprintf("expected %d rs, got %d rs", rsNum, len(rs.Items)))
	}
	deployments, err := clientSet.ExtensionsV1beta1().Deployments(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list deployments: %v", err)
	}
	if len(deployments.Items) != deploymentNum {
		ret = false
		By(fmt.Sprintf("expected %d Deployments, got %d Deployments", deploymentNum, len(deployments.Items)))
	}
	pods, err := clientSet.CoreV1().Pods(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list pods: %v", err)
	}
	if len(pods.Items) != podNum {
		ret = false
		By(fmt.Sprintf("expected %v Pods, got %d Pods", podNum, len(pods.Items)))
	}

	return ret, nil
}

func newGCPod(name string) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: new(int64),
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: imageutils.GetE2EImage(imageutils.Nginx),
				},
			},
		},
	}
}

// verifyRemainingReplicationControllersPods verifies if the number of the remaining replication
// controllers and pods are rcNum and podNum. It returns error if the
// communication with the API server fails.
func verifyRemainingReplicationControllersPods(f *framework.Framework, clientSet clientset.Interface, rcNum, podNum int) (bool, error) {
	rcClient := clientSet.CoreV1().ReplicationControllers(f.Namespace.Name)
	pods, err := clientSet.CoreV1().Pods(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list pods: %v", err)
	}
	var ret = true
	if len(pods.Items) != podNum {
		ret = false
		By(fmt.Sprintf("expected %d pods, got %d pods", podNum, len(pods.Items)))
	}
	rcs, err := rcClient.List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list replication controllers: %v", err)
	}
	if len(rcs.Items) != rcNum {
		ret = false
		By(fmt.Sprintf("expected %d RCs, got %d RCs", rcNum, len(rcs.Items)))
	}
	return ret, nil
}

// verifyRemainingCronJobsJobsPods verifies if the number of remaining cronjobs,
// jobs and pods. It returns error if the communication with the API server fails.
func verifyRemainingCronJobsJobsPods(f *framework.Framework, clientSet clientset.Interface,
	cjNum, jobNum, podNum int) (bool, error) {
	var ret = true

	cronJobs, err := f.ClientSet.BatchV1beta1().CronJobs(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list cronjobs: %v", err)
	}
	if len(cronJobs.Items) != cjNum {
		ret = false
		By(fmt.Sprintf("expected %d cronjobs, got %d cronjobs", cjNum, len(cronJobs.Items)))
	}

	jobs, err := f.ClientSet.BatchV1().Jobs(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list jobs: %v", err)
	}
	if len(jobs.Items) != jobNum {
		ret = false
		By(fmt.Sprintf("expected %d jobs, got %d jobs", jobNum, len(jobs.Items)))
	}

	pods, err := f.ClientSet.CoreV1().Pods(f.Namespace.Name).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("Failed to list pods: %v", err)
	}
	if len(pods.Items) != podNum {
		ret = false
		By(fmt.Sprintf("expected %d pods, got %d pods", podNum, len(pods.Items)))
	}

	return ret, nil
}

func gatherMetrics(f *framework.Framework) {
	By("Gathering metrics")
	var summary framework.TestDataSummary
	grabber, err := metrics.NewMetricsGrabber(f.ClientSet, f.KubemarkExternalClusterClientSet, false, false, true, false, false)
	if err != nil {
		framework.Logf("Failed to create MetricsGrabber. Skipping metrics gathering.")
	} else {
		received, err := grabber.Grab()
		if err != nil {
			framework.Logf("MetricsGrabber failed grab metrics. Skipping metrics gathering.")
		} else {
			summary = (*framework.MetricsForE2E)(&received)
			framework.Logf(summary.PrintHumanReadable())
		}
	}
}

func newCronJob(name, schedule string) *batchv1beta1.CronJob {
	parallelism := int32(1)
	completions := int32(1)
	return &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "CronJob",
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule: schedule,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Parallelism: &parallelism,
					Completions: &completions,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy:                 v1.RestartPolicyOnFailure,
							TerminationGracePeriodSeconds: &zero,
							Containers: []v1.Container{
								{
									Name:    "c",
									Image:   imageutils.GetE2EImage(imageutils.BusyBox),
									Command: []string{"sleep", "300"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// getUniqLabel returns a UniqLabel based on labeLkey and labelvalue.
func getUniqLabel(labelkey, labelvalue string) map[string]string {
	count := atomic.AddInt64(&lablecount, 1)
	uniqlabelkey := fmt.Sprintf("%s-%05d", labelkey, count)
	uniqlabelvalue := fmt.Sprintf("%s-%05d", labelvalue, count)
	return map[string]string{uniqlabelkey: uniqlabelvalue}
}

var _ = SIGDescribe("Garbage collector", func() {
	f := framework.NewDefaultFramework("gc")

	/*
		Release : v1.9
		Testname: Garbage Collector, delete replication controller, propagation policy background
		Description: Create a replication controller with 2 Pods. Once RC is created and the first Pod is created, delete RC with deleteOptions.PropagationPolicy set to Background. Deleting the Replication Controller MUST cause pods created by that RC to be deleted.
	*/
	framework.ConformanceIt("should delete pods created by rc when not orphaning", func() {
		clientSet := f.ClientSet
		rcClient := clientSet.CoreV1().ReplicationControllers(f.Namespace.Name)
		podClient := clientSet.CoreV1().Pods(f.Namespace.Name)
		rcName := "simpletest.rc"
		uniqLabels := getUniqLabel("gctest", "delete_pods")
		rc := newOwnerRC(f, rcName, 2, uniqLabels)
		By("create the rc")
		rc, err := rcClient.Create(rc)
		if err != nil {
			framework.Failf("Failed to create replication controller: %v", err)
		}
		// wait for rc to create some pods
		if err := wait.Poll(5*time.Second, 30*time.Second, func() (bool, error) {
			pods, err := podClient.List(metav1.ListOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to list pods: %v", err)
			}
			// We intentionally don't wait the number of pods to reach
			// rc.Spec.Replicas. We want to see if the garbage collector and the
			// rc manager work properly if the rc is deleted before it reaches
			// stasis.
			if len(pods.Items) > 0 {
				return true, nil
			} else {
				return false, nil
			}
		}); err != nil {
			framework.Failf("failed to wait for the rc to create some pods: %v", err)
		}
		By("delete the rc")
		deleteOptions := getBackgroundOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(rc.UID))
		if err := rcClient.Delete(rc.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the rc: %v", err)
		}
		By("wait for all pods to be garbage collected")
		// wait for the RCs and Pods to reach the expected numbers.
		if err := wait.Poll(5*time.Second, 60*time.Second, func() (bool, error) {
			return verifyRemainingReplicationControllersPods(f, clientSet, 0, 0)
		}); err != nil {
			framework.Failf("failed to wait for all pods to be deleted: %v", err)
			remainingPods, err := podClient.List(metav1.ListOptions{})
			if err != nil {
				framework.Failf("failed to list pods post mortem: %v", err)
			} else {
				framework.Failf("remaining pods are: %#v", remainingPods)
			}
		}
		gatherMetrics(f)
	})

	/*
		Release : v1.9
		Testname: Garbage Collector, delete replication controller, propagation policy orphan
		Description: Create a replication controller with maximum allocatable Pods between 10 and 100 replicas. Once RC is created and the all Pods are created, delete RC with deleteOptions.PropagationPolicy set to Orphan. Deleting the Replication Controller MUST cause pods created by that RC to be orphaned.
	*/
	framework.ConformanceIt("should orphan pods created by rc if delete options say so", func() {
		clientSet := f.ClientSet
		rcClient := clientSet.CoreV1().ReplicationControllers(f.Namespace.Name)
		podClient := clientSet.CoreV1().Pods(f.Namespace.Name)
		rcName := "simpletest.rc"
		uniqLabels := getUniqLabel("gctest", "orphan_pods")
		rc := newOwnerRC(f, rcName, estimateMaximumPods(clientSet, 10, 100), uniqLabels)
		By("create the rc")
		rc, err := rcClient.Create(rc)
		if err != nil {
			framework.Failf("Failed to create replication controller: %v", err)
		}
		// wait for rc to create pods
		if err := wait.Poll(5*time.Second, 30*time.Second, func() (bool, error) {
			rc, err := rcClient.Get(rc.Name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to get rc: %v", err)
			}
			if rc.Status.Replicas == *rc.Spec.Replicas {
				return true, nil
			} else {
				return false, nil
			}
		}); err != nil {
			framework.Failf("failed to wait for the rc.Status.Replicas to reach rc.Spec.Replicas: %v", err)
		}
		By("delete the rc")
		deleteOptions := getOrphanOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(rc.UID))
		if err := rcClient.Delete(rc.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the rc: %v", err)
		}
		By("wait for the rc to be deleted")
		// Orphaning the 100 pods takes 100 PATCH operations. The default qps of
		// a client is 5. If the qps is saturated, it will take 20s to orphan
		// the pods. However, apiserver takes hundreds of ms to finish one
		// PATCH, and the gc sends the patching in a single thread, so the
		// actual qps is less than 5. Also, the e2e tests are running in
		// parallel, the GC controller might get distracted by other tests.
		// According to the test logs, 120s is enough time.
		if err := wait.Poll(5*time.Second, 120*time.Second, func() (bool, error) {
			rcs, err := rcClient.List(metav1.ListOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to list rcs: %v", err)
			}
			if len(rcs.Items) != 0 {
				return false, nil
			}
			return true, nil
		}); err != nil {
			framework.Failf("%v", err)
		}
		By("wait for 30 seconds to see if the garbage collector mistakenly deletes the pods")
		time.Sleep(30 * time.Second)
		pods, err := podClient.List(metav1.ListOptions{})
		if err != nil {
			framework.Failf("Failed to list pods: %v", err)
		}
		if e, a := int(*(rc.Spec.Replicas)), len(pods.Items); e != a {
			framework.Failf("expect %d pods, got %d pods", e, a)
		}
		gatherMetrics(f)
	})

	// deleteOptions.OrphanDependents is deprecated in 1.7 and preferred to use the PropagationPolicy.
	// Discussion is tracked under https://github.com/kubernetes/kubernetes/issues/65427 to promote for conformance in future.
	It("should orphan pods created by rc if deleteOptions.OrphanDependents is nil", func() {
		clientSet := f.ClientSet
		rcClient := clientSet.CoreV1().ReplicationControllers(f.Namespace.Name)
		podClient := clientSet.CoreV1().Pods(f.Namespace.Name)
		rcName := "simpletest.rc"
		uniqLabels := getUniqLabel("gctest", "orphan_pods_nil_option")
		rc := newOwnerRC(f, rcName, 2, uniqLabels)
		By("create the rc")
		rc, err := rcClient.Create(rc)
		if err != nil {
			framework.Failf("Failed to create replication controller: %v", err)
		}
		// wait for rc to create some pods
		if err := wait.Poll(5*time.Second, 30*time.Second, func() (bool, error) {
			rc, err := rcClient.Get(rc.Name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to get rc: %v", err)
			}
			if rc.Status.Replicas == *rc.Spec.Replicas {
				return true, nil
			} else {
				return false, nil
			}
		}); err != nil {
			framework.Failf("failed to wait for the rc.Status.Replicas to reach rc.Spec.Replicas: %v", err)
		}
		By("delete the rc")
		deleteOptions := &metav1.DeleteOptions{}
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(rc.UID))
		if err := rcClient.Delete(rc.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the rc: %v", err)
		}
		By("wait for 30 seconds to see if the garbage collector mistakenly deletes the pods")
		time.Sleep(30 * time.Second)
		pods, err := podClient.List(metav1.ListOptions{})
		if err != nil {
			framework.Failf("Failed to list pods: %v", err)
		}
		if e, a := int(*(rc.Spec.Replicas)), len(pods.Items); e != a {
			framework.Failf("expect %d pods, got %d pods", e, a)
		}
		gatherMetrics(f)
	})

	/*
		Release : v1.9
		Testname: Garbage Collector, delete deployment,  propagation policy background
		Description: Create a deployment with a replicaset. Once replicaset is created , delete the deployment  with deleteOptions.PropagationPolicy set to Background. Deleting the deployment MUST delete the replicaset created by the deployment and also the Pods that belong to the deployments MUST be deleted.
	*/
	framework.ConformanceIt("should delete RS created by deployment when not orphaning", func() {
		clientSet := f.ClientSet
		deployClient := clientSet.ExtensionsV1beta1().Deployments(f.Namespace.Name)
		rsClient := clientSet.ExtensionsV1beta1().ReplicaSets(f.Namespace.Name)
		deploymentName := "simpletest.deployment"
		uniqLabels := getUniqLabel("gctest", "delete_rs")
		deployment := newOwnerDeployment(f, deploymentName, uniqLabels)
		By("create the deployment")
		createdDeployment, err := deployClient.Create(deployment)
		if err != nil {
			framework.Failf("Failed to create deployment: %v", err)
		}
		// wait for deployment to create some rs
		By("Wait for the Deployment to create new ReplicaSet")
		err = wait.PollImmediate(500*time.Millisecond, 1*time.Minute, func() (bool, error) {
			rsList, err := rsClient.List(metav1.ListOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to list rs: %v", err)
			}
			return len(rsList.Items) > 0, nil

		})
		if err != nil {
			framework.Failf("Failed to wait for the Deployment to create some ReplicaSet: %v", err)
		}

		By("delete the deployment")
		deleteOptions := getBackgroundOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(createdDeployment.UID))
		if err := deployClient.Delete(deployment.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the deployment: %v", err)
		}
		By("wait for all rs to be garbage collected")
		err = wait.PollImmediate(500*time.Millisecond, 1*time.Minute, func() (bool, error) {
			return verifyRemainingDeploymentsReplicaSetsPods(f, clientSet, deployment, 0, 0, 0)
		})
		if err != nil {
			errList := make([]error, 0)
			errList = append(errList, err)
			remainingRSs, err := rsClient.List(metav1.ListOptions{})
			if err != nil {
				errList = append(errList, fmt.Errorf("failed to list RSs post mortem: %v", err))
			} else {
				errList = append(errList, fmt.Errorf("remaining rs are: %#v", remainingRSs))
			}
			aggregatedError := utilerrors.NewAggregate(errList)
			framework.Failf("Failed to wait for all rs to be garbage collected: %v", aggregatedError)

		}

		gatherMetrics(f)
	})

	/*
		Release : v1.9
		Testname: Garbage Collector, delete deployment, propagation policy orphan
		Description: Create a deployment with a replicaset. Once replicaset is created , delete the deployment  with deleteOptions.PropagationPolicy set to Orphan. Deleting the deployment MUST cause the replicaset created by the deployment to be orphaned, also the Pods created by the deployments MUST be orphaned.
	*/
	framework.ConformanceIt("should orphan RS created by deployment when deleteOptions.PropagationPolicy is Orphan", func() {
		clientSet := f.ClientSet
		deployClient := clientSet.ExtensionsV1beta1().Deployments(f.Namespace.Name)
		rsClient := clientSet.ExtensionsV1beta1().ReplicaSets(f.Namespace.Name)
		deploymentName := "simpletest.deployment"
		uniqLabels := getUniqLabel("gctest", "orphan_rs")
		deployment := newOwnerDeployment(f, deploymentName, uniqLabels)
		By("create the deployment")
		createdDeployment, err := deployClient.Create(deployment)
		if err != nil {
			framework.Failf("Failed to create deployment: %v", err)
		}
		// wait for deployment to create some rs
		By("Wait for the Deployment to create new ReplicaSet")
		err = wait.PollImmediate(500*time.Millisecond, 1*time.Minute, func() (bool, error) {
			rsList, err := rsClient.List(metav1.ListOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to list rs: %v", err)
			}
			return len(rsList.Items) > 0, nil

		})
		if err != nil {
			framework.Failf("Failed to wait for the Deployment to create some ReplicaSet: %v", err)
		}

		By("delete the deployment")
		deleteOptions := getOrphanOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(createdDeployment.UID))
		if err := deployClient.Delete(deployment.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the deployment: %v", err)
		}
		By("wait for 30 seconds to see if the garbage collector mistakenly deletes the rs")
		time.Sleep(30 * time.Second)
		ok, err := verifyRemainingDeploymentsReplicaSetsPods(f, clientSet, deployment, 0, 1, 2)
		if err != nil {
			framework.Failf("Unexpected error while verifying remaining deployments, rs, and pods: %v", err)
		}
		if !ok {
			errList := make([]error, 0)
			remainingRSs, err := rsClient.List(metav1.ListOptions{})
			if err != nil {
				errList = append(errList, fmt.Errorf("failed to list RSs post mortem: %v", err))
			} else {
				errList = append(errList, fmt.Errorf("remaining rs post mortem: %#v", remainingRSs))
			}
			remainingDSs, err := deployClient.List(metav1.ListOptions{})
			if err != nil {
				errList = append(errList, fmt.Errorf("failed to list Deployments post mortem: %v", err))
			} else {
				errList = append(errList, fmt.Errorf("remaining deployment's post mortem: %#v", remainingDSs))
			}
			aggregatedError := utilerrors.NewAggregate(errList)
			framework.Failf("Failed to verify remaining deployments, rs, and pods: %v", aggregatedError)
		}
		rs, err := clientSet.ExtensionsV1beta1().ReplicaSets(f.Namespace.Name).List(metav1.ListOptions{})
		if err != nil {
			framework.Failf("Failed to list ReplicaSet %v", err)
		}
		for _, replicaSet := range rs.Items {
			if metav1.GetControllerOf(&replicaSet.ObjectMeta) != nil {
				framework.Failf("Found ReplicaSet with non nil ownerRef %v", replicaSet)
			}
		}

		gatherMetrics(f)
	})

	/*
		Release : v1.9
		Testname: Garbage Collector, delete replication controller, after owned pods
		Description: Create a replication controller with maximum allocatable Pods between 10 and 100 replicas. Once RC is created and the all Pods are created, delete RC with deleteOptions.PropagationPolicy set to Foreground. Deleting the Replication Controller MUST cause pods created by that RC to be deleted before the RC is deleted.
	*/
	framework.ConformanceIt("should keep the rc around until all its pods are deleted if the deleteOptions says so", func() {
		clientSet := f.ClientSet
		rcClient := clientSet.CoreV1().ReplicationControllers(f.Namespace.Name)
		podClient := clientSet.CoreV1().Pods(f.Namespace.Name)
		rcName := "simpletest.rc"
		uniqLabels := getUniqLabel("gctest", "delete_pods_foreground")
		rc := newOwnerRC(f, rcName, estimateMaximumPods(clientSet, 10, 100), uniqLabels)
		By("create the rc")
		rc, err := rcClient.Create(rc)
		if err != nil {
			framework.Failf("Failed to create replication controller: %v", err)
		}
		// wait for rc to create pods
		if err := wait.Poll(5*time.Second, 30*time.Second, func() (bool, error) {
			rc, err := rcClient.Get(rc.Name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to get rc: %v", err)
			}
			if rc.Status.Replicas == *rc.Spec.Replicas {
				return true, nil
			} else {
				return false, nil
			}
		}); err != nil {
			framework.Failf("failed to wait for the rc.Status.Replicas to reach rc.Spec.Replicas: %v", err)
		}
		By("delete the rc")
		deleteOptions := getForegroundOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(rc.UID))
		if err := rcClient.Delete(rc.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the rc: %v", err)
		}
		By("wait for the rc to be deleted")
		// default client QPS is 20, deleting each pod requires 2 requests, so 30s should be enough
		// TODO: 30s is enough assuming immediate processing of dependents following
		// owner deletion, but in practice there can be a long delay between owner
		// deletion and dependent deletion processing. For now, increase the timeout
		// and investigate the processing delay.
		if err := wait.Poll(1*time.Second, 60*time.Second, func() (bool, error) {
			_, err := rcClient.Get(rc.Name, metav1.GetOptions{})
			if err == nil {
				pods, _ := podClient.List(metav1.ListOptions{})
				framework.Logf("%d pods remaining", len(pods.Items))
				count := 0
				for _, pod := range pods.Items {
					if pod.ObjectMeta.DeletionTimestamp == nil {
						count++
					}
				}
				framework.Logf("%d pods has nil DeletionTimestamp", count)
				framework.Logf("")
				return false, nil
			} else {
				if errors.IsNotFound(err) {
					return true, nil
				} else {
					return false, err
				}
			}
		}); err != nil {
			pods, err2 := podClient.List(metav1.ListOptions{})
			if err2 != nil {
				framework.Failf("%v", err2)
			}
			framework.Logf("%d remaining pods are:", len(pods.Items))
			framework.Logf("The ObjectMeta of the remaining pods are:")
			for _, pod := range pods.Items {
				framework.Logf("%#v", pod.ObjectMeta)
			}
			framework.Failf("failed to delete the rc: %v", err)
		}
		// There shouldn't be any pods
		pods, err := podClient.List(metav1.ListOptions{})
		if err != nil {
			framework.Failf("%v", err)
		}
		if len(pods.Items) != 0 {
			framework.Failf("expected no pods, got %#v", pods)
		}
		gatherMetrics(f)
	})

	// TODO: this should be an integration test
	/*
		Release : v1.9
		Testname: Garbage Collector, multiple owners
		Description: Create a replication controller RC1, with maximum allocatable Pods between 10 and 100 replicas. Create second replication controller RC2 and set RC2 as owner for half of those replicas. Once RC1 is created and the all Pods are created, delete RC1 with deleteOptions.PropagationPolicy set to Foreground. Half of the Pods that has RC2 as owner MUST not be deleted but have a deletion timestamp. Deleting the Replication Controller MUST not delete Pods that are owned by multiple replication controllers.
	*/
	framework.ConformanceIt("should not delete dependents that have both valid owner and owner that's waiting for dependents to be deleted", func() {
		clientSet := f.ClientSet
		rcClient := clientSet.CoreV1().ReplicationControllers(f.Namespace.Name)
		podClient := clientSet.CoreV1().Pods(f.Namespace.Name)
		rc1Name := "simpletest-rc-to-be-deleted"
		replicas := int32(estimateMaximumPods(clientSet, 10, 100))
		halfReplicas := int(replicas / 2)
		uniqLabels_deleted := getUniqLabel("gctest_d", "valid_and_pending_owners_d")
		rc1 := newOwnerRC(f, rc1Name, replicas, uniqLabels_deleted)
		By("create the rc1")
		rc1, err := rcClient.Create(rc1)
		if err != nil {
			framework.Failf("Failed to create replication controller: %v", err)
		}
		rc2Name := "simpletest-rc-to-stay"
		uniqLabels_stay := getUniqLabel("gctest_s", "valid_and_pending_owners_s")
		rc2 := newOwnerRC(f, rc2Name, 0, uniqLabels_stay)
		By("create the rc2")
		rc2, err = rcClient.Create(rc2)
		if err != nil {
			framework.Failf("Failed to create replication controller: %v", err)
		}
		// wait for rc1 to be stable
		if err := wait.Poll(5*time.Second, 30*time.Second, func() (bool, error) {
			rc1, err := rcClient.Get(rc1.Name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to get rc: %v", err)
			}
			if rc1.Status.Replicas == *rc1.Spec.Replicas {
				return true, nil
			} else {
				return false, nil
			}
		}); err != nil {
			framework.Failf("failed to wait for the rc.Status.Replicas to reach rc.Spec.Replicas: %v", err)
		}
		By(fmt.Sprintf("set half of pods created by rc %s to have rc %s as owner as well", rc1Name, rc2Name))
		pods, err := podClient.List(metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		patch := fmt.Sprintf(`{"metadata":{"ownerReferences":[{"apiVersion":"v1","kind":"ReplicationController","name":"%s","uid":"%s"}]}}`, rc2.ObjectMeta.Name, rc2.ObjectMeta.UID)
		for i := 0; i < halfReplicas; i++ {
			pod := pods.Items[i]
			_, err := podClient.Patch(pod.Name, types.StrategicMergePatchType, []byte(patch))
			Expect(err).NotTo(HaveOccurred())
		}

		By(fmt.Sprintf("delete the rc %s", rc1Name))
		deleteOptions := getForegroundOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(rc1.UID))
		if err := rcClient.Delete(rc1.ObjectMeta.Name, deleteOptions); err != nil {
			framework.Failf("failed to delete the rc: %v", err)
		}
		By("wait for the rc to be deleted")
		// TODO: shorten the timeout when we make GC's periodic API rediscovery more efficient.
		// Tracked at https://github.com/kubernetes/kubernetes/issues/50046.
		if err := wait.Poll(5*time.Second, 90*time.Second, func() (bool, error) {
			_, err := rcClient.Get(rc1.Name, metav1.GetOptions{})
			if err == nil {
				pods, _ := podClient.List(metav1.ListOptions{})
				framework.Logf("%d pods remaining", len(pods.Items))
				count := 0
				for _, pod := range pods.Items {
					if pod.ObjectMeta.DeletionTimestamp == nil {
						count++
					}
				}
				framework.Logf("%d pods has nil DeletionTimestamp", count)
				framework.Logf("")
				return false, nil
			} else {
				if errors.IsNotFound(err) {
					return true, nil
				} else {
					return false, err
				}
			}
		}); err != nil {
			pods, err2 := podClient.List(metav1.ListOptions{})
			if err2 != nil {
				framework.Failf("%v", err2)
			}
			framework.Logf("%d remaining pods are:", len(pods.Items))
			framework.Logf("ObjectMeta of remaining pods are:")
			for _, pod := range pods.Items {
				framework.Logf("%#v", pod.ObjectMeta)
			}
			framework.Failf("failed to delete rc %s, err: %v", rc1Name, err)
		}
		// half of the pods should still exist,
		pods, err = podClient.List(metav1.ListOptions{})
		if err != nil {
			framework.Failf("%v", err)
		}
		if len(pods.Items) != halfReplicas {
			framework.Failf("expected %d pods, got %d", halfReplicas, len(pods.Items))
		}
		for _, pod := range pods.Items {
			if pod.ObjectMeta.DeletionTimestamp != nil {
				framework.Failf("expected pod DeletionTimestamp to be nil, got %#v", pod.ObjectMeta)
			}
			// they should only have 1 ownerReference left
			if len(pod.ObjectMeta.OwnerReferences) != 1 {
				framework.Failf("expected pod to only have 1 owner, got %#v", pod.ObjectMeta.OwnerReferences)
			}
		}
		gatherMetrics(f)
	})

	// TODO: should be an integration test
	/*
		Release : v1.9
		Testname: Garbage Collector, dependency cycle
		Description: Create three pods, patch them with Owner references such that pod1 has pod3, pod2 has pod1 and pod3 has pod2 as owner references respectively. Delete pod1 MUST delete all pods. The dependency cycle MUST not block the garbage collection.
	*/
	framework.ConformanceIt("should not be blocked by dependency circle", func() {
		clientSet := f.ClientSet
		podClient := clientSet.CoreV1().Pods(f.Namespace.Name)
		pod1 := newGCPod("pod1")
		pod1, err := podClient.Create(pod1)
		Expect(err).NotTo(HaveOccurred())
		pod2 := newGCPod("pod2")
		pod2, err = podClient.Create(pod2)
		Expect(err).NotTo(HaveOccurred())
		pod3 := newGCPod("pod3")
		pod3, err = podClient.Create(pod3)
		Expect(err).NotTo(HaveOccurred())
		// create circular dependency
		addRefPatch := func(name string, uid types.UID) []byte {
			return []byte(fmt.Sprintf(`{"metadata":{"ownerReferences":[{"apiVersion":"v1","kind":"Pod","name":"%s","uid":"%s","controller":true,"blockOwnerDeletion":true}]}}`, name, uid))
		}
		pod1, err = podClient.Patch(pod1.Name, types.StrategicMergePatchType, addRefPatch(pod3.Name, pod3.UID))
		Expect(err).NotTo(HaveOccurred())
		framework.Logf("pod1.ObjectMeta.OwnerReferences=%#v", pod1.ObjectMeta.OwnerReferences)
		pod2, err = podClient.Patch(pod2.Name, types.StrategicMergePatchType, addRefPatch(pod1.Name, pod1.UID))
		Expect(err).NotTo(HaveOccurred())
		framework.Logf("pod2.ObjectMeta.OwnerReferences=%#v", pod2.ObjectMeta.OwnerReferences)
		pod3, err = podClient.Patch(pod3.Name, types.StrategicMergePatchType, addRefPatch(pod2.Name, pod2.UID))
		Expect(err).NotTo(HaveOccurred())
		framework.Logf("pod3.ObjectMeta.OwnerReferences=%#v", pod3.ObjectMeta.OwnerReferences)
		// delete one pod, should result in the deletion of all pods
		deleteOptions := getForegroundOptions()
		deleteOptions.Preconditions = metav1.NewUIDPreconditions(string(pod1.UID))
		err = podClient.Delete(pod1.ObjectMeta.Name, deleteOptions)
		Expect(err).NotTo(HaveOccurred())
		var pods *v1.PodList
		var err2 error
		// TODO: shorten the timeout when we make GC's periodic API rediscovery more efficient.
		// Tracked at https://github.com/kubernetes/kubernetes/issues/50046.
		if err := wait.Poll(5*time.Second, 90*time.Second, func() (bool, error) {
			pods, err2 = podClient.List(metav1.ListOptions{})
			if err2 != nil {
				return false, fmt.Errorf("Failed to list pods: %v", err)
			}
			if len(pods.Items) == 0 {
				return true, nil
			} else {
				return false, nil
			}
		}); err != nil {
			framework.Logf("pods are %#v", pods.Items)
			framework.Failf("failed to wait for all pods to be deleted: %v", err)
		}
	})

	It("should support cascading deletion of custom resources", func() {
		config, err := framework.LoadConfig()
		if err != nil {
			framework.Failf("failed to load config: %v", err)
		}

		apiExtensionClient, err := apiextensionsclientset.NewForConfig(config)
		if err != nil {
			framework.Failf("failed to initialize apiExtensionClient: %v", err)
		}

		// Create a random custom resource definition and ensure it's available for
		// use.
		definition := apiextensionstestserver.NewRandomNameCustomResourceDefinition(apiextensionsv1beta1.ClusterScoped)
		defer func() {
			err = apiextensionstestserver.DeleteCustomResourceDefinition(definition, apiExtensionClient)
			if err != nil && !errors.IsNotFound(err) {
				framework.Failf("failed to delete CustomResourceDefinition: %v", err)
			}
		}()
		definition, err = apiextensionstestserver.CreateNewCustomResourceDefinition(definition, apiExtensionClient, f.DynamicClient)
		if err != nil {
			framework.Failf("failed to create CustomResourceDefinition: %v", err)
		}

		// Get a client for the custom resource.
		gvr := schema.GroupVersionResource{Group: definition.Spec.Group, Version: definition.Spec.Version, Resource: definition.Spec.Names.Plural}
		resourceClient := f.DynamicClient.Resource(gvr)

		apiVersion := definition.Spec.Group + "/" + definition.Spec.Version

		// Create a custom owner resource.
		ownerName := names.SimpleNameGenerator.GenerateName("owner")
		owner := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": apiVersion,
				"kind":       definition.Spec.Names.Kind,
				"metadata": map[string]interface{}{
					"name": ownerName,
				},
			},
		}
		persistedOwner, err := resourceClient.Create(owner, metav1.CreateOptions{})
		if err != nil {
			framework.Failf("failed to create owner resource %q: %v", ownerName, err)
		}
		framework.Logf("created owner resource %q", ownerName)

		// Create a custom dependent resource.
		dependentName := names.SimpleNameGenerator.GenerateName("dependent")
		dependent := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": apiVersion,
				"kind":       definition.Spec.Names.Kind,
				"metadata": map[string]interface{}{
					"name": dependentName,
					"ownerReferences": []interface{}{
						map[string]interface{}{
							"uid":        string(persistedOwner.GetUID()),
							"apiVersion": apiVersion,
							"kind":       definition.Spec.Names.Kind,
							"name":       ownerName,
						},
					},
				},
			},
		}
		persistedDependent, err := resourceClient.Create(dependent, metav1.CreateOptions{})
		if err != nil {
			framework.Failf("failed to create dependent resource %q: %v", dependentName, err)
		}
		framework.Logf("created dependent resource %q", dependentName)

		// Delete the owner.
		background := metav1.DeletePropagationBackground
		err = resourceClient.Delete(ownerName, &metav1.DeleteOptions{PropagationPolicy: &background})
		if err != nil {
			framework.Failf("failed to delete owner resource %q: %v", ownerName, err)
		}

		// Ensure the dependent is deleted.
		if err := wait.Poll(5*time.Second, 60*time.Second, func() (bool, error) {
			_, err := resourceClient.Get(dependentName, metav1.GetOptions{})
			return errors.IsNotFound(err), nil
		}); err != nil {
			framework.Logf("owner: %#v", persistedOwner)
			framework.Logf("dependent: %#v", persistedDependent)
			framework.Failf("failed waiting for dependent resource %q to be deleted", dependentName)
		}

		// Ensure the owner is deleted.
		_, err = resourceClient.Get(ownerName, metav1.GetOptions{})
		if err == nil {
			framework.Failf("expected owner resource %q to be deleted", ownerName)
		} else {
			if !errors.IsNotFound(err) {
				framework.Failf("unexpected error getting owner resource %q: %v", ownerName, err)
			}
		}
	})

	It("should support orphan deletion of custom resources", func() {
		config, err := framework.LoadConfig()
		if err != nil {
			framework.Failf("failed to load config: %v", err)
		}

		apiExtensionClient, err := apiextensionsclientset.NewForConfig(config)
		if err != nil {
			framework.Failf("failed to initialize apiExtensionClient: %v", err)
		}

		// Create a random custom resource definition and ensure it's available for
		// use.
		definition := apiextensionstestserver.NewRandomNameCustomResourceDefinition(apiextensionsv1beta1.ClusterScoped)
		defer func() {
			err = apiextensionstestserver.DeleteCustomResourceDefinition(definition, apiExtensionClient)
			if err != nil && !errors.IsNotFound(err) {
				framework.Failf("failed to delete CustomResourceDefinition: %v", err)
			}
		}()
		definition, err = apiextensionstestserver.CreateNewCustomResourceDefinition(definition, apiExtensionClient, f.DynamicClient)
		if err != nil {
			framework.Failf("failed to create CustomResourceDefinition: %v", err)
		}

		// Get a client for the custom resource.
		gvr := schema.GroupVersionResource{Group: definition.Spec.Group, Version: definition.Spec.Version, Resource: definition.Spec.Names.Plural}
		resourceClient := f.DynamicClient.Resource(gvr)

		apiVersion := definition.Spec.Group + "/" + definition.Spec.Version

		// Create a custom owner resource.
		ownerName := names.SimpleNameGenerator.GenerateName("owner")
		owner := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": apiVersion,
				"kind":       definition.Spec.Names.Kind,
				"metadata": map[string]interface{}{
					"name": ownerName,
				},
			},
		}
		persistedOwner, err := resourceClient.Create(owner, metav1.CreateOptions{})
		if err != nil {
			framework.Failf("failed to create owner resource %q: %v", ownerName, err)
		}
		framework.Logf("created owner resource %q", ownerName)

		// Create a custom dependent resource.
		dependentName := names.SimpleNameGenerator.GenerateName("dependent")
		dependent := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": apiVersion,
				"kind":       definition.Spec.Names.Kind,
				"metadata": map[string]interface{}{
					"name": dependentName,
					"ownerReferences": []map[string]string{
						{
							"uid":        string(persistedOwner.GetUID()),
							"apiVersion": apiVersion,
							"kind":       definition.Spec.Names.Kind,
							"name":       ownerName,
						},
					},
				},
			},
		}
		_, err = resourceClient.Create(dependent, metav1.CreateOptions{})
		if err != nil {
			framework.Failf("failed to create dependent resource %q: %v", dependentName, err)
		}
		framework.Logf("created dependent resource %q", dependentName)

		// Delete the owner and orphan the dependent.
		err = resourceClient.Delete(ownerName, getOrphanOptions())
		if err != nil {
			framework.Failf("failed to delete owner resource %q: %v", ownerName, err)
		}

		By("wait for the owner to be deleted")
		if err := wait.Poll(5*time.Second, 120*time.Second, func() (bool, error) {
			_, err = resourceClient.Get(ownerName, metav1.GetOptions{})
			if err == nil {
				return false, nil
			}
			if err != nil && !errors.IsNotFound(err) {
				return false, fmt.Errorf("Failed to get owner: %v", err)
			}
			return true, nil
		}); err != nil {
			framework.Failf("timeout in waiting for the owner to be deleted: %v", err)
		}

		// Wait 30s and ensure the dependent is not deleted.
		By("wait for 30 seconds to see if the garbage collector mistakenly deletes the dependent crd")
		if err := wait.Poll(5*time.Second, 30*time.Second, func() (bool, error) {
			_, err := resourceClient.Get(dependentName, metav1.GetOptions{})
			return false, err
		}); err != nil && err != wait.ErrWaitTimeout {
			framework.Failf("failed to ensure the dependent is not deleted: %v", err)
		}
	})

	It("should delete jobs and pods created by cronjob", func() {
		framework.SkipIfMissingResource(f.DynamicClient, CronJobGroupVersionResource, f.Namespace.Name)

		By("Create the cronjob")
		cronJob := newCronJob("simple", "*/1 * * * ?")
		cronJob, err := f.ClientSet.BatchV1beta1().CronJobs(f.Namespace.Name).Create(cronJob)
		Expect(err).NotTo(HaveOccurred())

		By("Wait for the CronJob to create new Job")
		err = wait.PollImmediate(500*time.Millisecond, 2*time.Minute, func() (bool, error) {
			jobs, err := f.ClientSet.BatchV1().Jobs(f.Namespace.Name).List(metav1.ListOptions{})
			if err != nil {
				return false, fmt.Errorf("Failed to list jobs: %v", err)
			}
			return len(jobs.Items) > 0, nil
		})
		if err != nil {
			framework.Failf("Failed to wait for the CronJob to create some Jobs: %v", err)
		}

		By("Delete the cronjob")
		if err := f.ClientSet.BatchV1beta1().CronJobs(f.Namespace.Name).Delete(cronJob.Name, getBackgroundOptions()); err != nil {
			framework.Failf("Failed to delete the CronJob: %v", err)
		}
		By("Verify if cronjob does not leave jobs nor pods behind")
		err = wait.PollImmediate(500*time.Millisecond, 1*time.Minute, func() (bool, error) {
			return verifyRemainingCronJobsJobsPods(f, f.ClientSet, 0, 0, 0)
		})
		if err != nil {
			framework.Failf("Failed to wait for all jobs and pods to be deleted: %v", err)
		}

		gatherMetrics(f)
	})
})
