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
	"time"

	timetypes "github.com/docker/docker/api/types/time"
	"golang.org/x/net/context"
)

// ContainerRestart stops and starts a container again.
// It makes the daemon to wait for the container to be up again for
// a specific amount of time, given the timeout.
func (cli *Client) ContainerRestart(ctx context.Context, containerID string, timeout *time.Duration) error {
	query := url.Values{}
	if timeout != nil {
		query.Set("t", timetypes.DurationToSecondsString(*timeout))
	}
	resp, err := cli.post(ctx, "/containers/"+containerID+"/restart", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
