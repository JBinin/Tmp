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

package markmaster

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	"k8s.io/kubernetes/pkg/util/node"
)

func TestMarkMaster(t *testing.T) {
	// Note: this test takes advantage of the deterministic marshalling of
	// JSON provided by strategicpatch so that "expectedPatch" can use a
	// string equality test instead of a logical JSON equality test. That
	// will need to change if strategicpatch's behavior changes in the
	// future.
	tests := []struct {
		name           string
		existingLabel  string
		existingTaints []v1.Taint
		newTaints      []v1.Taint
		expectedPatch  string
	}{
		{
			"master label and taint missing",
			"",
			nil,
			[]v1.Taint{kubeadmconstants.MasterTaint},
			"{\"metadata\":{\"labels\":{\"node-role.kubernetes.io/master\":\"\"}},\"spec\":{\"taints\":[{\"effect\":\"NoSchedule\",\"key\":\"node-role.kubernetes.io/master\"}]}}",
		},
		{
			"master label and taint missing but taint not wanted",
			"",
			nil,
			nil,
			"{\"metadata\":{\"labels\":{\"node-role.kubernetes.io/master\":\"\"}}}",
		},
		{
			"master label missing",
			"",
			[]v1.Taint{kubeadmconstants.MasterTaint},
			[]v1.Taint{kubeadmconstants.MasterTaint},
			"{\"metadata\":{\"labels\":{\"node-role.kubernetes.io/master\":\"\"}}}",
		},
		{
			"master taint missing",
			kubeadmconstants.LabelNodeRoleMaster,
			nil,
			[]v1.Taint{kubeadmconstants.MasterTaint},
			"{\"spec\":{\"taints\":[{\"effect\":\"NoSchedule\",\"key\":\"node-role.kubernetes.io/master\"}]}}",
		},
		{
			"nothing missing",
			kubeadmconstants.LabelNodeRoleMaster,
			[]v1.Taint{kubeadmconstants.MasterTaint},
			[]v1.Taint{kubeadmconstants.MasterTaint},
			"{}",
		},
		{
			"has taint and no new taints wanted",
			kubeadmconstants.LabelNodeRoleMaster,
			[]v1.Taint{
				{
					Key:    "node.cloudprovider.kubernetes.io/uninitialized",
					Effect: v1.TaintEffectNoSchedule,
				},
			},
			nil,
			"{}",
		},
		{
			"has taint and should merge with wanted taint",
			kubeadmconstants.LabelNodeRoleMaster,
			[]v1.Taint{
				{
					Key:    "node.cloudprovider.kubernetes.io/uninitialized",
					Effect: v1.TaintEffectNoSchedule,
				},
			},
			[]v1.Taint{kubeadmconstants.MasterTaint},
			"{\"spec\":{\"taints\":[{\"effect\":\"NoSchedule\",\"key\":\"node-role.kubernetes.io/master\"},{\"effect\":\"NoSchedule\",\"key\":\"node.cloudprovider.kubernetes.io/uninitialized\"}]}}",
		},
	}

	for _, tc := range tests {
		hostname, err := node.GetHostname("")
		if err != nil {
			t.Fatalf("MarkMaster(%s): unexpected error: %v", tc.name, err)
		}
		masterNode := &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: hostname,
				Labels: map[string]string{
					kubeletapis.LabelHostname: hostname,
				},
			},
		}

		if tc.existingLabel != "" {
			masterNode.ObjectMeta.Labels[tc.existingLabel] = ""
		}

		if tc.existingTaints != nil {
			masterNode.Spec.Taints = tc.existingTaints
		}

		jsonNode, err := json.Marshal(masterNode)
		if err != nil {
			t.Fatalf("MarkMaster(%s): unexpected encoding error: %v", tc.name, err)
		}

		var patchRequest string
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if req.URL.Path != "/api/v1/nodes/"+hostname {
				t.Errorf("MarkMaster(%s): request for unexpected HTTP resource: %v", tc.name, req.URL.Path)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			switch req.Method {
			case "GET":
			case "PATCH":
				patchRequest = toString(req.Body)
			default:
				t.Errorf("MarkMaster(%s): request for unexpected HTTP verb: %v", tc.name, req.Method)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonNode)
		}))
		defer s.Close()

		cs, err := clientset.NewForConfig(&restclient.Config{Host: s.URL})
		if err != nil {
			t.Fatalf("MarkMaster(%s): unexpected error building clientset: %v", tc.name, err)
		}

		if err := MarkMaster(cs, hostname, tc.newTaints); err != nil {
			t.Errorf("MarkMaster(%s) returned unexpected error: %v", tc.name, err)
		}

		if tc.expectedPatch != patchRequest {
			t.Errorf("MarkMaster(%s) wanted patch %v, got %v", tc.name, tc.expectedPatch, patchRequest)
		}
	}
}

func toString(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}
