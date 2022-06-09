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
	"strconv"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerResize changes the size of the tty for a container.
func (cli *Client) ContainerResize(ctx context.Context, containerID string, options types.ResizeOptions) error {
	return cli.resize(ctx, "/containers/"+containerID, options.Height, options.Width)
}

// ContainerExecResize changes the size of the tty for an exec process running inside a container.
func (cli *Client) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	return cli.resize(ctx, "/exec/"+execID, options.Height, options.Width)
}

func (cli *Client) resize(ctx context.Context, basePath string, height, width uint) error {
	query := url.Values{}
	query.Set("h", strconv.Itoa(int(height)))
	query.Set("w", strconv.Itoa(int(width)))

	resp, err := cli.post(ctx, basePath+"/resize", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
