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

	"github.com/docker/docker/api/types/versions"
	"golang.org/x/net/context"
)

// VolumeRemove removes a volume from the docker host.
func (cli *Client) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	query := url.Values{}
	if versions.GreaterThanOrEqualTo(cli.version, "1.25") {
		if force {
			query.Set("force", "1")
		}
	}
	resp, err := cli.delete(ctx, "/volumes/"+volumeID, query, nil)
	ensureReaderClosed(resp)
	return err
}
