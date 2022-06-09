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
package transport

import (
	"io"
	"net/http"
)

// httpTransport holds an http.RoundTripper
// and information about the scheme and address the transport
// sends request to.
type httpTransport struct {
	http.RoundTripper
	scheme string
	addr   string
}

// NewHTTPTransport creates a new httpTransport.
func NewHTTPTransport(r http.RoundTripper, scheme, addr string) Transport {
	return httpTransport{
		RoundTripper: r,
		scheme:       scheme,
		addr:         addr,
	}
}

// NewRequest creates a new http.Request and sets the URL
// scheme and address with the transport's fields.
func (t httpTransport) NewRequest(path string, data io.Reader) (*http.Request, error) {
	req, err := newHTTPRequest(path, data)
	if err != nil {
		return nil, err
	}
	req.URL.Scheme = t.scheme
	req.URL.Host = t.addr
	return req, nil
}
