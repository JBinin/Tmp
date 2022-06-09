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
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

// VolumeCreate creates a volume in the docker host.
func (cli *Client) VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error) {
	var volume types.Volume
	resp, err := cli.post(ctx, "/volumes/create", nil, options, nil)
	if err != nil {
		return volume, err
	}
	err = json.NewDecoder(resp.body).Decode(&volume)
	ensureReaderClosed(resp)
	return volume, err
}
