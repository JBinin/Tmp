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
package autorest

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// NewRetriableRequest returns a wrapper around an HTTP request that support retry logic.
func NewRetriableRequest(req *http.Request) *RetriableRequest {
	return &RetriableRequest{req: req}
}

// Request returns the wrapped HTTP request.
func (rr *RetriableRequest) Request() *http.Request {
	return rr.req
}

func (rr *RetriableRequest) prepareFromByteReader() (err error) {
	// fall back to making a copy (only do this once)
	b := []byte{}
	if rr.req.ContentLength > 0 {
		b = make([]byte, rr.req.ContentLength)
		_, err = io.ReadFull(rr.req.Body, b)
		if err != nil {
			return err
		}
	} else {
		b, err = ioutil.ReadAll(rr.req.Body)
		if err != nil {
			return err
		}
	}
	rr.br = bytes.NewReader(b)
	rr.req.Body = ioutil.NopCloser(rr.br)
	return err
}
