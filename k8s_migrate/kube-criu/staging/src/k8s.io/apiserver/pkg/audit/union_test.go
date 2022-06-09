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

package audit

import (
	"strconv"
	"testing"

	"k8s.io/apimachinery/pkg/types"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
)

type fakeBackend struct {
	events []*auditinternal.Event
}

func (f *fakeBackend) ProcessEvents(events ...*auditinternal.Event) {
	f.events = append(f.events, events...)
}

func (f *fakeBackend) Run(stopCh <-chan struct{}) error {
	return nil
}

func (f *fakeBackend) Shutdown() {
	// Nothing to do here.
}

func (f *fakeBackend) String() string {
	return ""
}

func TestUnion(t *testing.T) {
	backends := []Backend{
		new(fakeBackend),
		new(fakeBackend),
		new(fakeBackend),
	}

	b := Union(backends...)

	n := 5

	var events []*auditinternal.Event
	for i := 0; i < n; i++ {
		events = append(events, &auditinternal.Event{
			AuditID: types.UID(strconv.Itoa(i)),
		})
	}
	b.ProcessEvents(events...)

	for i, b := range backends {
		// so we can inspect the underlying events.
		backend := b.(*fakeBackend)

		if got := len(backend.events); got != n {
			t.Errorf("backend %d wanted %d events, got %d", i, n, got)
			continue
		}
		for j, event := range backend.events {
			wantID := types.UID(strconv.Itoa(j))
			if event.AuditID != wantID {
				t.Errorf("backend %d event %d wanted id %s, got %s", i, j, wantID, event.AuditID)
			}
		}
	}
}
