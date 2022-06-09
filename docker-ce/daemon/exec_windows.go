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
	"github.com/docker/docker/daemon/exec"
	"github.com/docker/docker/libcontainerd"
)

func execSetPlatformOpt(c *container.Container, ec *exec.Config, p *libcontainerd.Process) error {
	// Process arguments need to be escaped before sending to OCI.
	p.Args = escapeArgs(p.Args)
	p.User.Username = ec.User
	return nil
}
