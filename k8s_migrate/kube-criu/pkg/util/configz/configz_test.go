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

package configz

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigz(t *testing.T) {
	v, err := New("testing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	v.Set("blah")

	s := httptest.NewServer(http.HandlerFunc(handle))
	defer s.Close()

	resp, err := http.Get(s.URL + "/configz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(body) != `{"testing":"blah"}` {
		t.Fatalf("unexpected output: %s", body)
	}

	v.Set("bing")
	resp, err = http.Get(s.URL + "/configz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(body) != `{"testing":"bing"}` {
		t.Fatalf("unexpected output: %s", body)
	}

	Delete("testing")
	resp, err = http.Get(s.URL + "/configz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(body) != `{}` {
		t.Fatalf("unexpected output: %s", body)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected Content-Type: %s", resp.Header.Get("Content-Type"))
	}
}
