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

package filters

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
)

type recorder struct {
	lock  sync.Mutex
	count int
}

func (r *recorder) Record() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.count++
}

func (r *recorder) Count() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.count
}

func TestTimeout(t *testing.T) {
	sendResponse := make(chan struct{}, 1)
	writeErrors := make(chan error, 1)
	timeout := make(chan time.Time, 1)
	resp := "test response"
	timeoutErr := apierrors.NewServerTimeout(schema.GroupResource{Group: "foo", Resource: "bar"}, "get", 0)
	record := &recorder{}

	ts := httptest.NewServer(WithTimeout(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			<-sendResponse
			_, err := w.Write([]byte(resp))
			writeErrors <- err
		}),
		func(req *http.Request) (*http.Request, <-chan time.Time, func(), *apierrors.StatusError) {
			return req, timeout, record.Record, timeoutErr
		}))
	defer ts.Close()

	// No timeouts
	sendResponse <- struct{}{}
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("got res.StatusCode %d; expected %d", res.StatusCode, http.StatusOK)
	}
	body, _ := ioutil.ReadAll(res.Body)
	if string(body) != resp {
		t.Errorf("got body %q; expected %q", string(body), resp)
	}
	if err := <-writeErrors; err != nil {
		t.Errorf("got unexpected Write error on first request: %v", err)
	}
	if record.Count() != 0 {
		t.Errorf("invoked record method: %#v", record)
	}

	// Times out
	timeout <- time.Time{}
	res, err = http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusGatewayTimeout {
		t.Errorf("got res.StatusCode %d; expected %d", res.StatusCode, http.StatusServiceUnavailable)
	}
	body, _ = ioutil.ReadAll(res.Body)
	status := &metav1.Status{}
	if err := json.Unmarshal(body, status); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(status, &timeoutErr.ErrStatus) {
		t.Errorf("unexpected object: %s", diff.ObjectReflectDiff(&timeoutErr.ErrStatus, status))
	}
	if record.Count() != 1 {
		t.Errorf("did not invoke record method: %#v", record)
	}

	// Now try to send a response
	sendResponse <- struct{}{}
	if err := <-writeErrors; err != http.ErrHandlerTimeout {
		t.Errorf("got Write error of %v; expected %v", err, http.ErrHandlerTimeout)
	}
}
