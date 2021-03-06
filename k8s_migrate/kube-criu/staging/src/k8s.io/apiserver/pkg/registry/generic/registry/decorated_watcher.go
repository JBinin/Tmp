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

package registry

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type decoratedWatcher struct {
	w         watch.Interface
	decorator ObjectFunc
	cancel    context.CancelFunc
	resultCh  chan watch.Event
}

func newDecoratedWatcher(w watch.Interface, decorator ObjectFunc) *decoratedWatcher {
	ctx, cancel := context.WithCancel(context.Background())
	d := &decoratedWatcher{
		w:         w,
		decorator: decorator,
		cancel:    cancel,
		resultCh:  make(chan watch.Event),
	}
	go d.run(ctx)
	return d
}

func (d *decoratedWatcher) run(ctx context.Context) {
	var recv, send watch.Event
	for {
		select {
		case recv = <-d.w.ResultChan():
			switch recv.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				err := d.decorator(recv.Object)
				if err != nil {
					send = makeStatusErrorEvent(err)
					break
				}
				send = recv
			case watch.Error:
				send = recv
			}
			select {
			case d.resultCh <- send:
				if send.Type == watch.Error {
					d.cancel()
				}
			case <-ctx.Done():
			}
		case <-ctx.Done():
			d.w.Stop()
			close(d.resultCh)
			return
		}
	}
}

func (d *decoratedWatcher) Stop() {
	d.cancel()
}

func (d *decoratedWatcher) ResultChan() <-chan watch.Event {
	return d.resultCh
}

func makeStatusErrorEvent(err error) watch.Event {
	status := &metav1.Status{
		Status:  metav1.StatusFailure,
		Message: err.Error(),
		Code:    http.StatusInternalServerError,
		Reason:  metav1.StatusReasonInternalError,
	}
	return watch.Event{
		Type:   watch.Error,
		Object: status,
	}
}
