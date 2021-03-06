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
Copyright 2015 The Kubernetes Authors.

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

package exec

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"
	core "k8s.io/client-go/testing"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
)

// newAllowEscalatingExec returns `admission.Interface` that allows execution on
// "hostIPC", "hostPID" and "privileged".
func newAllowEscalatingExec() *DenyExec {
	return &DenyExec{
		Handler:    admission.NewHandler(admission.Connect),
		hostIPC:    false,
		hostPID:    false,
		privileged: false,
	}
}

func TestAdmission(t *testing.T) {
	privPod := validPod("privileged")
	priv := true
	privPod.Spec.Containers[0].SecurityContext = &api.SecurityContext{
		Privileged: &priv,
	}

	hostPIDPod := validPod("hostPID")
	hostPIDPod.Spec.SecurityContext = &api.PodSecurityContext{}
	hostPIDPod.Spec.SecurityContext.HostPID = true

	hostIPCPod := validPod("hostIPC")
	hostIPCPod.Spec.SecurityContext = &api.PodSecurityContext{}
	hostIPCPod.Spec.SecurityContext.HostIPC = true

	testCases := map[string]struct {
		pod          *api.Pod
		shouldAccept bool
	}{
		"priv": {
			shouldAccept: false,
			pod:          privPod,
		},
		"hostPID": {
			shouldAccept: false,
			pod:          hostPIDPod,
		},
		"hostIPC": {
			shouldAccept: false,
			pod:          hostIPCPod,
		},
		"non privileged": {
			shouldAccept: true,
			pod:          validPod("nonPrivileged"),
		},
	}

	// Get the direct object though to allow testAdmission to inject the client
	handler := NewDenyEscalatingExec()

	for _, tc := range testCases {
		testAdmission(t, tc.pod, handler, tc.shouldAccept)
	}

	// run with a permissive config and all cases should pass
	handler = newAllowEscalatingExec()

	for _, tc := range testCases {
		testAdmission(t, tc.pod, handler, true)
	}

	// run against an init container
	handler = NewDenyEscalatingExec()

	for _, tc := range testCases {
		tc.pod.Spec.InitContainers = tc.pod.Spec.Containers
		tc.pod.Spec.Containers = nil
		testAdmission(t, tc.pod, handler, tc.shouldAccept)
	}

	// run with a permissive config and all cases should pass
	handler = newAllowEscalatingExec()

	for _, tc := range testCases {
		testAdmission(t, tc.pod, handler, true)
	}
}

func testAdmission(t *testing.T, pod *api.Pod, handler *DenyExec, shouldAccept bool) {
	mockClient := &fake.Clientset{}
	mockClient.AddReactor("get", "pods", func(action core.Action) (bool, runtime.Object, error) {
		if action.(core.GetAction).GetName() == pod.Name {
			return true, pod, nil
		}
		t.Errorf("Unexpected API call: %#v", action)
		return true, nil, nil
	})

	handler.SetInternalKubeClientSet(mockClient)
	admission.ValidateInitialization(handler)

	// pods/exec
	{
		err := handler.Validate(admission.NewAttributesRecord(nil, nil, api.Kind("Pod").WithVersion("version"), "test", pod.Name, api.Resource("pods").WithVersion("version"), "exec", admission.Connect, false, nil))
		if shouldAccept && err != nil {
			t.Errorf("Unexpected error returned from admission handler: %v", err)
		}
		if !shouldAccept && err == nil {
			t.Errorf("An error was expected from the admission handler. Received nil")
		}
	}

	// pods/attach
	{
		err := handler.Validate(admission.NewAttributesRecord(nil, nil, api.Kind("Pod").WithVersion("version"), "test", pod.Name, api.Resource("pods").WithVersion("version"), "attach", admission.Connect, false, nil))
		if shouldAccept && err != nil {
			t.Errorf("Unexpected error returned from admission handler: %v", err)
		}
		if !shouldAccept && err == nil {
			t.Errorf("An error was expected from the admission handler. Received nil")
		}
	}
}

// Test to ensure legacy admission controller works as expected.
func TestDenyExecOnPrivileged(t *testing.T) {
	privPod := validPod("privileged")
	priv := true
	privPod.Spec.Containers[0].SecurityContext = &api.SecurityContext{
		Privileged: &priv,
	}

	hostPIDPod := validPod("hostPID")
	hostPIDPod.Spec.SecurityContext = &api.PodSecurityContext{}
	hostPIDPod.Spec.SecurityContext.HostPID = true

	hostIPCPod := validPod("hostIPC")
	hostIPCPod.Spec.SecurityContext = &api.PodSecurityContext{}
	hostIPCPod.Spec.SecurityContext.HostIPC = true

	testCases := map[string]struct {
		pod          *api.Pod
		shouldAccept bool
	}{
		"priv": {
			shouldAccept: false,
			pod:          privPod,
		},
		"hostPID": {
			shouldAccept: true,
			pod:          hostPIDPod,
		},
		"hostIPC": {
			shouldAccept: true,
			pod:          hostIPCPod,
		},
		"non privileged": {
			shouldAccept: true,
			pod:          validPod("nonPrivileged"),
		},
	}

	// Get the direct object though to allow testAdmission to inject the client
	handler := NewDenyExecOnPrivileged()

	for _, tc := range testCases {
		testAdmission(t, tc.pod, handler, tc.shouldAccept)
	}

	// test init containers
	for _, tc := range testCases {
		tc.pod.Spec.InitContainers = tc.pod.Spec.Containers
		tc.pod.Spec.Containers = nil
		testAdmission(t, tc.pod, handler, tc.shouldAccept)
	}
}

func validPod(name string) *api.Pod {
	return &api.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "test"},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{Name: "ctr1", Image: "image"},
				{Name: "ctr2", Image: "image2"},
			},
		},
	}
}
