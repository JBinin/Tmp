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
Copyright 2018 The Kubernetes Authors.

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

package polymorphichelpers

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func TestMapBasedSelectorForObject(t *testing.T) {
	tests := []struct {
		object         runtime.Object
		expectSelector string
		expectErr      bool
	}{
		{
			object: &api.ReplicationController{
				Spec: api.ReplicationControllerSpec{
					Selector: map[string]string{
						"foo": "bar",
					},
				},
			},
			expectSelector: "foo=bar",
		},
		{
			object:    &api.Pod{},
			expectErr: true,
		},
		{
			object: &api.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "bar",
					},
				},
			},
			expectSelector: "foo=bar",
		},
		{
			object: &api.Service{
				Spec: api.ServiceSpec{
					Selector: map[string]string{
						"foo": "bar",
					},
				},
			},
			expectSelector: "foo=bar",
		},
		{
			object:    &api.Service{},
			expectErr: true,
		},
		{
			object: &extensions.Deployment{
				Spec: extensions.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"foo": "bar",
						},
					},
				},
			},
			expectSelector: "foo=bar",
		},
		{
			object: &extensions.Deployment{
				Spec: extensions.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key: "foo",
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			object: &extensions.ReplicaSet{
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"foo": "bar",
						},
					},
				},
			},
			expectSelector: "foo=bar",
		},
		{
			object: &extensions.ReplicaSet{
				Spec: extensions.ReplicaSetSpec{
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key: "foo",
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			object:    &api.Node{},
			expectErr: true,
		},
	}

	for _, test := range tests {
		actual, err := mapBasedSelectorForObject(test.object)
		if test.expectErr && err == nil {
			t.Error("unexpected non-error")
		}
		if !test.expectErr && err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if actual != test.expectSelector {
			t.Errorf("expected selector %q, but got %q", test.expectSelector, actual)
		}
	}
}
