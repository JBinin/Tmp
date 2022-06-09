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
package registry

import (
	"net/url"
	"strings"

	"github.com/docker/go-connections/tlsconfig"
)

func (s *DefaultService) lookupV2Endpoints(hostname string) (endpoints []APIEndpoint, err error) {
	tlsConfig := tlsconfig.ServerDefault()
	if hostname == DefaultNamespace || hostname == IndexHostname {
		// v2 mirrors
		for _, mirror := range s.config.Mirrors {
			if !strings.HasPrefix(mirror, "http://") && !strings.HasPrefix(mirror, "https://") {
				mirror = "https://" + mirror
			}
			mirrorURL, err := url.Parse(mirror)
			if err != nil {
				return nil, err
			}
			mirrorTLSConfig, err := s.tlsConfigForMirror(mirrorURL)
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, APIEndpoint{
				URL: mirrorURL,
				// guess mirrors are v2
				Version:      APIVersion2,
				Mirror:       true,
				TrimHostname: true,
				TLSConfig:    mirrorTLSConfig,
			})
		}
		// v2 registry
		endpoints = append(endpoints, APIEndpoint{
			URL:          DefaultV2Registry,
			Version:      APIVersion2,
			Official:     true,
			TrimHostname: true,
			TLSConfig:    tlsConfig,
		})

		return endpoints, nil
	}

	tlsConfig, err = s.tlsConfig(hostname)
	if err != nil {
		return nil, err
	}

	endpoints = []APIEndpoint{
		{
			URL: &url.URL{
				Scheme: "https",
				Host:   hostname,
			},
			Version:      APIVersion2,
			TrimHostname: true,
			TLSConfig:    tlsConfig,
		},
	}

	if tlsConfig.InsecureSkipVerify {
		endpoints = append(endpoints, APIEndpoint{
			URL: &url.URL{
				Scheme: "http",
				Host:   hostname,
			},
			Version:      APIVersion2,
			TrimHostname: true,
			// used to check if supposed to be secure via InsecureSkipVerify
			TLSConfig: tlsConfig,
		})
	}

	return endpoints, nil
}
