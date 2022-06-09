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
// +build linux

package fs

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/system"
)

type DevicesGroup struct {
}

func (s *DevicesGroup) Name() string {
	return "devices"
}

func (s *DevicesGroup) Apply(d *cgroupData) error {
	_, err := d.join("devices")
	if err != nil {
		// We will return error even it's `not found` error, devices
		// cgroup is hard requirement for container's security.
		return err
	}
	return nil
}

func (s *DevicesGroup) Set(path string, cgroup *configs.Cgroup) error {
	if system.RunningInUserNS() {
		return nil
	}

	devices := cgroup.Resources.Devices
	if len(devices) > 0 {
		for _, dev := range devices {
			file := "devices.deny"
			if dev.Allow {
				file = "devices.allow"
			}
			if err := writeFile(path, file, dev.CgroupString()); err != nil {
				return err
			}
		}
		return nil
	}
	if cgroup.Resources.AllowAllDevices != nil {
		if *cgroup.Resources.AllowAllDevices == false {
			if err := writeFile(path, "devices.deny", "a"); err != nil {
				return err
			}

			for _, dev := range cgroup.Resources.AllowedDevices {
				if err := writeFile(path, "devices.allow", dev.CgroupString()); err != nil {
					return err
				}
			}
			return nil
		}

		if err := writeFile(path, "devices.allow", "a"); err != nil {
			return err
		}
	}

	for _, dev := range cgroup.Resources.DeniedDevices {
		if err := writeFile(path, "devices.deny", dev.CgroupString()); err != nil {
			return err
		}
	}

	return nil
}

func (s *DevicesGroup) Remove(d *cgroupData) error {
	return removePath(d.path("devices"))
}

func (s *DevicesGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
