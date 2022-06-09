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
Copyright 2017 The Kubernetes Authors.

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

package util

import (
	"fmt"
	"net/url"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FromApiserverCache modifies <opts> so that the GET request will
// be served from apiserver cache instead of from etcd.
func FromApiserverCache(opts *metav1.GetOptions) {
	opts.ResourceVersion = "0"
}

func parseEndpoint(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}

	if u.Scheme == "tcp" {
		return "tcp", u.Host, nil
	} else if u.Scheme == "unix" {
		return "unix", u.Path, nil
	} else if u.Scheme == "" {
		return "", "", fmt.Errorf("Using %q as endpoint is deprecated, please consider using full url format", endpoint)
	} else {
		return u.Scheme, "", fmt.Errorf("protocol %q not supported", u.Scheme)
	}
}
