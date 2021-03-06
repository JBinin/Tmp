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

package auth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/rest/fake"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	cmdtesting "k8s.io/kubernetes/pkg/kubectl/cmd/testing"
)

func TestRunAccessCheck(t *testing.T) {
	tests := []struct {
		name      string
		o         *CanIOptions
		args      []string
		allowed   bool
		serverErr error

		expectedBodyStrings []string
	}{
		{
			name:    "restmapping for args",
			o:       &CanIOptions{},
			args:    []string{"get", "replicaset"},
			allowed: true,
			expectedBodyStrings: []string{
				`{"resourceAttributes":{"namespace":"test","verb":"get","group":"extensions","resource":"replicasets"}}`,
			},
		},
		{
			name:    "simple success",
			o:       &CanIOptions{},
			args:    []string{"get", "deployments.extensions/foo"},
			allowed: true,
			expectedBodyStrings: []string{
				`{"resourceAttributes":{"namespace":"test","verb":"get","group":"extensions","resource":"deployments","name":"foo"}}`,
			},
		},
		{
			name: "all namespaces",
			o: &CanIOptions{
				AllNamespaces: true,
			},
			args:    []string{"get", "deployments.extensions/foo"},
			allowed: true,
			expectedBodyStrings: []string{
				`{"resourceAttributes":{"verb":"get","group":"extensions","resource":"deployments","name":"foo"}}`,
			},
		},
		{
			name: "disallowed",
			o: &CanIOptions{
				AllNamespaces: true,
			},
			args:    []string{"get", "deployments.extensions/foo"},
			allowed: false,
			expectedBodyStrings: []string{
				`{"resourceAttributes":{"verb":"get","group":"extensions","resource":"deployments","name":"foo"}}`,
			},
		},
		{
			name: "forcedError",
			o: &CanIOptions{
				AllNamespaces: true,
			},
			args:      []string{"get", "deployments.extensions/foo"},
			allowed:   false,
			serverErr: fmt.Errorf("forcedError"),
			expectedBodyStrings: []string{
				`{"resourceAttributes":{"verb":"get","group":"extensions","resource":"deployments","name":"foo"}}`,
			},
		},
		{
			name: "sub resource",
			o: &CanIOptions{
				AllNamespaces: true,
				Subresource:   "log",
			},
			args:    []string{"get", "pods"},
			allowed: true,
			expectedBodyStrings: []string{
				`{"resourceAttributes":{"verb":"get","resource":"pods","subresource":"log"}}`,
			},
		},
		{
			name:    "nonResourceURL",
			o:       &CanIOptions{},
			args:    []string{"get", "/logs"},
			allowed: true,
			expectedBodyStrings: []string{
				`{"nonResourceAttributes":{"path":"/logs","verb":"get"}}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.o.Out = ioutil.Discard
			test.o.ErrOut = ioutil.Discard

			tf := cmdtesting.NewTestFactory().WithNamespace("test")
			defer tf.Cleanup()

			ns := legacyscheme.Codecs

			tf.Client = &fake.RESTClient{
				GroupVersion:         schema.GroupVersion{Group: "", Version: "v1"},
				NegotiatedSerializer: ns,
				Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
					expectPath := "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews"
					if req.URL.Path != expectPath {
						t.Errorf("%s: expected %v, got %v", test.name, expectPath, req.URL.Path)
						return nil, nil
					}
					bodyBits, err := ioutil.ReadAll(req.Body)
					if err != nil {
						t.Errorf("%s: %v", test.name, err)
						return nil, nil
					}
					body := string(bodyBits)

					for _, expectedBody := range test.expectedBodyStrings {
						if !strings.Contains(body, expectedBody) {
							t.Errorf("%s expecting %s in %s", test.name, expectedBody, body)
						}
					}

					return &http.Response{
							StatusCode: http.StatusOK,
							Body: ioutil.NopCloser(bytes.NewBufferString(
								fmt.Sprintf(`{"kind":"SelfSubjectAccessReview","apiVersion":"authorization.k8s.io/v1","status":{"allowed":%v}}`, test.allowed),
							)),
						},
						test.serverErr
				}),
			}
			tf.ClientConfigVal = &restclient.Config{ContentConfig: restclient.ContentConfig{GroupVersion: &schema.GroupVersion{Group: "", Version: "v1"}}}

			if err := test.o.Complete(tf, test.args); err != nil {
				t.Errorf("%s: %v", test.name, err)
				return
			}

			actualAllowed, err := test.o.RunAccessCheck()
			switch {
			case test.serverErr == nil && err == nil:
				// pass
			case err != nil && test.serverErr != nil && strings.Contains(err.Error(), test.serverErr.Error()):
				// pass
			default:
				t.Errorf("%s: expected %v, got %v", test.name, test.serverErr, err)
				return
			}
			if actualAllowed != test.allowed {
				t.Errorf("%s: expected %v, got %v", test.name, test.allowed, actualAllowed)
				return
			}
		})
	}
}
