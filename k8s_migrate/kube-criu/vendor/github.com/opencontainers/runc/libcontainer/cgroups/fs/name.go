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

type NameGroup struct {
	GroupName string
	Join      bool
}

func (s *NameGroup) Name() string {
	return s.GroupName
}

func (s *NameGroup) Apply(d *cgroupData) error {
	if s.Join {
		// ignore errors if the named cgroup does not exist
		d.join(s.GroupName)
	}
	return nil
}

func (s *NameGroup) Set(path string, cgroup *configs.Cgroup) error {
	return nil
}

func (s *NameGroup) Remove(d *cgroupData) error {
	if s.Join {
		removePath(d.path(s.GroupName))
	}
	return nil
}

func (s *NameGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
