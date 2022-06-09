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

	"github.com/docker/docker/api/types"
)

// ImageLoad loads an image in the docker host from the client host.
// It's up to the caller to close the io.ReadCloser in the
// ImageLoadResponse returned by this function.
func (cli *Client) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	v := url.Values{}
	v.Set("quiet", "0")
	if quiet {
		v.Set("quiet", "1")
	}
	headers := map[string][]string{"Content-Type": {"application/x-tar"}}
	resp, err := cli.postRaw(ctx, "/images/load", v, input, headers)
	if err != nil {
		return types.ImageLoadResponse{}, err
	}
	return types.ImageLoadResponse{
		Body: resp.body,
		JSON: resp.header.Get("Content-Type") == "application/json",
	}, nil
}
