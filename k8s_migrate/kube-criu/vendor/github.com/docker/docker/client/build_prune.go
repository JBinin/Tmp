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
	"fmt"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// BuildCachePrune requests the daemon to delete unused cache data
func (cli *Client) BuildCachePrune(ctx context.Context) (*types.BuildCachePruneReport, error) {
	if err := cli.NewVersionError("1.31", "build prune"); err != nil {
		return nil, err
	}

	report := types.BuildCachePruneReport{}

	serverResp, err := cli.post(ctx, "/build/prune", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil {
		return nil, fmt.Errorf("Error retrieving disk usage: %v", err)
	}

	return &report, nil
}
