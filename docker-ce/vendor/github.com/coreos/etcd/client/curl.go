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
// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	cURLDebug = false
)

func EnablecURLDebug() {
	cURLDebug = true
}

func DisablecURLDebug() {
	cURLDebug = false
}

// printcURL prints the cURL equivalent request to stderr.
// It returns an error if the body of the request cannot
// be read.
// The caller MUST cancel the request if there is an error.
func printcURL(req *http.Request) error {
	if !cURLDebug {
		return nil
	}
	var (
		command string
		b       []byte
		err     error
	)

	if req.URL != nil {
		command = fmt.Sprintf("curl -X %s %s", req.Method, req.URL.String())
	}

	if req.Body != nil {
		b, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		command += fmt.Sprintf(" -d %q", string(b))
	}

	fmt.Fprintf(os.Stderr, "cURL Command: %s\n", command)

	// reset body
	body := bytes.NewBuffer(b)
	req.Body = ioutil.NopCloser(body)

	return nil
}
