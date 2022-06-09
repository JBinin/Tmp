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
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// Info returns information about the docker server.
func (cli *Client) Info(ctx context.Context) (types.Info, error) {
	var info types.Info
	serverResp, err := cli.get(ctx, "/info", url.Values{}, nil)
	if err != nil {
		return info, err
	}
	defer ensureReaderClosed(serverResp)

	if err := json.NewDecoder(serverResp.body).Decode(&info); err != nil {
		return info, fmt.Errorf("Error reading remote info: %v", err)
	}

	return info, nil
}
