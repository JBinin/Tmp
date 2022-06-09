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
package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/caps"
	"github.com/docker/docker/daemon/exec"
	"github.com/docker/docker/libcontainerd"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func execSetPlatformOpt(c *container.Container, ec *exec.Config, p *libcontainerd.Process) error {
	if len(ec.User) > 0 {
		uid, gid, additionalGids, err := getUser(c, ec.User)
		if err != nil {
			return err
		}
		p.User = &specs.User{
			UID:            uid,
			GID:            gid,
			AdditionalGids: additionalGids,
		}
	}
	if ec.Privileged {
		p.Capabilities = caps.GetAllCapabilities()
	}
	return nil
}
