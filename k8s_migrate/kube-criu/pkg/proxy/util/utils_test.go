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
	"net"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	fake "k8s.io/kubernetes/pkg/proxy/util/testing"
)

func TestShouldSkipService(t *testing.T) {
	testCases := []struct {
		service    *v1.Service
		svcName    types.NamespacedName
		shouldSkip bool
	}{
		{
			// Cluster IP is None
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"},
				Spec: v1.ServiceSpec{
					ClusterIP: v1.ClusterIPNone,
				},
			},
			svcName:    types.NamespacedName{Namespace: "foo", Name: "bar"},
			shouldSkip: true,
		},
		{
			// Cluster IP is empty
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"},
				Spec: v1.ServiceSpec{
					ClusterIP: "",
				},
			},
			svcName:    types.NamespacedName{Namespace: "foo", Name: "bar"},
			shouldSkip: true,
		},
		{
			// ExternalName type service
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"},
				Spec: v1.ServiceSpec{
					ClusterIP: "1.2.3.4",
					Type:      v1.ServiceTypeExternalName,
				},
			},
			svcName:    types.NamespacedName{Namespace: "foo", Name: "bar"},
			shouldSkip: true,
		},
		{
			// ClusterIP type service with ClusterIP set
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"},
				Spec: v1.ServiceSpec{
					ClusterIP: "1.2.3.4",
					Type:      v1.ServiceTypeClusterIP,
				},
			},
			svcName:    types.NamespacedName{Namespace: "foo", Name: "bar"},
			shouldSkip: false,
		},
		{
			// NodePort type service with ClusterIP set
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"},
				Spec: v1.ServiceSpec{
					ClusterIP: "1.2.3.4",
					Type:      v1.ServiceTypeNodePort,
				},
			},
			svcName:    types.NamespacedName{Namespace: "foo", Name: "bar"},
			shouldSkip: false,
		},
		{
			// LoadBalancer type service with ClusterIP set
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foo", Name: "bar"},
				Spec: v1.ServiceSpec{
					ClusterIP: "1.2.3.4",
					Type:      v1.ServiceTypeLoadBalancer,
				},
			},
			svcName:    types.NamespacedName{Namespace: "foo", Name: "bar"},
			shouldSkip: false,
		},
	}

	for i := range testCases {
		skip := ShouldSkipService(testCases[i].svcName, testCases[i].service)
		if skip != testCases[i].shouldSkip {
			t.Errorf("case %d: expect %v, got %v", i, testCases[i].shouldSkip, skip)
		}
	}
}

type InterfaceAddrsPair struct {
	itf   net.Interface
	addrs []net.Addr
}

func TestGetNodeAddressses(t *testing.T) {
	testCases := []struct {
		cidrs         []string
		nw            *fake.FakeNetwork
		itfAddrsPairs []InterfaceAddrsPair
		expected      sets.String
	}{
		{ // case 0
			cidrs: []string{"10.20.30.0/24"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "10.20.30.51/24"}},
				},
				{
					itf:   net.Interface{Index: 2, MTU: 0, Name: "eth1", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "100.200.201.1/24"}},
				},
			},
			expected: sets.NewString("10.20.30.51"),
		},
		{ // case 1
			cidrs: []string{"0.0.0.0/0"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "10.20.30.51/24"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "127.0.0.1/8"}},
				},
			},
			expected: sets.NewString("0.0.0.0/0"),
		},
		{ // case 2
			cidrs: []string{"2001:db8::/32", "::1/128"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "2001:db8::1/32"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "::1/128"}},
				},
			},
			expected: sets.NewString("2001:db8::1", "::1"),
		},
		{ // case 3
			cidrs: []string{"::/0"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "2001:db8::1/32"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "::1/128"}},
				},
			},
			expected: sets.NewString("::/0"),
		},
		{ // case 4
			cidrs: []string{"127.0.0.1/32"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "10.20.30.51/24"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "127.0.0.1/8"}},
				},
			},
			expected: sets.NewString("127.0.0.1"),
		},
		{ // case 5
			cidrs: []string{"127.0.0.0/8"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "127.0.1.1/8"}},
				},
			},
			expected: sets.NewString("127.0.1.1"),
		},
		{ // case 6
			cidrs: []string{"10.20.30.0/24", "100.200.201.0/24"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "10.20.30.51/24"}},
				},
				{
					itf:   net.Interface{Index: 2, MTU: 0, Name: "eth1", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "100.200.201.1/24"}},
				},
			},
			expected: sets.NewString("10.20.30.51", "100.200.201.1"),
		},
		{ // case 7
			cidrs: []string{"10.20.30.0/24", "100.200.201.0/24"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "192.168.1.2/24"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "127.0.0.1/8"}},
				},
			},
			expected: sets.NewString(),
		},
		{ // case 8
			cidrs: []string{},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "192.168.1.2/24"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "127.0.0.1/8"}},
				},
			},
			expected: sets.NewString("0.0.0.0/0", "::/0"),
		},
		{ // case 9
			cidrs: []string{},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "2001:db8::1/32"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "::1/128"}},
				},
			},
			expected: sets.NewString("0.0.0.0/0", "::/0"),
		},
		{ // case 9
			cidrs: []string{"1.2.3.0/24", "0.0.0.0/0"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "1.2.3.4/30"}},
				},
			},
			expected: sets.NewString("0.0.0.0/0"),
		},
		{ // case 10
			cidrs: []string{"0.0.0.0/0", "1.2.3.0/24", "::1/128"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "1.2.3.4/30"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "::1/128"}},
				},
			},
			expected: sets.NewString("0.0.0.0/0", "::1"),
		},
		{ // case 11
			cidrs: []string{"::/0", "1.2.3.0/24", "::1/128"},
			nw:    fake.NewFakeNetwork(),
			itfAddrsPairs: []InterfaceAddrsPair{
				{
					itf:   net.Interface{Index: 0, MTU: 0, Name: "eth0", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "1.2.3.4/30"}},
				},
				{
					itf:   net.Interface{Index: 1, MTU: 0, Name: "lo", HardwareAddr: nil, Flags: 0},
					addrs: []net.Addr{fake.AddrStruct{Val: "::1/128"}},
				},
			},
			expected: sets.NewString("::/0", "1.2.3.4"),
		},
	}

	for i := range testCases {
		for _, pair := range testCases[i].itfAddrsPairs {
			testCases[i].nw.AddInterfaceAddr(&pair.itf, pair.addrs)
		}
		addrList, err := GetNodeAddresses(testCases[i].cidrs, testCases[i].nw)
		if err != nil {
			t.Errorf("case [%d], unexpected error: %v", i, err)
		}
		if !addrList.Equal(testCases[i].expected) {
			t.Errorf("case [%d], unexpected mismatch, expected: %v, got: %v", i, testCases[i].expected, addrList)
		}
	}
}
