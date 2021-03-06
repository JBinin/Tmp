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
package middleware

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// CORSMiddleware injects CORS headers to each request
// when it's configured.
type CORSMiddleware struct {
	defaultHeaders string
}

// NewCORSMiddleware creates a new CORSMiddleware with default headers.
func NewCORSMiddleware(d string) CORSMiddleware {
	return CORSMiddleware{defaultHeaders: d}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (c CORSMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		// If "api-cors-header" is not given, but "api-enable-cors" is true, we set cors to "*"
		// otherwise, all head values will be passed to HTTP handler
		corsHeaders := c.defaultHeaders
		if corsHeaders == "" {
			corsHeaders = "*"
		}

		logrus.Debugf("CORS header is enabled and set to: %s", corsHeaders)
		w.Header().Add("Access-Control-Allow-Origin", corsHeaders)
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Registry-Auth")
		w.Header().Add("Access-Control-Allow-Methods", "HEAD, GET, POST, DELETE, PUT, OPTIONS")
		return handler(ctx, w, r, vars)
	}
}
