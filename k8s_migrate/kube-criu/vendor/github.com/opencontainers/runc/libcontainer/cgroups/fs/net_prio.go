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
)

type NetPrioGroup struct {
}

func (s *NetPrioGroup) Name() string {
	return "net_prio"
}

func (s *NetPrioGroup) Apply(d *cgroupData) error {
	_, err := d.join("net_prio")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *NetPrioGroup) Set(path string, cgroup *configs.Cgroup) error {
	for _, prioMap := range cgroup.Resources.NetPrioIfpriomap {
		if err := writeFile(path, "net_prio.ifpriomap", prioMap.CgroupString()); err != nil {
			return err
		}
	}

	return nil
}

func (s *NetPrioGroup) Remove(d *cgroupData) error {
	return removePath(d.path("net_prio"))
}

func (s *NetPrioGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
