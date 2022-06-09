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
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.8

package http2

import (
	"io"
	"net/http"
)

func configureServer18(h1 *http.Server, h2 *Server) error {
	// No IdleTimeout to sync prior to Go 1.8.
	return nil
}

func shouldLogPanic(panicValue interface{}) bool {
	return panicValue != nil
}

func reqGetBody(req *http.Request) func() (io.ReadCloser, error) {
	return nil
}

func reqBodyIsNoBody(io.ReadCloser) bool { return false }
