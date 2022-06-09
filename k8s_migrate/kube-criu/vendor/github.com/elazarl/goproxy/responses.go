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
package goproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Will generate a valid http response to the given request the response will have
// the given contentType, and http status.
// Typical usage, refuse to process requests to local addresses:
//
//	proxy.OnRequest(IsLocalHost()).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request,*http.Response) {
//		return nil,NewResponse(r,goproxy.ContentTypeHtml,http.StatusUnauthorized,
//			`<!doctype html><html><head><title>Can't use proxy for local addresses</title></head><body/></html>`)
//	})
func NewResponse(r *http.Request, contentType string, status int, body string) *http.Response {
	resp := &http.Response{}
	resp.Request = r
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", contentType)
	resp.StatusCode = status
	buf := bytes.NewBufferString(body)
	resp.ContentLength = int64(buf.Len())
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

const (
	ContentTypeText = "text/plain"
	ContentTypeHtml = "text/html"
)

// Alias for NewResponse(r,ContentTypeText,http.StatusAccepted,text)
func TextResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, ContentTypeText, http.StatusAccepted, text)
}
