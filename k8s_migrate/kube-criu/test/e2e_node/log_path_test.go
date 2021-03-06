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

package e2e_node

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/kubelet"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

const (
	logString = "This is the expected log content of this node e2e test"

	logPodName    = "logger-pod"
	logContName   = "logger-container"
	checkPodName  = "checker-pod"
	checkContName = "checker-container"
)

var _ = framework.KubeDescribe("ContainerLogPath [NodeConformance]", func() {
	f := framework.NewDefaultFramework("kubelet-container-log-path")
	Describe("Pod with a container", func() {
		Context("printed log to stdout", func() {
			BeforeEach(func() {
				if framework.TestContext.ContainerRuntime == "docker" {
					// Container Log Path support requires JSON logging driver.
					// It does not work when Docker daemon is logging to journald.
					d, err := getDockerLoggingDriver()
					framework.ExpectNoError(err)
					if d != "json-file" {
						framework.Skipf("Skipping because Docker daemon is using a logging driver other than \"json-file\": %s", d)
					}
					// Even if JSON logging is in use, this test fails if SELinux support
					// is enabled, since the isolation provided by the SELinux policy
					// prevents processes running inside Docker containers (under SELinux
					// type svirt_lxc_net_t) from accessing the log files which are owned
					// by Docker (and labeled with the container_var_lib_t type.)
					//
					// Therefore, let's also skip this test when running with SELinux
					// support enabled.
					e, err := isDockerSELinuxSupportEnabled()
					framework.ExpectNoError(err)
					if e {
						framework.Skipf("Skipping because Docker daemon is running with SELinux support enabled")
					}
				}
			})
			It("should print log to correct log path", func() {
				podClient := f.PodClient()
				ns := f.Namespace.Name

				logDirVolumeName := "log-dir-vol"
				logDir := kubelet.ContainerLogsDir

				logPod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: logPodName,
					},
					Spec: v1.PodSpec{
						// this pod is expected to exit successfully
						RestartPolicy: v1.RestartPolicyNever,
						Containers: []v1.Container{
							{
								Image:   busyboxImage,
								Name:    logContName,
								Command: []string{"sh", "-c", "echo " + logString},
							},
						},
					},
				}

				podClient.Create(logPod)
				err := framework.WaitForPodSuccessInNamespace(f.ClientSet, logPodName, ns)
				framework.ExpectNoError(err, "Failed waiting for pod: %s to enter success state", logPodName)

				// get containerID from created Pod
				createdLogPod, err := podClient.Get(logPodName, metav1.GetOptions{})
				logConID := kubecontainer.ParseContainerID(createdLogPod.Status.ContainerStatuses[0].ContainerID)
				framework.ExpectNoError(err, "Failed to get pod: %s", logPodName)

				expectedlogFile := logDir + "/" + logPodName + "_" + ns + "_" + logContName + "-" + logConID.ID + ".log"

				hostPathType := new(v1.HostPathType)
				*hostPathType = v1.HostPathType(string(v1.HostPathFileOrCreate))

				checkPod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: checkPodName,
					},
					Spec: v1.PodSpec{
						// this pod is expected to exit successfully
						RestartPolicy: v1.RestartPolicyNever,
						Containers: []v1.Container{
							{
								Image: busyboxImage,
								Name:  checkContName,
								// If we find expected log file and contains right content, exit 0
								// else, keep checking until test timeout
								Command: []string{"sh", "-c", "while true; do if [ -e " + expectedlogFile + " ] && grep -q " + logString + " " + expectedlogFile + "; then exit 0; fi; sleep 1; done"},
								VolumeMounts: []v1.VolumeMount{
									{
										Name: logDirVolumeName,
										// mount ContainerLogsDir to the same path in container
										MountPath: expectedlogFile,
										ReadOnly:  true,
									},
								},
							},
						},
						Volumes: []v1.Volume{
							{
								Name: logDirVolumeName,
								VolumeSource: v1.VolumeSource{
									HostPath: &v1.HostPathVolumeSource{
										Path: expectedlogFile,
										Type: hostPathType,
									},
								},
							},
						},
					},
				}

				podClient.Create(checkPod)
				err = framework.WaitForPodSuccessInNamespace(f.ClientSet, checkPodName, ns)
				framework.ExpectNoError(err, "Failed waiting for pod: %s to enter success state", checkPodName)
			})
		})
	})
})
