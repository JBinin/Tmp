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
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/api/server/httputils"
	"golang.org/x/net/context"
)

func TestVersionMiddleware(t *testing.T) {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		if httputils.VersionFromContext(ctx) == "" {
			t.Fatalf("Expected version, got empty string")
		}
		return nil
	}

	defaultVersion := "1.10.0"
	minVersion := "1.2.0"
	m := NewVersionMiddleware(defaultVersion, defaultVersion, minVersion)
	h := m.WrapHandler(handler)

	req, _ := http.NewRequest("GET", "/containers/json", nil)
	resp := httptest.NewRecorder()
	ctx := context.Background()
	if err := h(ctx, resp, req, map[string]string{}); err != nil {
		t.Fatal(err)
	}
}

func TestVersionMiddlewareWithErrors(t *testing.T) {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		if httputils.VersionFromContext(ctx) == "" {
			t.Fatalf("Expected version, got empty string")
		}
		return nil
	}

	defaultVersion := "1.10.0"
	minVersion := "1.2.0"
	m := NewVersionMiddleware(defaultVersion, defaultVersion, minVersion)
	h := m.WrapHandler(handler)

	req, _ := http.NewRequest("GET", "/containers/json", nil)
	resp := httptest.NewRecorder()
	ctx := context.Background()

	vars := map[string]string{"version": "0.1"}
	err := h(ctx, resp, req, vars)

	if !strings.Contains(err.Error(), "client version 0.1 is too old. Minimum supported API version is 1.2.0") {
		t.Fatalf("Expected too old client error, got %v", err)
	}
}
