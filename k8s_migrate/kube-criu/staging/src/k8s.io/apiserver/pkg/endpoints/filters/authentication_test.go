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

package filters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

func TestAuthenticateRequest(t *testing.T) {
	success := make(chan struct{})
	auth := WithAuthentication(
		http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			user, ok := genericapirequest.UserFrom(ctx)
			if user == nil || !ok {
				t.Errorf("no user stored in context: %#v", ctx)
			}
			if req.Header.Get("Authorization") != "" {
				t.Errorf("Authorization header should be removed from request on success: %#v", req)
			}
			close(success)
		}),
		authenticator.RequestFunc(func(req *http.Request) (user.Info, bool, error) {
			if req.Header.Get("Authorization") == "Something" {
				return &user.DefaultInfo{Name: "user"}, true, nil
			}
			return nil, false, errors.New("Authorization header is missing.")
		}),
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Errorf("unexpected call to failed")
		}),
	)

	auth.ServeHTTP(httptest.NewRecorder(), &http.Request{Header: map[string][]string{"Authorization": {"Something"}}})

	<-success
}

func TestAuthenticateRequestFailed(t *testing.T) {
	failed := make(chan struct{})
	auth := WithAuthentication(
		http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			t.Errorf("unexpected call to handler")
		}),
		authenticator.RequestFunc(func(req *http.Request) (user.Info, bool, error) {
			return nil, false, nil
		}),
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			close(failed)
		}),
	)

	auth.ServeHTTP(httptest.NewRecorder(), &http.Request{})

	<-failed
}

func TestAuthenticateRequestError(t *testing.T) {
	failed := make(chan struct{})
	auth := WithAuthentication(
		http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			t.Errorf("unexpected call to handler")
		}),
		authenticator.RequestFunc(func(req *http.Request) (user.Info, bool, error) {
			return nil, false, errors.New("failure")
		}),
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			close(failed)
		}),
	)

	auth.ServeHTTP(httptest.NewRecorder(), &http.Request{})

	<-failed
}
