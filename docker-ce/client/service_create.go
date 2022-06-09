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
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ServiceCreate creates a new Service.
func (cli *Client) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	var headers map[string][]string

	if options.EncodedRegistryAuth != "" {
		headers = map[string][]string{
			"X-Registry-Auth": {options.EncodedRegistryAuth},
		}
	}

	var response types.ServiceCreateResponse
	resp, err := cli.post(ctx, "/services/create", nil, service, headers)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
