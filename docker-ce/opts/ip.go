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
	"fmt"
	"net"
)

// IPOpt holds an IP. It is used to store values from CLI flags.
type IPOpt struct {
	*net.IP
}

// NewIPOpt creates a new IPOpt from a reference net.IP and a
// string representation of an IP. If the string is not a valid
// IP it will fallback to the specified reference.
func NewIPOpt(ref *net.IP, defaultVal string) *IPOpt {
	o := &IPOpt{
		IP: ref,
	}
	o.Set(defaultVal)
	return o
}

// Set sets an IPv4 or IPv6 address from a given string. If the given
// string is not parseable as an IP address it returns an error.
func (o *IPOpt) Set(val string) error {
	ip := net.ParseIP(val)
	if ip == nil {
		return fmt.Errorf("%s is not an ip address", val)
	}
	*o.IP = ip
	return nil
}

// String returns the IP address stored in the IPOpt. If stored IP is a
// nil pointer, it returns an empty string.
func (o *IPOpt) String() string {
	if *o.IP == nil {
		return ""
	}
	return o.IP.String()
}

// Type returns the type of the option
func (o *IPOpt) Type() string {
	return "ip"
}
