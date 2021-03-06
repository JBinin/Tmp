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

package kubectl

import (
	"reflect"
	"testing"

	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterRoleBindingGenerate(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]interface{}
		expected  *rbacv1beta1.ClusterRoleBinding
		expectErr bool
	}{
		{
			name: "valid case 1",
			params: map[string]interface{}{
				"name":           "foo",
				"clusterrole":    "admin",
				"user":           []string{"user"},
				"group":          []string{"group"},
				"serviceaccount": []string{"ns1:name1"},
			},
			expected: &rbacv1beta1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				RoleRef: rbacv1beta1.RoleRef{
					APIGroup: rbacv1beta1.GroupName,
					Kind:     "ClusterRole",
					Name:     "admin",
				},
				Subjects: []rbacv1beta1.Subject{
					{
						APIGroup: rbacv1beta1.GroupName,
						Kind:     rbacv1beta1.UserKind,
						Name:     "user",
					},
					{
						APIGroup: rbacv1beta1.GroupName,
						Kind:     rbacv1beta1.GroupKind,
						Name:     "group",
					},
					{
						Kind:      rbacv1beta1.ServiceAccountKind,
						APIGroup:  "",
						Namespace: "ns1",
						Name:      "name1",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid case 2",
			params: map[string]interface{}{
				"name":           "foo",
				"clusterrole":    "admin",
				"user":           []string{"user1", "user2"},
				"group":          []string{"group1", "group2"},
				"serviceaccount": []string{"ns1:name1", "ns2:name2"},
			},
			expected: &rbacv1beta1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				RoleRef: rbacv1beta1.RoleRef{
					APIGroup: rbacv1beta1.GroupName,
					Kind:     "ClusterRole",
					Name:     "admin",
				},
				Subjects: []rbacv1beta1.Subject{
					{
						APIGroup: rbacv1beta1.GroupName,
						Kind:     rbacv1beta1.UserKind,
						Name:     "user1",
					},
					{
						APIGroup: rbacv1beta1.GroupName,
						Kind:     rbacv1beta1.UserKind,
						Name:     "user2",
					},
					{
						APIGroup: rbacv1beta1.GroupName,
						Kind:     rbacv1beta1.GroupKind,
						Name:     "group1",
					},
					{
						APIGroup: rbacv1beta1.GroupName,
						Kind:     rbacv1beta1.GroupKind,
						Name:     "group2",
					},
					{
						Kind:      rbacv1beta1.ServiceAccountKind,
						APIGroup:  "",
						Namespace: "ns1",
						Name:      "name1",
					},
					{
						Kind:      rbacv1beta1.ServiceAccountKind,
						APIGroup:  "",
						Namespace: "ns2",
						Name:      "name2",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid case 3",
			params: map[string]interface{}{
				"name":        "foo",
				"clusterrole": "admin",
			},
			expected: &rbacv1beta1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				RoleRef: rbacv1beta1.RoleRef{
					APIGroup: rbacv1beta1.GroupName,
					Kind:     "ClusterRole",
					Name:     "admin",
				},
			},
			expectErr: false,
		},
		{
			name: "invalid serviceaccount, expected format: <namespace:name>",
			params: map[string]interface{}{
				"name":           "role",
				"clusterrole":    "admin",
				"user":           []string{"user"},
				"group":          []string{"group"},
				"serviceaccount": []string{"ns1-name1"},
			},
			expectErr: true,
		},
		{
			name: "name must be specified",
			params: map[string]interface{}{
				"name":           "",
				"clusterrole":    "admin",
				"user":           []string{"user"},
				"group":          []string{"group"},
				"serviceaccount": []string{"ns1:name1"},
			},
			expectErr: true,
		},
		{
			name: "clusterrole must be specified",
			params: map[string]interface{}{
				"name":           "foo",
				"clusterrole":    "",
				"user":           []string{"user"},
				"group":          []string{"group"},
				"serviceaccount": []string{"ns1:name1"},
			},
			expectErr: true,
		},
		{
			name: "expected user []string",
			params: map[string]interface{}{
				"name":           "role",
				"clusterrole":    "admin",
				"user":           "user",
				"group":          []string{"group"},
				"serviceaccount": []string{"ns1:name1"},
			},
			expectErr: true,
		},
		{
			name: "expected group []string",
			params: map[string]interface{}{
				"name":           "role",
				"clusterrole":    "admin",
				"user":           []string{"user"},
				"group":          "group",
				"serviceaccount": []string{"ns1:name1"},
			},
			expectErr: true,
		},
		{
			name: "expected serviceaccount []string",
			params: map[string]interface{}{
				"name":           "role",
				"clusterrole":    "admin",
				"user":           []string{"user"},
				"group":          []string{"group"},
				"serviceaccount": "ns1",
			},
			expectErr: true,
		},
	}
	generator := ClusterRoleBindingGeneratorV1{}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := generator.Generate(tt.params)
			if !tt.expectErr && err != nil {
				t.Errorf("[%d] unexpected error: %v", i, err)
			}
			if tt.expectErr && err != nil {
				return
			}
			if tt.expectErr && err == nil {
				t.Errorf("[%s] expect error, got nil", tt.name)
			}
			if !reflect.DeepEqual(obj.(*rbacv1beta1.ClusterRoleBinding), tt.expected) {
				t.Errorf("\n[%s] want:\n%#v\ngot:\n%#v", tt.name, tt.expected, obj.(*rbacv1beta1.ClusterRoleBinding))
			}
		})
	}
}
