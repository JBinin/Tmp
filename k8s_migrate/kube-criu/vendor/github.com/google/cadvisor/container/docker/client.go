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
// Copyright 2015 Google Inc. All Rights Reserved.
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

// Handler for /validate content.
// Validates cadvisor dependencies - kernel, os, docker setup.

package docker

import (
	"net/http"
	"sync"

	dclient "github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

var (
	dockerClient     *dclient.Client
	dockerClientErr  error
	dockerClientOnce sync.Once
)

// Client creates a Docker API client based on the given Docker flags
func Client() (*dclient.Client, error) {
	dockerClientOnce.Do(func() {
		var client *http.Client
		if *ArgDockerTLS {
			client = &http.Client{}
			options := tlsconfig.Options{
				CAFile:             *ArgDockerCA,
				CertFile:           *ArgDockerCert,
				KeyFile:            *ArgDockerKey,
				InsecureSkipVerify: false,
			}
			tlsc, err := tlsconfig.Client(options)
			if err != nil {
				dockerClientErr = err
				return
			}
			client.Transport = &http.Transport{
				TLSClientConfig: tlsc,
			}
		}
		dockerClient, dockerClientErr = dclient.NewClient(*ArgDockerEndpoint,
			"",
			client,
			nil)

	})
	return dockerClient, dockerClientErr
}
