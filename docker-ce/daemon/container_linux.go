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
//+build !windows

package daemon

import (
	"github.com/docker/docker/container"
)

func (daemon *Daemon) saveApparmorConfig(container *container.Container) error {
	container.AppArmorProfile = "" //we don't care about the previous value.

	if !daemon.apparmorEnabled {
		return nil // if apparmor is disabled there is nothing to do here.
	}

	if err := parseSecurityOpt(container, container.HostConfig); err != nil {
		return err
	}

	if !container.HostConfig.Privileged {
		if container.AppArmorProfile == "" {
			container.AppArmorProfile = defaultApparmorProfile
		}

	} else {
		container.AppArmorProfile = "unconfined"
	}
	return nil
}
