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

// ContainerExecCreate creates a new exec configuration to run an exec process.
func (cli *Client) ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) {
	var response types.IDResponse

	if err := cli.NewVersionError("1.25", "env"); len(config.Env) != 0 && err != nil {
		return response, err
	}

	resp, err := cli.post(ctx, "/containers/"+container+"/exec", nil, config, nil)
	if err != nil {
		return response, err
	}
	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}

// ContainerExecStart starts an exec process already created in the docker host.
func (cli *Client) ContainerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	resp, err := cli.post(ctx, "/exec/"+execID+"/start", nil, config, nil)
	ensureReaderClosed(resp)
	return err
}

// ContainerExecAttach attaches a connection to an exec process in the server.
// It returns a types.HijackedConnection with the hijacked connection
// and the a reader to get output. It's up to the called to close
// the hijacked connection by calling types.HijackedResponse.Close.
func (cli *Client) ContainerExecAttach(ctx context.Context, execID string, config types.ExecConfig) (types.HijackedResponse, error) {
	headers := map[string][]string{"Content-Type": {"application/json"}}
	return cli.postHijacked(ctx, "/exec/"+execID+"/start", nil, config, headers)
}

// ContainerExecInspect returns information about a specific exec process on the docker host.
func (cli *Client) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	var response types.ContainerExecInspect
	resp, err := cli.get(ctx, "/exec/"+execID+"/json", nil, nil)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
