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

// SwarmGetUnlockKey retrieves the swarm's unlock key.
func (cli *Client) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error) {
	serverResp, err := cli.get(ctx, "/swarm/unlockkey", nil, nil)
	if err != nil {
		return types.SwarmUnlockKeyResponse{}, err
	}

	var response types.SwarmUnlockKeyResponse
	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
}
