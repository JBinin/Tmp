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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

var _ = framework.KubeDescribe("Docker Containers", func() {
	f := framework.NewDefaultFramework("containers")

	/*
		Release : v1.9
		Testname: Docker containers, without command and arguments
		Description: Default command and arguments from the docker image entrypoint MUST be used when Pod does not specify the container command
	*/
	framework.ConformanceIt("should use the image defaults if command and args are blank [NodeConformance]", func() {
		f.TestContainerOutput("use defaults", entrypointTestPod(), 0, []string{
			"[/ep default arguments]",
		})
	})

	/*
		Release : v1.9
		Testname: Docker containers, with arguments
		Description: Default command and  from the docker image entrypoint MUST be used when Pod does not specify the container command but the arguments from Pod spec MUST override when specified.
	*/
	framework.ConformanceIt("should be able to override the image's default arguments (docker cmd) [NodeConformance]", func() {
		pod := entrypointTestPod()
		pod.Spec.Containers[0].Args = []string{"override", "arguments"}

		f.TestContainerOutput("override arguments", pod, 0, []string{
			"[/ep override arguments]",
		})
	})

	// Note: when you override the entrypoint, the image's arguments (docker cmd)
	// are ignored.
	/*
		Release : v1.9
		Testname: Docker containers, with command
		Description: Default command from the docker image entrypoint MUST NOT be used when Pod specifies the container command.  Command from Pod spec MUST override the command in the image.
	*/
	framework.ConformanceIt("should be able to override the image's default command (docker entrypoint) [NodeConformance]", func() {
		pod := entrypointTestPod()
		pod.Spec.Containers[0].Command = []string{"/ep-2"}

		f.TestContainerOutput("override command", pod, 0, []string{
			"[/ep-2]",
		})
	})

	/*
		Release : v1.9
		Testname: Docker containers, with command and arguments
		Description: Default command and arguments from the docker image entrypoint MUST NOT be used when Pod specifies the container command and arguments.  Command and arguments from Pod spec MUST override the command and arguments in the image.
	*/
	framework.ConformanceIt("should be able to override the image's default command and arguments [NodeConformance]", func() {
		pod := entrypointTestPod()
		pod.Spec.Containers[0].Command = []string{"/ep-2"}
		pod.Spec.Containers[0].Args = []string{"override", "arguments"}

		f.TestContainerOutput("override all", pod, 0, []string{
			"[/ep-2 override arguments]",
		})
	})
})

const testContainerName = "test-container"

// Return a prototypical entrypoint test pod
func entrypointTestPod() *v1.Pod {
	podName := "client-containers-" + string(uuid.NewUUID())

	one := int64(1)
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  testContainerName,
					Image: imageutils.GetE2EImage(imageutils.EntrypointTester),
				},
			},
			RestartPolicy:                 v1.RestartPolicyNever,
			TerminationGracePeriodSeconds: &one,
		},
	}
}
