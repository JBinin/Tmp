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
	"io"
	"net/url"

	"golang.org/x/net/context"
)

// ContainerExport retrieves the raw contents of a container
// and returns them as an io.ReadCloser. It's up to the caller
// to close the stream.
func (cli *Client) ContainerExport(ctx context.Context, containerID string) (io.ReadCloser, error) {
	serverResp, err := cli.get(ctx, "/containers/"+containerID+"/export", url.Values{}, nil)
	if err != nil {
		return nil, err
	}

	return serverResp.body, nil
}
