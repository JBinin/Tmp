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
	"encoding/json"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ImageRemove removes an image from the docker host.
func (cli *Client) ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	query := url.Values{}

	if options.Force {
		query.Set("force", "1")
	}
	if !options.PruneChildren {
		query.Set("noprune", "1")
	}

	resp, err := cli.delete(ctx, "/images/"+imageID, query, nil)
	if err != nil {
		return nil, err
	}

	var dels []types.ImageDeleteResponseItem
	err = json.NewDecoder(resp.body).Decode(&dels)
	ensureReaderClosed(resp)
	return dels, err
}
