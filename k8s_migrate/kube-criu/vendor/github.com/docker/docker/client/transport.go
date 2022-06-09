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
package client

import (
	"crypto/tls"
	"net/http"
)

// transportFunc allows us to inject a mock transport for testing. We define it
// here so we can detect the tlsconfig and return nil for only this type.
type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

// resolveTLSConfig attempts to resolve the TLS configuration from the
// RoundTripper.
func resolveTLSConfig(transport http.RoundTripper) *tls.Config {
	switch tr := transport.(type) {
	case *http.Transport:
		return tr.TLSClientConfig
	default:
		return nil
	}
}
