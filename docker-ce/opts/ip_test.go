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
package opts

import (
	"net"
	"testing"
)

func TestIpOptString(t *testing.T) {
	addresses := []string{"", "0.0.0.0"}
	var ip net.IP

	for _, address := range addresses {
		stringAddress := NewIPOpt(&ip, address).String()
		if stringAddress != address {
			t.Fatalf("IpOpt string should be `%s`, not `%s`", address, stringAddress)
		}
	}
}

func TestNewIpOptInvalidDefaultVal(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	defaultVal := "Not an ip"

	ipOpt := NewIPOpt(&ip, defaultVal)

	expected := "127.0.0.1"
	if ipOpt.String() != expected {
		t.Fatalf("Expected [%v], got [%v]", expected, ipOpt.String())
	}
}

func TestNewIpOptValidDefaultVal(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	defaultVal := "192.168.1.1"

	ipOpt := NewIPOpt(&ip, defaultVal)

	expected := "192.168.1.1"
	if ipOpt.String() != expected {
		t.Fatalf("Expected [%v], got [%v]", expected, ipOpt.String())
	}
}

func TestIpOptSetInvalidVal(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	ipOpt := &IPOpt{IP: &ip}

	invalidIP := "invalid ip"
	expectedError := "invalid ip is not an ip address"
	err := ipOpt.Set(invalidIP)
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected an Error with [%v], got [%v]", expectedError, err.Error())
	}
}
