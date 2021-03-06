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
Copyright 2014 The Kubernetes Authors.

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

package node

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"

	. "github.com/onsi/ginkgo"
)

// partially cloned from webserver.go
type State struct {
	Received map[string]int
}

func testPreStop(c clientset.Interface, ns string) {
	// This is the server that will receive the preStop notification
	podDescr := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "server",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "server",
					Image: imageutils.GetE2EImage(imageutils.Nettest),
					Ports: []v1.ContainerPort{{ContainerPort: 8080}},
				},
			},
		},
	}
	By(fmt.Sprintf("Creating server pod %s in namespace %s", podDescr.Name, ns))
	podDescr, err := c.CoreV1().Pods(ns).Create(podDescr)
	framework.ExpectNoError(err, fmt.Sprintf("creating pod %s", podDescr.Name))

	// At the end of the test, clean up by removing the pod.
	defer func() {
		By("Deleting the server pod")
		c.CoreV1().Pods(ns).Delete(podDescr.Name, nil)
	}()

	By("Waiting for pods to come up.")
	err = framework.WaitForPodRunningInNamespace(c, podDescr)
	framework.ExpectNoError(err, "waiting for server pod to start")

	val := "{\"Source\": \"prestop\"}"

	podOut, err := c.CoreV1().Pods(ns).Get(podDescr.Name, metav1.GetOptions{})
	framework.ExpectNoError(err, "getting pod info")

	preStopDescr := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tester",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "tester",
					Image:   imageutils.GetE2EImage(imageutils.BusyBox),
					Command: []string{"sleep", "600"},
					Lifecycle: &v1.Lifecycle{
						PreStop: &v1.Handler{
							Exec: &v1.ExecAction{
								Command: []string{
									"wget", "-O-", "--post-data=" + val, fmt.Sprintf("http://%s:8080/write", podOut.Status.PodIP),
								},
							},
						},
					},
				},
			},
		},
	}

	By(fmt.Sprintf("Creating tester pod %s in namespace %s", preStopDescr.Name, ns))
	preStopDescr, err = c.CoreV1().Pods(ns).Create(preStopDescr)
	framework.ExpectNoError(err, fmt.Sprintf("creating pod %s", preStopDescr.Name))
	deletePreStop := true

	// At the end of the test, clean up by removing the pod.
	defer func() {
		if deletePreStop {
			By("Deleting the tester pod")
			c.CoreV1().Pods(ns).Delete(preStopDescr.Name, nil)
		}
	}()

	err = framework.WaitForPodRunningInNamespace(c, preStopDescr)
	framework.ExpectNoError(err, "waiting for tester pod to start")

	// Delete the pod with the preStop handler.
	By("Deleting pre-stop pod")
	if err := c.CoreV1().Pods(ns).Delete(preStopDescr.Name, nil); err == nil {
		deletePreStop = false
	}
	framework.ExpectNoError(err, fmt.Sprintf("deleting pod: %s", preStopDescr.Name))

	// Validate that the server received the web poke.
	err = wait.Poll(time.Second*5, time.Second*60, func() (bool, error) {

		ctx, cancel := context.WithTimeout(context.Background(), framework.SingleCallTimeout)
		defer cancel()

		var body []byte
		body, err = c.CoreV1().RESTClient().Get().
			Context(ctx).
			Namespace(ns).
			Resource("pods").
			SubResource("proxy").
			Name(podDescr.Name).
			Suffix("read").
			DoRaw()

		if err != nil {
			if ctx.Err() != nil {
				framework.Failf("Error validating prestop: %v", err)
				return true, err
			}
			By(fmt.Sprintf("Error validating prestop: %v", err))
		} else {
			framework.Logf("Saw: %s", string(body))
			state := State{}
			err := json.Unmarshal(body, &state)
			if err != nil {
				framework.Logf("Error parsing: %v", err)
				return false, nil
			}
			if state.Received["prestop"] != 0 {
				return true, nil
			}
		}
		return false, nil
	})
	framework.ExpectNoError(err, "validating pre-stop.")
}

var _ = SIGDescribe("PreStop", func() {
	f := framework.NewDefaultFramework("prestop")

	/*
		Release : v1.9
		Testname: Pods, prestop hook
		Description: Create a server pod with a rest endpoint '/write' that changes state.Received field. Create a Pod with a pre-stop handle that posts to the /write endpoint on the server Pod. Verify that the Pod with pre-stop hook is running. Delete the Pod with the pre-stop hook. Before the Pod is deleted, pre-stop handler MUST be called when configured. Verify that the Pod is deleted and a call to prestop hook is verified by checking the status received on the server Pod.
	*/
	framework.ConformanceIt("should call prestop when killing a pod ", func() {
		testPreStop(f.ClientSet, f.Namespace.Name)
	})
})
