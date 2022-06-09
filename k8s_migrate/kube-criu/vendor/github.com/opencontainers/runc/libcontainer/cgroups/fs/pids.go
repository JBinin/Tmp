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
	"path/filepath"
	"strconv"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type PidsGroup struct {
}

func (s *PidsGroup) Name() string {
	return "pids"
}

func (s *PidsGroup) Apply(d *cgroupData) error {
	_, err := d.join("pids")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	return nil
}

func (s *PidsGroup) Set(path string, cgroup *configs.Cgroup) error {
	if cgroup.Resources.PidsLimit != 0 {
		// "max" is the fallback value.
		limit := "max"

		if cgroup.Resources.PidsLimit > 0 {
			limit = strconv.FormatInt(cgroup.Resources.PidsLimit, 10)
		}

		if err := writeFile(path, "pids.max", limit); err != nil {
			return err
		}
	}

	return nil
}

func (s *PidsGroup) Remove(d *cgroupData) error {
	return removePath(d.path("pids"))
}

func (s *PidsGroup) GetStats(path string, stats *cgroups.Stats) error {
	current, err := getCgroupParamUint(path, "pids.current")
	if err != nil {
		return fmt.Errorf("failed to parse pids.current - %s", err)
	}

	maxString, err := getCgroupParamString(path, "pids.max")
	if err != nil {
		return fmt.Errorf("failed to parse pids.max - %s", err)
	}

	// Default if pids.max == "max" is 0 -- which represents "no limit".
	var max uint64
	if maxString != "max" {
		max, err = parseUint(maxString, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse pids.max - unable to parse %q as a uint from Cgroup file %q", maxString, filepath.Join(path, "pids.max"))
		}
	}

	stats.PidsStats.Current = current
	stats.PidsStats.Limit = max
	return nil
}
