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

package openapi

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/spec"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/endpoints/request"
)

// Downloader is the OpenAPI downloader type. It will try to download spec from /swagger.json endpoint.
type Downloader struct {
}

// NewDownloader creates a new OpenAPI Downloader.
func NewDownloader() Downloader {
	return Downloader{}
}

// inMemoryResponseWriter is a http.Writer that keep the response in memory.
type inMemoryResponseWriter struct {
	writeHeaderCalled bool
	header            http.Header
	respCode          int
	data              []byte
}

func newInMemoryResponseWriter() *inMemoryResponseWriter {
	return &inMemoryResponseWriter{header: http.Header{}}
}

func (r *inMemoryResponseWriter) Header() http.Header {
	return r.header
}

func (r *inMemoryResponseWriter) WriteHeader(code int) {
	r.writeHeaderCalled = true
	r.respCode = code
}

func (r *inMemoryResponseWriter) Write(in []byte) (int, error) {
	if !r.writeHeaderCalled {
		r.WriteHeader(http.StatusOK)
	}
	r.data = append(r.data, in...)
	return len(in), nil
}

func (r *inMemoryResponseWriter) String() string {
	s := fmt.Sprintf("ResponseCode: %d", r.respCode)
	if r.data != nil {
		s += fmt.Sprintf(", Body: %s", string(r.data))
	}
	if r.header != nil {
		s += fmt.Sprintf(", Header: %s", r.header)
	}
	return s
}

func (s *Downloader) handlerWithUser(handler http.Handler, info user.Info) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(request.WithUser(req.Context(), info))
		handler.ServeHTTP(w, req)
	})
}

func etagFor(data []byte) string {
	return fmt.Sprintf("%s%X\"", locallyGeneratedEtagPrefix, sha512.Sum512(data))
}

// Download downloads openAPI spec from /swagger.json endpoint of the given handler.
// httpStatus is only valid if err == nil
func (s *Downloader) Download(handler http.Handler, etag string) (returnSpec *spec.Swagger, newEtag string, httpStatus int, err error) {
	handler = s.handlerWithUser(handler, &user.DefaultInfo{Name: aggregatorUser})
	handler = http.TimeoutHandler(handler, specDownloadTimeout, "request timed out")

	req, err := http.NewRequest("GET", "/openapi/v2", nil)
	if err != nil {
		return nil, "", 0, err
	}
	req.Header.Add("Accept", "application/json")

	// Only pass eTag if it is not generated locally
	if len(etag) > 0 && !strings.HasPrefix(etag, locallyGeneratedEtagPrefix) {
		req.Header.Add("If-None-Match", etag)
	}

	writer := newInMemoryResponseWriter()
	handler.ServeHTTP(writer, req)

	// single endpoint not found/registered in old server, try to fetch old endpoint
	// TODO(roycaihw): remove this in 1.11
	if writer.respCode == http.StatusForbidden || writer.respCode == http.StatusNotFound {
		req, err = http.NewRequest("GET", "/swagger.json", nil)
		if err != nil {
			return nil, "", 0, err
		}

		// Only pass eTag if it is not generated locally
		if len(etag) > 0 && !strings.HasPrefix(etag, locallyGeneratedEtagPrefix) {
			req.Header.Add("If-None-Match", etag)
		}

		writer = newInMemoryResponseWriter()
		handler.ServeHTTP(writer, req)
	}

	switch writer.respCode {
	case http.StatusNotModified:
		if len(etag) == 0 {
			return nil, etag, http.StatusNotModified, fmt.Errorf("http.StatusNotModified is not allowed in absence of etag")
		}
		return nil, etag, http.StatusNotModified, nil
	case http.StatusNotFound:
		// Gracefully skip 404, assuming the server won't provide any spec
		return nil, "", http.StatusNotFound, nil
	case http.StatusOK:
		openAPISpec := &spec.Swagger{}
		if err := json.Unmarshal(writer.data, openAPISpec); err != nil {
			return nil, "", 0, err
		}
		newEtag = writer.Header().Get("Etag")
		if len(newEtag) == 0 {
			newEtag = etagFor(writer.data)
			if len(etag) > 0 && strings.HasPrefix(etag, locallyGeneratedEtagPrefix) {
				// The function call with an etag and server does not report an etag.
				// That means this server does not support etag and the etag that passed
				// to the function generated previously by us. Just compare etags and
				// return StatusNotModified if they are the same.
				if etag == newEtag {
					return nil, etag, http.StatusNotModified, nil
				}
			}
		}
		return openAPISpec, newEtag, http.StatusOK, nil
	default:
		return nil, "", 0, fmt.Errorf("failed to retrieve openAPI spec, http error: %s", writer.String())
	}
}
