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
package main

import (
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestStopContainerWithRestartPolicyAlways(c *check.C) {
	dockerCmd(c, "run", "--name", "verifyRestart1", "-d", "--restart=always", "busybox", "false")
	dockerCmd(c, "run", "--name", "verifyRestart2", "-d", "--restart=always", "busybox", "false")

	c.Assert(waitRun("verifyRestart1"), checker.IsNil)
	c.Assert(waitRun("verifyRestart2"), checker.IsNil)

	dockerCmd(c, "stop", "verifyRestart1")
	dockerCmd(c, "stop", "verifyRestart2")
}
