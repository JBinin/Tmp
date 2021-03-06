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
/**
 *  Copyright 2015 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package cachecontrol

import (
	"github.com/pquerna/cachecontrol/cacheobject"

	"net/http"
	"time"
)

type Options struct {
	// Set to True for a prviate cache, which is not shared amoung users (eg, in a browser)
	// Set to False for a "shared" cache, which is more common in a server context.
	PrivateCache bool
}

// Given an HTTP Request, the future Status Code, and an ResponseWriter,
// determine the possible reasons a response SHOULD NOT be cached.
func CachableResponseWriter(req *http.Request,
	statusCode int,
	resp http.ResponseWriter,
	opts Options) ([]cacheobject.Reason, time.Time, error) {
	return cacheobject.UsingRequestResponse(req, statusCode, resp.Header(), opts.PrivateCache)
}

// Given an HTTP Request and Response, determine the possible reasons a response SHOULD NOT
// be cached.
func CachableResponse(req *http.Request,
	resp *http.Response,
	opts Options) ([]cacheobject.Reason, time.Time, error) {
	return cacheobject.UsingRequestResponse(req, resp.StatusCode, resp.Header, opts.PrivateCache)
}
