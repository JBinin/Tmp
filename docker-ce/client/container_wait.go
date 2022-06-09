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

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types/container"
)

// ContainerWait pauses execution until a container exits.
// It returns the API status code as response of its readiness.
func (cli *Client) ContainerWait(ctx context.Context, containerID string) (int64, error) {
	resp, err := cli.post(ctx, "/containers/"+containerID+"/wait", nil, nil, nil)
	if err != nil {
		return -1, err
	}
	defer ensureReaderClosed(resp)

	var res container.ContainerWaitOKBody
	if err := json.NewDecoder(resp.body).Decode(&res); err != nil {
		return -1, err
	}

	return res.StatusCode, nil
}
