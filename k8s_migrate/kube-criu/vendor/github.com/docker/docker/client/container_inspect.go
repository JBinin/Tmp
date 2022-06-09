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
	"net/url"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerInspect returns the container information.
func (cli *Client) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	serverResp, err := cli.get(ctx, "/containers/"+containerID+"/json", nil, nil)
	if err != nil {
		if serverResp.statusCode == http.StatusNotFound {
			return types.ContainerJSON{}, containerNotFoundError{containerID}
		}
		return types.ContainerJSON{}, err
	}

	var response types.ContainerJSON
	err = json.NewDecoder(serverResp.body).Decode(&response)
	ensureReaderClosed(serverResp)
	return response, err
}

// ContainerInspectWithRaw returns the container information and its raw representation.
func (cli *Client) ContainerInspectWithRaw(ctx context.Context, containerID string, getSize bool) (types.ContainerJSON, []byte, error) {
	query := url.Values{}
	if getSize {
		query.Set("size", "1")
	}
	serverResp, err := cli.get(ctx, "/containers/"+containerID+"/json", query, nil)
	if err != nil {
		if serverResp.statusCode == http.StatusNotFound {
			return types.ContainerJSON{}, nil, containerNotFoundError{containerID}
		}
		return types.ContainerJSON{}, nil, err
	}
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil {
		return types.ContainerJSON{}, nil, err
	}

	var response types.ContainerJSON
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
}
