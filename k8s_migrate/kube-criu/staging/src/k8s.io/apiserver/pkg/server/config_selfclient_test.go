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
Copyright 2016 The Kubernetes Authors.

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

package server

import (
	"net"
	"testing"
)

func TestLoopbackHostPort(t *testing.T) {
	host, port, err := LoopbackHostPort("1.2.3.4:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != "1.2.3.4" {
		t.Fatalf("expected 1.2.3.4 as host, got %q", host)
	}
	if port != "443" {
		t.Fatalf("expected 443 as port, got %q", port)
	}

	host, port, err = LoopbackHostPort("0.0.0.0:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		t.Fatalf("expected host to be loopback, got %q", host)
	}
	if port != "443" {
		t.Fatalf("expected 443 as port, got %q", port)
	}

	host, port, err = LoopbackHostPort("[ff06:0:0:0:0:0:0:c3]:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != "ff06:0:0:0:0:0:0:c3" {
		t.Fatalf("expected ff06:0:0:0:0:0:0:c3 as host, got %q", host)
	}
	if port != "443" {
		t.Fatalf("expected 443 as port, got %q", port)
	}

	host, port, err = LoopbackHostPort("[::]:443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() || ip.To4() != nil {
		t.Fatalf("expected IPv6 host to be loopback, got %q", host)
	}

	if port != "443" {
		t.Fatalf("expected 443 as port, got %q", port)
	}
}
