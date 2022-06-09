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

package filters

import (
	"errors"
	"net/http"

	utilwaitgroup "k8s.io/apimachinery/pkg/util/waitgroup"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
)

// WithWaitGroup adds all non long-running requests to wait group, which is used for graceful shutdown.
func WithWaitGroup(handler http.Handler, longRunning apirequest.LongRunningRequestCheck, wg *utilwaitgroup.SafeWaitGroup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestInfo, ok := apirequest.RequestInfoFrom(ctx)
		if !ok {
			// if this happens, the handler chain isn't setup correctly because there is no request info
			responsewriters.InternalError(w, req, errors.New("no RequestInfo found in the context"))
			return
		}

		if !longRunning(req, requestInfo) {
			if err := wg.Add(1); err != nil {
				http.Error(w, "apiserver is shutting down.", http.StatusInternalServerError)
				return
			}
			defer wg.Done()
		}

		handler.ServeHTTP(w, req)
	})
}
