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
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// ImageBuild sends request to the daemon to build images.
// The Body in the response implement an io.ReadCloser and it's up to the caller to
// close it.
func (cli *Client) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	query, err := cli.imageBuildOptionsToQuery(options)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	headers := http.Header(make(map[string][]string))
	buf, err := json.Marshal(options.AuthConfigs)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	headers.Add("X-Registry-Config", base64.URLEncoding.EncodeToString(buf))
	headers.Set("Content-Type", "application/x-tar")

	serverResp, err := cli.postRaw(ctx, "/build", query, buildContext, headers)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	osType := getDockerOS(serverResp.header.Get("Server"))

	return types.ImageBuildResponse{
		Body:   serverResp.body,
		OSType: osType,
	}, nil
}

func (cli *Client) imageBuildOptionsToQuery(options types.ImageBuildOptions) (url.Values, error) {
	query := url.Values{
		"t":           options.Tags,
		"securityopt": options.SecurityOpt,
		"extrahosts":  options.ExtraHosts,
	}
	if options.SuppressOutput {
		query.Set("q", "1")
	}
	if options.RemoteContext != "" {
		query.Set("remote", options.RemoteContext)
	}
	if options.NoCache {
		query.Set("nocache", "1")
	}
	if options.Remove {
		query.Set("rm", "1")
	} else {
		query.Set("rm", "0")
	}

	if options.ForceRemove {
		query.Set("forcerm", "1")
	}

	if options.PullParent {
		query.Set("pull", "1")
	}

	if options.Squash {
		if err := cli.NewVersionError("1.25", "squash"); err != nil {
			return query, err
		}
		query.Set("squash", "1")
	}

	if !container.Isolation.IsDefault(options.Isolation) {
		query.Set("isolation", string(options.Isolation))
	}

	query.Set("cpusetcpus", options.CPUSetCPUs)
	query.Set("networkmode", options.NetworkMode)
	query.Set("cpusetmems", options.CPUSetMems)
	query.Set("cpushares", strconv.FormatInt(options.CPUShares, 10))
	query.Set("cpuquota", strconv.FormatInt(options.CPUQuota, 10))
	query.Set("cpuperiod", strconv.FormatInt(options.CPUPeriod, 10))
	query.Set("memory", strconv.FormatInt(options.Memory, 10))
	query.Set("memswap", strconv.FormatInt(options.MemorySwap, 10))
	query.Set("cgroupparent", options.CgroupParent)
	query.Set("shmsize", strconv.FormatInt(options.ShmSize, 10))
	query.Set("dockerfile", options.Dockerfile)
	query.Set("target", options.Target)

	ulimitsJSON, err := json.Marshal(options.Ulimits)
	if err != nil {
		return query, err
	}
	query.Set("ulimits", string(ulimitsJSON))

	buildArgsJSON, err := json.Marshal(options.BuildArgs)
	if err != nil {
		return query, err
	}
	query.Set("buildargs", string(buildArgsJSON))

	labelsJSON, err := json.Marshal(options.Labels)
	if err != nil {
		return query, err
	}
	query.Set("labels", string(labelsJSON))

	cacheFromJSON, err := json.Marshal(options.CacheFrom)
	if err != nil {
		return query, err
	}
	query.Set("cachefrom", string(cacheFromJSON))
	if options.SessionID != "" {
		query.Set("session", options.SessionID)
	}

	return query, nil
}
