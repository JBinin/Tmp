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

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// NetworkCreate creates a new network in the docker host.
func (cli *Client) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	networkCreateRequest := types.NetworkCreateRequest{
		NetworkCreate: options,
		Name:          name,
	}
	var response types.NetworkCreateResponse
	serverResp, err := cli.post(ctx, "/networks/create", nil, networkCreateRequest, nil)
	if err != nil {
		return response, err
	}

	json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
}
