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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ConfigList returns the list of configs.
func (cli *Client) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) {
	if err := cli.NewVersionError("1.30", "config list"); err != nil {
		return nil, err
	}
	query := url.Values{}

	if options.Filters.Len() > 0 {
		filterJSON, err := filters.ToParam(options.Filters)
		if err != nil {
			return nil, err
		}

		query.Set("filters", filterJSON)
	}

	resp, err := cli.get(ctx, "/configs", query, nil)
	if err != nil {
		return nil, err
	}

	var configs []swarm.Config
	err = json.NewDecoder(resp.body).Decode(&configs)
	ensureReaderClosed(resp)
	return configs, err
}
