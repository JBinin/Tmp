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

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ContainerRemove kills and removes a container from the docker host.
func (cli *Client) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	query := url.Values{}
	if options.RemoveVolumes {
		query.Set("v", "1")
	}
	if options.RemoveLinks {
		query.Set("link", "1")
	}

	if options.Force {
		query.Set("force", "1")
	}

	resp, err := cli.delete(ctx, "/containers/"+containerID, query, nil)
	ensureReaderClosed(resp)
	return err
}
