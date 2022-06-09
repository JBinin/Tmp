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
	"fmt"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SwarmUpdate updates the swarm.
func (cli *Client) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, flags swarm.UpdateFlags) error {
	query := url.Values{}
	query.Set("version", strconv.FormatUint(version.Index, 10))
	query.Set("rotateWorkerToken", fmt.Sprintf("%v", flags.RotateWorkerToken))
	query.Set("rotateManagerToken", fmt.Sprintf("%v", flags.RotateManagerToken))
	query.Set("rotateManagerUnlockKey", fmt.Sprintf("%v", flags.RotateManagerUnlockKey))
	resp, err := cli.post(ctx, "/swarm/update", query, swarm, nil)
	ensureReaderClosed(resp)
	return err
}
