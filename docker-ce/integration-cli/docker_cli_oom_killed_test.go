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
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestInspectOomKilledTrue(c *check.C) {
	testRequires(c, DaemonIsLinux, memoryLimitSupport, swapMemorySupport)

	name := "testoomkilled"
	_, exitCode, _ := dockerCmdWithError("run", "--name", name, "--memory", "32MB", "busybox", "sh", "-c", "x=a; while true; do x=$x$x$x$x; done")

	c.Assert(exitCode, checker.Equals, 137, check.Commentf("OOM exit should be 137"))

	oomKilled := inspectField(c, name, "State.OOMKilled")
	c.Assert(oomKilled, checker.Equals, "true")
}

func (s *DockerSuite) TestInspectOomKilledFalse(c *check.C) {
	testRequires(c, DaemonIsLinux, memoryLimitSupport, swapMemorySupport)

	name := "testoomkilled"
	dockerCmd(c, "run", "--name", name, "--memory", "32MB", "busybox", "sh", "-c", "echo hello world")

	oomKilled := inspectField(c, name, "State.OOMKilled")
	c.Assert(oomKilled, checker.Equals, "false")
}
