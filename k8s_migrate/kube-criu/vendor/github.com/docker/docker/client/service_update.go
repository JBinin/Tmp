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
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// ServiceUpdate updates a Service.
func (cli *Client) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) {
	var (
		query   = url.Values{}
		distErr error
	)

	headers := map[string][]string{
		"version": {cli.version},
	}

	if options.EncodedRegistryAuth != "" {
		headers["X-Registry-Auth"] = []string{options.EncodedRegistryAuth}
	}

	if options.RegistryAuthFrom != "" {
		query.Set("registryAuthFrom", options.RegistryAuthFrom)
	}

	if options.Rollback != "" {
		query.Set("rollback", options.Rollback)
	}

	query.Set("version", strconv.FormatUint(version.Index, 10))

	if err := validateServiceSpec(service); err != nil {
		return types.ServiceUpdateResponse{}, err
	}

	var imgPlatforms []swarm.Platform
	// ensure that the image is tagged
	if service.TaskTemplate.ContainerSpec != nil {
		if taggedImg := imageWithTagString(service.TaskTemplate.ContainerSpec.Image); taggedImg != "" {
			service.TaskTemplate.ContainerSpec.Image = taggedImg
		}
		if options.QueryRegistry {
			var img string
			img, imgPlatforms, distErr = imageDigestAndPlatforms(ctx, cli, service.TaskTemplate.ContainerSpec.Image, options.EncodedRegistryAuth)
			if img != "" {
				service.TaskTemplate.ContainerSpec.Image = img
			}
		}
	}

	// ensure that the image is tagged
	if service.TaskTemplate.PluginSpec != nil {
		if taggedImg := imageWithTagString(service.TaskTemplate.PluginSpec.Remote); taggedImg != "" {
			service.TaskTemplate.PluginSpec.Remote = taggedImg
		}
		if options.QueryRegistry {
			var img string
			img, imgPlatforms, distErr = imageDigestAndPlatforms(ctx, cli, service.TaskTemplate.PluginSpec.Remote, options.EncodedRegistryAuth)
			if img != "" {
				service.TaskTemplate.PluginSpec.Remote = img
			}
		}
	}

	if service.TaskTemplate.Placement == nil && len(imgPlatforms) > 0 {
		service.TaskTemplate.Placement = &swarm.Placement{}
	}
	if len(imgPlatforms) > 0 {
		service.TaskTemplate.Placement.Platforms = imgPlatforms
	}

	var response types.ServiceUpdateResponse
	resp, err := cli.post(ctx, "/services/"+serviceID+"/update", query, service, headers)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)

	if distErr != nil {
		response.Warnings = append(response.Warnings, digestWarning(service.TaskTemplate.ContainerSpec.Image))
	}

	ensureReaderClosed(resp)
	return response, err
}
