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
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerAttach attaches a connection to a container in the server.
// It returns a types.HijackedConnection with the hijacked connection
// and the a reader to get output. It's up to the called to close
// the hijacked connection by calling types.HijackedResponse.Close.
func (cli *Client) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	query := url.Values{}
	if options.Stream {
		query.Set("stream", "1")
	}
	if options.Stdin {
		query.Set("stdin", "1")
	}
	if options.Stdout {
		query.Set("stdout", "1")
	}
	if options.Stderr {
		query.Set("stderr", "1")
	}
	if options.DetachKeys != "" {
		query.Set("detachKeys", options.DetachKeys)
	}
	if options.Logs {
		query.Set("logs", "1")
	}

	headers := map[string][]string{"Content-Type": {"text/plain"}}
	return cli.postHijacked(ctx, "/containers/"+container+"/attach", query, nil, headers)
}
