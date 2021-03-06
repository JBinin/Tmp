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
	"net/http"
	"runtime/debug"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/server/httplog"
)

// WithPanicRecovery wraps an http Handler to recover and log panics.
func WithPanicRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer runtime.HandleCrash(func(err interface{}) {
			http.Error(w, "This request caused apiserver to panic. Look in the logs for details.", http.StatusInternalServerError)
			glog.Errorf("apiserver panic'd on %v %v: %v\n%s\n", req.Method, req.RequestURI, err, debug.Stack())
		})

		logger := httplog.NewLogged(req, &w)
		defer logger.Log()

		// Dispatch to the internal handler
		handler.ServeHTTP(w, req)
	})
}
