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

package e2e_node

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = framework.KubeDescribe("Docker features [Feature:Docker][Legacy:Docker]", func() {
	f := framework.NewDefaultFramework("docker-feature-test")

	BeforeEach(func() {
		framework.RunIfContainerRuntimeIs("docker")
	})

	Context("when live-restore is enabled [Serial] [Slow] [Disruptive]", func() {
		It("containers should not be disrupted when the daemon shuts down and restarts", func() {
			const (
				podName       = "live-restore-test-pod"
				containerName = "live-restore-test-container"
			)

			isSupported, err := isDockerLiveRestoreSupported()
			framework.ExpectNoError(err)
			if !isSupported {
				framework.Skipf("Docker live-restore is not supported.")
			}
			isEnabled, err := isDockerLiveRestoreEnabled()
			framework.ExpectNoError(err)
			if !isEnabled {
				framework.Skipf("Docker live-restore is not enabled.")
			}

			By("Create the test pod.")
			pod := f.PodClient().CreateSync(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: podName},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:  containerName,
						Image: imageutils.GetE2EImage(imageutils.Nginx),
					}},
				},
			})

			By("Ensure that the container is running before Docker is down.")
			Eventually(func() bool {
				return isContainerRunning(pod.Status.PodIP)
			}).Should(BeTrue())

			startTime1, err := getContainerStartTime(f, podName, containerName)
			framework.ExpectNoError(err)

			By("Stop Docker daemon.")
			framework.ExpectNoError(stopDockerDaemon())
			isDockerDown := true
			defer func() {
				if isDockerDown {
					By("Start Docker daemon.")
					framework.ExpectNoError(startDockerDaemon())
				}
			}()

			By("Ensure that the container is running after Docker is down.")
			Consistently(func() bool {
				return isContainerRunning(pod.Status.PodIP)
			}).Should(BeTrue())

			By("Start Docker daemon.")
			framework.ExpectNoError(startDockerDaemon())
			isDockerDown = false

			By("Ensure that the container is running after Docker has restarted.")
			Consistently(func() bool {
				return isContainerRunning(pod.Status.PodIP)
			}).Should(BeTrue())

			By("Ensure that the container has not been restarted after Docker is restarted.")
			Consistently(func() bool {
				startTime2, err := getContainerStartTime(f, podName, containerName)
				framework.ExpectNoError(err)
				return startTime1 == startTime2
			}, 3*time.Second, time.Second).Should(BeTrue())
		})
	})
})

// isContainerRunning returns true if the container is running by checking
// whether the server is responding, and false otherwise.
func isContainerRunning(podIP string) bool {
	output, err := runCommand("curl", podIP)
	if err != nil {
		return false
	}
	return strings.Contains(output, "Welcome to nginx!")
}

// getContainerStartTime returns the start time of the container with the
// containerName of the pod having the podName.
func getContainerStartTime(f *framework.Framework, podName, containerName string) (time.Time, error) {
	pod, err := f.PodClient().Get(podName, metav1.GetOptions{})
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get pod %q: %v", podName, err)
	}
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name != containerName {
			continue
		}
		if status.State.Running == nil {
			return time.Time{}, fmt.Errorf("%v/%v is not running", podName, containerName)
		}
		return status.State.Running.StartedAt.Time, nil
	}
	return time.Time{}, fmt.Errorf("failed to find %v/%v", podName, containerName)
}
