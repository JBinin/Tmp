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

package kubelet

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig"
	"k8s.io/kubernetes/pkg/util/version"
)

func TestCreateConfigMap(t *testing.T) {
	nodeName := "fake-node"
	client := fake.NewSimpleClientset()
	cfg := &kubeadmapi.InitConfiguration{
		NodeRegistration:  kubeadmapi.NodeRegistrationOptions{Name: nodeName},
		KubernetesVersion: "v1.11.0",
		ComponentConfigs: kubeadmapi.ComponentConfigs{
			Kubelet: &kubeletconfig.KubeletConfiguration{},
		},
	}

	client.PrependReactor("get", "nodes", func(action core.Action) (bool, runtime.Object, error) {
		return true, &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: nodeName,
			},
			Spec: v1.NodeSpec{},
		}, nil
	})
	client.PrependReactor("create", "roles", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "rolebindings", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})

	if err := CreateConfigMap(cfg, client); err != nil {
		t.Errorf("CreateConfigMap: unexpected error %v", err)
	}
}

func TestCreateConfigMapRBACRules(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("create", "roles", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "rolebindings", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})

	if err := createConfigMapRBACRules(client, version.MustParseSemantic("v1.11.0")); err != nil {
		t.Errorf("createConfigMapRBACRules: unexpected error %v", err)
	}
}
