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

// ImageSave retrieves one or more images from the docker host as an io.ReadCloser.
// It's up to the caller to store the images and close the stream.
func (cli *Client) ImageSave(ctx context.Context, imageIDs []string) (io.ReadCloser, error) {
	query := url.Values{
		"names": imageIDs,
	}

	resp, err := cli.get(ctx, "/images/get", query, nil)
	if err != nil {
		return nil, err
	}
	return resp.body, nil
}
