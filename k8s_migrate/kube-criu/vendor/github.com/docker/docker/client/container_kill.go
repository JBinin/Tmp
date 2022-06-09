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
	"net/url"

	"golang.org/x/net/context"
)

// ContainerKill terminates the container process but does not remove the container from the docker host.
func (cli *Client) ContainerKill(ctx context.Context, containerID, signal string) error {
	query := url.Values{}
	query.Set("signal", signal)

	resp, err := cli.post(ctx, "/containers/"+containerID+"/kill", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
