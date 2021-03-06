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
Copyright 2014 The Kubernetes Authors.

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

package rest

import (
	"crypto/tls"
	"errors"
	"net/http"

	"k8s.io/client-go/plugin/pkg/client/auth/exec"
	"k8s.io/client-go/transport"
)

// TLSConfigFor returns a tls.Config that will provide the transport level security defined
// by the provided Config. Will return nil if no transport level security is requested.
func TLSConfigFor(config *Config) (*tls.Config, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.TLSConfigFor(cfg)
}

// TransportFor returns an http.RoundTripper that will provide the authentication
// or transport level security defined by the provided Config. Will return the
// default http.DefaultTransport if no special case behavior is needed.
func TransportFor(config *Config) (http.RoundTripper, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.New(cfg)
}

// HTTPWrappersForConfig wraps a round tripper with any relevant layered behavior from the
// config. Exposed to allow more clients that need HTTP-like behavior but then must hijack
// the underlying connection (like WebSocket or HTTP2 clients). Pure HTTP clients should use
// the higher level TransportFor or RESTClientFor methods.
func HTTPWrappersForConfig(config *Config, rt http.RoundTripper) (http.RoundTripper, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.HTTPWrappersForConfig(cfg, rt)
}

// TransportConfig converts a client config to an appropriate transport config.
func (c *Config) TransportConfig() (*transport.Config, error) {
	conf := &transport.Config{
		UserAgent:     c.UserAgent,
		Transport:     c.Transport,
		WrapTransport: c.WrapTransport,
		TLS: transport.TLSConfig{
			Insecure:   c.Insecure,
			ServerName: c.ServerName,
			CAFile:     c.CAFile,
			CAData:     c.CAData,
			CertFile:   c.CertFile,
			CertData:   c.CertData,
			KeyFile:    c.KeyFile,
			KeyData:    c.KeyData,
		},
		Username:    c.Username,
		Password:    c.Password,
		BearerToken: c.BearerToken,
		Impersonate: transport.ImpersonationConfig{
			UserName: c.Impersonate.UserName,
			Groups:   c.Impersonate.Groups,
			Extra:    c.Impersonate.Extra,
		},
		Dial: c.Dial,
	}

	if c.ExecProvider != nil && c.AuthProvider != nil {
		return nil, errors.New("execProvider and authProvider cannot be used in combination")
	}

	if c.ExecProvider != nil {
		provider, err := exec.GetAuthenticator(c.ExecProvider)
		if err != nil {
			return nil, err
		}
		if err := provider.UpdateTransportConfig(conf); err != nil {
			return nil, err
		}
	}
	if c.AuthProvider != nil {
		provider, err := GetAuthProvider(c.Host, c.AuthProvider, c.AuthConfigPersister)
		if err != nil {
			return nil, err
		}
		wt := conf.WrapTransport
		if wt != nil {
			conf.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
				return provider.WrapTransport(wt(rt))
			}
		} else {
			conf.WrapTransport = provider.WrapTransport
		}
	}
	return conf, nil
}
