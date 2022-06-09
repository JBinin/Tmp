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
package swag

import (
	"net"
	"strconv"
)

// SplitHostPort splits a network address into a host and a port.
// The port is -1 when there is no port to be found
func SplitHostPort(addr string) (host string, port int, err error) {
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		return "", -1, err
	}
	if p == "" {
		return "", -1, &net.AddrError{Err: "missing port in address", Addr: addr}
	}

	pi, err := strconv.Atoi(p)
	if err != nil {
		return "", -1, err
	}
	return h, pi, nil
}
