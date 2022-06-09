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

	registrytypes "github.com/docker/docker/api/types/registry"
	"golang.org/x/net/context"
)

// DistributionInspect returns the image digest with full Manifest
func (cli *Client) DistributionInspect(ctx context.Context, image, encodedRegistryAuth string) (registrytypes.DistributionInspect, error) {
	// Contact the registry to retrieve digest and platform information
	var distributionInspect registrytypes.DistributionInspect

	if err := cli.NewVersionError("1.30", "distribution inspect"); err != nil {
		return distributionInspect, err
	}
	var headers map[string][]string

	if encodedRegistryAuth != "" {
		headers = map[string][]string{
			"X-Registry-Auth": {encodedRegistryAuth},
		}
	}

	resp, err := cli.get(ctx, "/distribution/"+image+"/json", url.Values{}, headers)
	if err != nil {
		return distributionInspect, err
	}

	err = json.NewDecoder(resp.body).Decode(&distributionInspect)
	ensureReaderClosed(resp)
	return distributionInspect, err
}
