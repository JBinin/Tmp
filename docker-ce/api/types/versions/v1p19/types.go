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
// Package v1p19 provides specific API types for the API version 1, patch 19.
package v1p19

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/versions/v1p20"
	"github.com/docker/go-connections/nat"
)

// ContainerJSON is a backcompatibility struct for APIs prior to 1.20.
// Note this is not used by the Windows daemon.
type ContainerJSON struct {
	*types.ContainerJSONBase
	Volumes         map[string]string
	VolumesRW       map[string]bool
	Config          *ContainerConfig
	NetworkSettings *v1p20.NetworkSettings
}

// ContainerConfig is a backcompatibility struct for APIs prior to 1.20.
type ContainerConfig struct {
	*container.Config

	MacAddress      string
	NetworkDisabled bool
	ExposedPorts    map[nat.Port]struct{}

	// backward compatibility, they now live in HostConfig
	VolumeDriver string
	Memory       int64
	MemorySwap   int64
	CPUShares    int64  `json:"CpuShares"`
	CPUSet       string `json:"Cpuset"`
}
