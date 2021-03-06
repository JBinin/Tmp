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

package rest

import (
	"path"
	"testing"

	"k8s.io/api/core/v1"
)

func TestValidatesHostParameter(t *testing.T) {
	testCases := []struct {
		Host    string
		APIPath string

		URL string
		Err bool
	}{
		{"127.0.0.1", "", "http://127.0.0.1/" + v1.SchemeGroupVersion.Version, false},
		{"127.0.0.1:8080", "", "http://127.0.0.1:8080/" + v1.SchemeGroupVersion.Version, false},
		{"foo.bar.com", "", "http://foo.bar.com/" + v1.SchemeGroupVersion.Version, false},
		{"http://host/prefix", "", "http://host/prefix/" + v1.SchemeGroupVersion.Version, false},
		{"http://host", "", "http://host/" + v1.SchemeGroupVersion.Version, false},
		{"http://host", "/", "http://host/" + v1.SchemeGroupVersion.Version, false},
		{"http://host", "/other", "http://host/other/" + v1.SchemeGroupVersion.Version, false},
		{"host/server", "", "", true},
	}
	for i, testCase := range testCases {
		u, versionedAPIPath, err := DefaultServerURL(testCase.Host, testCase.APIPath, v1.SchemeGroupVersion, false)
		switch {
		case err == nil && testCase.Err:
			t.Errorf("expected error but was nil")
			continue
		case err != nil && !testCase.Err:
			t.Errorf("unexpected error %v", err)
			continue
		case err != nil:
			continue
		}
		u.Path = path.Join(u.Path, versionedAPIPath)
		if e, a := testCase.URL, u.String(); e != a {
			t.Errorf("%d: expected host %s, got %s", i, e, a)
			continue
		}
	}
}
