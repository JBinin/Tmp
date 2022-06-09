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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types/swarm"

	"golang.org/x/net/context"
)

// TaskInspectWithRaw returns the task information and its raw representation..
func (cli *Client) TaskInspectWithRaw(ctx context.Context, taskID string) (swarm.Task, []byte, error) {
	serverResp, err := cli.get(ctx, "/tasks/"+taskID, nil, nil)
	if err != nil {
		if serverResp.statusCode == http.StatusNotFound {
			return swarm.Task{}, nil, taskNotFoundError{taskID}
		}
		return swarm.Task{}, nil, err
	}
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil {
		return swarm.Task{}, nil, err
	}

	var response swarm.Task
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
}
