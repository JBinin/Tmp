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
package client

import (
	"net"
	"net/http"

	"golang.org/x/net/context"
)

// DialSession returns a connection that can be used communication with daemon
func (cli *Client) DialSession(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
	req, err := http.NewRequest("POST", "/session", nil)
	if err != nil {
		return nil, err
	}
	req = cli.addHeaders(req, meta)

	return cli.setupHijackConn(req, proto)
}
