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

package kubectl

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSecretForDockerRegistryGenerate(t *testing.T) {
	username, password, email, server := "test-user", "test-password", "test-user@example.org", "https://index.docker.io/v1/"
	secretData, err := handleDockerCfgJsonContent(username, password, email, server)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	secretDataNoEmail, err := handleDockerCfgJsonContent(username, password, "", server)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tests := []struct {
		name      string
		params    map[string]interface{}
		expected  *v1.Secret
		expectErr bool
	}{
		{
			name: "test-valid-use",
			params: map[string]interface{}{
				"name":            "foo",
				"docker-server":   server,
				"docker-username": username,
				"docker-password": password,
				"docker-email":    email,
			},
			expected: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Data: map[string][]byte{
					v1.DockerConfigJsonKey: secretData,
				},
				Type: v1.SecretTypeDockerConfigJson,
			},
			expectErr: false,
		},
		{
			name: "test-valid-use-append-hash",
			params: map[string]interface{}{
				"name":            "foo",
				"docker-server":   server,
				"docker-username": username,
				"docker-password": password,
				"docker-email":    email,
				"append-hash":     true,
			},
			expected: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo-548cm7fgdh",
				},
				Data: map[string][]byte{
					v1.DockerConfigJsonKey: secretData,
				},
				Type: v1.SecretTypeDockerConfigJson,
			},
			expectErr: false,
		},
		{
			name: "test-valid-use-no-email",
			params: map[string]interface{}{
				"name":            "foo",
				"docker-server":   server,
				"docker-username": username,
				"docker-password": password,
			},
			expected: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Data: map[string][]byte{
					v1.DockerConfigJsonKey: secretDataNoEmail,
				},
				Type: v1.SecretTypeDockerConfigJson,
			},
			expectErr: false,
		},
		{
			name: "test-missing-required-param",
			params: map[string]interface{}{
				"name":            "foo",
				"docker-server":   server,
				"docker-password": password,
				"docker-email":    email,
			},
			expectErr: true,
		},
	}

	generator := SecretForDockerRegistryGeneratorV1{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := generator.Generate(tt.params)
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectErr && err != nil {
				return
			}
			if !reflect.DeepEqual(obj.(*v1.Secret), tt.expected) {
				t.Errorf("\nexpected:\n%#v\nsaw:\n%#v", tt.expected, obj.(*v1.Secret))
			}
		})
	}
}
