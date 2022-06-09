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

type PerfEventGroup struct {
}

func (s *PerfEventGroup) Name() string {
	return "perf_event"
}

func (s *PerfEventGroup) Apply(d *cgroupData) error {
	// we just want to join this group even though we don't set anything
	if _, err := d.join("perf_event"); err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *PerfEventGroup) Set(path string, cgroup *configs.Cgroup) error {
	return nil
}

func (s *PerfEventGroup) Remove(d *cgroupData) error {
	return removePath(d.path("perf_event"))
}

func (s *PerfEventGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
