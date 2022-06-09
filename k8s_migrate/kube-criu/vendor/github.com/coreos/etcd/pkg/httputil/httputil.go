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
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// borrowed from golang/net/context/ctxhttp/cancelreq.go

// Package httputil provides HTTP utility functions.
package httputil

import (
	"io"
	"io/ioutil"
	"net/http"
)

// GracefulClose drains http.Response.Body until it hits EOF
// and closes it. This prevents TCP/TLS connections from closing,
// therefore available for reuse.
func GracefulClose(resp *http.Response) {
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}
