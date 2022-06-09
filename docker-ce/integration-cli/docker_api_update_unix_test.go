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
// +build !windows

package main

import (
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestAPIUpdateContainer(c *check.C) {
	testRequires(c, DaemonIsLinux)
	testRequires(c, memoryLimitSupport)
	testRequires(c, swapMemorySupport)

	name := "apiUpdateContainer"
	hostConfig := map[string]interface{}{
		"Memory":     314572800,
		"MemorySwap": 524288000,
	}
	dockerCmd(c, "run", "-d", "--name", name, "-m", "200M", "busybox", "top")
	_, _, err := request.SockRequest("POST", "/containers/"+name+"/update", hostConfig, daemonHost())
	c.Assert(err, check.IsNil)

	c.Assert(inspectField(c, name, "HostConfig.Memory"), checker.Equals, "314572800")
	file := "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	out, _ := dockerCmd(c, "exec", name, "cat", file)
	c.Assert(strings.TrimSpace(out), checker.Equals, "314572800")

	c.Assert(inspectField(c, name, "HostConfig.MemorySwap"), checker.Equals, "524288000")
	file = "/sys/fs/cgroup/memory/memory.memsw.limit_in_bytes"
	out, _ = dockerCmd(c, "exec", name, "cat", file)
	c.Assert(strings.TrimSpace(out), checker.Equals, "524288000")
}
