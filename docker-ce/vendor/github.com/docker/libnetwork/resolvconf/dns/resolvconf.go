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
package dns

import (
	"regexp"
)

// IPLocalhost is a regex pattern for IPv4 or IPv6 loopback range.
const IPLocalhost = `((127\.([0-9]{1,3}\.){2}[0-9]{1,3})|(::1)$)`

// IPv4Localhost is a regex pattern for IPv4 localhost address range.
const IPv4Localhost = `(127\.([0-9]{1,3}\.){2}[0-9]{1,3})`

var localhostIPRegexp = regexp.MustCompile(IPLocalhost)
var localhostIPv4Regexp = regexp.MustCompile(IPv4Localhost)

// IsLocalhost returns true if ip matches the localhost IP regular expression.
// Used for determining if nameserver settings are being passed which are
// localhost addresses
func IsLocalhost(ip string) bool {
	return localhostIPRegexp.MatchString(ip)
}

// IsIPv4Localhost returns true if ip matches the IPv4 localhost regular expression.
func IsIPv4Localhost(ip string) bool {
	return localhostIPv4Regexp.MatchString(ip)
}
