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

// SecretList returns the list of secrets.
func (cli *Client) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
	if err := cli.NewVersionError("1.25", "secret list"); err != nil {
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

	resp, err := cli.get(ctx, "/secrets", query, nil)
	if err != nil {
		return nil, err
	}

	var secrets []swarm.Secret
	err = json.NewDecoder(resp.body).Decode(&secrets)
	ensureReaderClosed(resp)
	return secrets, err
}
