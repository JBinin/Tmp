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

package discovery

import (
	"net/http"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
)

type cacheRoundTripper struct {
	rt *httpcache.Transport
}

// newCacheRoundTripper creates a roundtripper that reads the ETag on
// response headers and send the If-None-Match header on subsequent
// corresponding requests.
func newCacheRoundTripper(cacheDir string, rt http.RoundTripper) http.RoundTripper {
	d := diskv.New(diskv.Options{
		BasePath: cacheDir,
		TempDir:  filepath.Join(cacheDir, ".diskv-temp"),
	})
	t := httpcache.NewTransport(diskcache.NewWithDiskv(d))
	t.Transport = rt

	return &cacheRoundTripper{rt: t}
}

func (rt *cacheRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.rt.RoundTrip(req)
}

func (rt *cacheRoundTripper) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := rt.rt.Transport.(canceler); ok {
		cr.CancelRequest(req)
	} else {
		glog.Errorf("CancelRequest not implemented by %T", rt.rt.Transport)
	}
}

func (rt *cacheRoundTripper) WrappedRoundTripper() http.RoundTripper { return rt.rt.Transport }
