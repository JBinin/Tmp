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
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

// ContainerTop shows process information from within a container.
func (cli *Client) ContainerTop(ctx context.Context, containerID string, arguments []string) (container.ContainerTopOKBody, error) {
	var response container.ContainerTopOKBody
	query := url.Values{}
	if len(arguments) > 0 {
		query.Set("ps_args", strings.Join(arguments, " "))
	}

	resp, err := cli.get(ctx, "/containers/"+containerID+"/top", query, nil)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
