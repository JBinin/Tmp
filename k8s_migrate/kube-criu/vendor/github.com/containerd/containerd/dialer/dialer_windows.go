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
package dialer

import (
	"net"
	"os"
	"syscall"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func isNoent(err error) bool {
	if err != nil {
		if oerr, ok := err.(*os.PathError); ok {
			if oerr.Err == syscall.ENOENT {
				return true
			}
		}
	}
	return false
}

func dialer(address string, timeout time.Duration) (net.Conn, error) {
	return winio.DialPipe(address, &timeout)
}

// DialAddress returns the dial address
func DialAddress(address string) string {
	return address
}
