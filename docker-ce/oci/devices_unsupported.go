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
// +build !linux

package oci

import (
	"errors"

	"github.com/opencontainers/runc/libcontainer/configs"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Device transforms a libcontainer configs.Device to a specs.Device object.
// Not implemented
func Device(d *configs.Device) specs.Device { return specs.Device{} }

// DevicesFromPath computes a list of devices and device permissions from paths (pathOnHost and pathInContainer) and cgroup permissions.
// Not implemented
func DevicesFromPath(pathOnHost, pathInContainer, cgroupPermissions string) (devs []specs.Device, devPermissions []specs.DeviceCgroup, err error) {
	return nil, nil, errors.New("oci/devices: unsupported platform")
}
