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
// +build !linux static_build

package systemd

import (
	"fmt"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type Manager struct {
	Cgroups *configs.Cgroup
	Paths   map[string]string
}

func UseSystemd() bool {
	return false
}

func (m *Manager) Apply(pid int) error {
	return fmt.Errorf("Systemd not supported")
}

func (m *Manager) GetPids() ([]int, error) {
	return nil, fmt.Errorf("Systemd not supported")
}

func (m *Manager) GetAllPids() ([]int, error) {
	return nil, fmt.Errorf("Systemd not supported")
}

func (m *Manager) Destroy() error {
	return fmt.Errorf("Systemd not supported")
}

func (m *Manager) GetPaths() map[string]string {
	return nil
}

func (m *Manager) GetStats() (*cgroups.Stats, error) {
	return nil, fmt.Errorf("Systemd not supported")
}

func (m *Manager) Set(container *configs.Config) error {
	return fmt.Errorf("Systemd not supported")
}

func (m *Manager) Freeze(state configs.FreezerState) error {
	return fmt.Errorf("Systemd not supported")
}

func Freeze(c *configs.Cgroup, state configs.FreezerState) error {
	return fmt.Errorf("Systemd not supported")
}
