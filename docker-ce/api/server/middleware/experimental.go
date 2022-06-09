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

	"golang.org/x/net/context"
)

// ExperimentalMiddleware is a the middleware in charge of adding the
// 'Docker-Experimental' header to every outgoing request
type ExperimentalMiddleware struct {
	experimental string
}

// NewExperimentalMiddleware creates a new ExperimentalMiddleware
func NewExperimentalMiddleware(experimentalEnabled bool) ExperimentalMiddleware {
	if experimentalEnabled {
		return ExperimentalMiddleware{"true"}
	}
	return ExperimentalMiddleware{"false"}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (e ExperimentalMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		w.Header().Set("Docker-Experimental", e.experimental)
		return handler(ctx, w, r, vars)
	}
}
