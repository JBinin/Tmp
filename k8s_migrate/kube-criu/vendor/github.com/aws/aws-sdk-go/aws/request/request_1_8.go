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
// +build go1.8

package request

import (
	"net/http"
)

// NoBody is a http.NoBody reader instructing Go HTTP client to not include
// and body in the HTTP request.
var NoBody = http.NoBody

// ResetBody rewinds the request body back to its starting position, and
// sets the HTTP Request body reference. When the body is read prior
// to being sent in the HTTP request it will need to be rewound.
//
// ResetBody will automatically be called by the SDK's build handler, but if
// the request is being used directly ResetBody must be called before the request
// is Sent.  SetStringBody, SetBufferBody, and SetReaderBody will automatically
// call ResetBody.
//
// Will also set the Go 1.8's http.Request.GetBody member to allow retrying
// PUT/POST redirects.
func (r *Request) ResetBody() {
	body, err := r.getNextRequestBody()
	if err != nil {
		r.Error = err
		return
	}

	r.HTTPRequest.Body = body
	r.HTTPRequest.GetBody = r.getNextRequestBody
}
