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

package capabilities

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	defaultCap := Capabilities{
		AllowPrivileged: false,
		PrivilegedSources: PrivilegedSources{
			HostNetworkSources: []string{},
			HostPIDSources:     []string{},
			HostIPCSources:     []string{},
		},
	}

	res := Get()
	if !reflect.DeepEqual(defaultCap, res) {
		t.Fatalf("expected Capabilities: %#v, got a non-default: %#v", defaultCap, res)
	}

	cap := Capabilities{
		PrivilegedSources: PrivilegedSources{
			HostNetworkSources: []string{"A", "B"},
		},
	}
	SetForTests(cap)

	res = Get()
	if !reflect.DeepEqual(cap, res) {
		t.Fatalf("expected Capabilities: %#v , got a different: %#v", cap, res)
	}
}
