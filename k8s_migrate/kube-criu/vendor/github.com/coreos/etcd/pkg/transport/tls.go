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
// Copyright 2016 The etcd Authors
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

package transport

import (
	"fmt"
	"strings"
	"time"
)

// ValidateSecureEndpoints scans the given endpoints against tls info, returning only those
// endpoints that could be validated as secure.
func ValidateSecureEndpoints(tlsInfo TLSInfo, eps []string) ([]string, error) {
	t, err := NewTransport(tlsInfo, 5*time.Second)
	if err != nil {
		return nil, err
	}
	var errs []string
	var endpoints []string
	for _, ep := range eps {
		if !strings.HasPrefix(ep, "https://") {
			errs = append(errs, fmt.Sprintf("%q is insecure", ep))
			continue
		}
		conn, cerr := t.Dial("tcp", ep[len("https://"):])
		if cerr != nil {
			errs = append(errs, fmt.Sprintf("%q failed to dial (%v)", ep, cerr))
			continue
		}
		conn.Close()
		endpoints = append(endpoints, ep)
	}
	if len(errs) != 0 {
		err = fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return endpoints, err
}
