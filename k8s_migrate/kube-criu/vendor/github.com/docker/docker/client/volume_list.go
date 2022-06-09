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

	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

// VolumeList returns the volumes configured in the docker host.
func (cli *Client) VolumeList(ctx context.Context, filter filters.Args) (volumetypes.VolumesListOKBody, error) {
	var volumes volumetypes.VolumesListOKBody
	query := url.Values{}

	if filter.Len() > 0 {
		filterJSON, err := filters.ToParamWithVersion(cli.version, filter)
		if err != nil {
			return volumes, err
		}
		query.Set("filters", filterJSON)
	}
	resp, err := cli.get(ctx, "/volumes", query, nil)
	if err != nil {
		return volumes, err
	}

	err = json.NewDecoder(resp.body).Decode(&volumes)
	ensureReaderClosed(resp)
	return volumes, err
}
