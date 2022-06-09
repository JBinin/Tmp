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
	"fmt"
	"strings"
	"time"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type FreezerGroup struct {
}

func (s *FreezerGroup) Name() string {
	return "freezer"
}

func (s *FreezerGroup) Apply(d *cgroupData) error {
	_, err := d.join("freezer")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *FreezerGroup) Set(path string, cgroup *configs.Cgroup) error {
	switch cgroup.Resources.Freezer {
	case configs.Frozen, configs.Thawed:
		for {
			// In case this loop does not exit because it doesn't get the expected
			// state, let's write again this state, hoping it's going to be properly
			// set this time. Otherwise, this loop could run infinitely, waiting for
			// a state change that would never happen.
			if err := writeFile(path, "freezer.state", string(cgroup.Resources.Freezer)); err != nil {
				return err
			}

			state, err := readFile(path, "freezer.state")
			if err != nil {
				return err
			}
			if strings.TrimSpace(state) == string(cgroup.Resources.Freezer) {
				break
			}

			time.Sleep(1 * time.Millisecond)
		}
	case configs.Undefined:
		return nil
	default:
		return fmt.Errorf("Invalid argument '%s' to freezer.state", string(cgroup.Resources.Freezer))
	}

	return nil
}

func (s *FreezerGroup) Remove(d *cgroupData) error {
	return removePath(d.path("freezer"))
}

func (s *FreezerGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}
