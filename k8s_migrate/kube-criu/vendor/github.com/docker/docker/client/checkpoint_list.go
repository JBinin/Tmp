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
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// CheckpointList returns the checkpoints of the given container in the docker host
func (cli *Client) CheckpointList(ctx context.Context, container string, options types.CheckpointListOptions) ([]types.Checkpoint, error) {
	var checkpoints []types.Checkpoint

	query := url.Values{}
	if options.CheckpointDir != "" {
		query.Set("dir", options.CheckpointDir)
	}

	resp, err := cli.get(ctx, "/containers/"+container+"/checkpoints", query, nil)
	if err != nil {
		if resp.statusCode == http.StatusNotFound {
			return checkpoints, containerNotFoundError{container}
		}
		return checkpoints, err
	}

	err = json.NewDecoder(resp.body).Decode(&checkpoints)
	ensureReaderClosed(resp)
	return checkpoints, err
}
