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
/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transport

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
)

// Config holds various options for establishing a transport.
type Config struct {
	// UserAgent is an optional field that specifies the caller of this
	// request.
	UserAgent string

	// The base TLS configuration for this transport.
	TLS TLSConfig

	// Username and password for basic authentication
	Username string
	Password string

	// Bearer token for authentication
	BearerToken string

	// Impersonate is the config that this Config will impersonate using
	Impersonate ImpersonationConfig

	// Transport may be used for custom HTTP behavior. This attribute may
	// not be specified with the TLS client certificate options. Use
	// WrapTransport for most client level operations.
	Transport http.RoundTripper

	// WrapTransport will be invoked for custom HTTP behavior after the
	// underlying transport is initialized (either the transport created
	// from TLSClientConfig, Transport, or http.DefaultTransport). The
	// config may layer other RoundTrippers on top of the returned
	// RoundTripper.
	WrapTransport func(rt http.RoundTripper) http.RoundTripper

	// Dial specifies the dial function for creating unencrypted TCP connections.
	Dial func(ctx context.Context, network, address string) (net.Conn, error)
}

// ImpersonationConfig has all the available impersonation options
type ImpersonationConfig struct {
	// UserName matches user.Info.GetName()
	UserName string
	// Groups matches user.Info.GetGroups()
	Groups []string
	// Extra matches user.Info.GetExtra()
	Extra map[string][]string
}

// HasCA returns whether the configuration has a certificate authority or not.
func (c *Config) HasCA() bool {
	return len(c.TLS.CAData) > 0 || len(c.TLS.CAFile) > 0
}

// HasBasicAuth returns whether the configuration has basic authentication or not.
func (c *Config) HasBasicAuth() bool {
	return len(c.Username) != 0
}

// HasTokenAuth returns whether the configuration has token authentication or not.
func (c *Config) HasTokenAuth() bool {
	return len(c.BearerToken) != 0
}

// HasCertAuth returns whether the configuration has certificate authentication or not.
func (c *Config) HasCertAuth() bool {
	return (len(c.TLS.CertData) != 0 || len(c.TLS.CertFile) != 0) && (len(c.TLS.KeyData) != 0 || len(c.TLS.KeyFile) != 0)
}

// HasCertCallbacks returns whether the configuration has certificate callback or not.
func (c *Config) HasCertCallback() bool {
	return c.TLS.GetCert != nil
}

// TLSConfig holds the information needed to set up a TLS transport.
type TLSConfig struct {
	CAFile   string // Path of the PEM-encoded server trusted root certificates.
	CertFile string // Path of the PEM-encoded client certificate.
	KeyFile  string // Path of the PEM-encoded client key.

	Insecure   bool   // Server should be accessed without verifying the certificate. For testing only.
	ServerName string // Override for the server name passed to the server for SNI and used to verify certificates.

	CAData   []byte // Bytes of the PEM-encoded server trusted root certificates. Supercedes CAFile.
	CertData []byte // Bytes of the PEM-encoded client certificate. Supercedes CertFile.
	KeyData  []byte // Bytes of the PEM-encoded client key. Supercedes KeyFile.

	GetCert func() (*tls.Certificate, error) // Callback that returns a TLS client certificate. CertData, CertFile, KeyData and KeyFile supercede this field.
}
