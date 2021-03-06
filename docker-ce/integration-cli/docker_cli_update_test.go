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
	"strings"
	"time"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestUpdateRestartPolicy(c *check.C) {
	out, _ := dockerCmd(c, "run", "-d", "--restart=on-failure:3", "busybox", "sh", "-c", "sleep 1 && false")
	timeout := 60 * time.Second
	if testEnv.DaemonPlatform() == "windows" {
		timeout = 180 * time.Second
	}

	id := strings.TrimSpace(string(out))

	// update restart policy to on-failure:5
	dockerCmd(c, "update", "--restart=on-failure:5", id)

	err := waitExited(id, timeout)
	c.Assert(err, checker.IsNil)

	count := inspectField(c, id, "RestartCount")
	c.Assert(count, checker.Equals, "5")

	maximumRetryCount := inspectField(c, id, "HostConfig.RestartPolicy.MaximumRetryCount")
	c.Assert(maximumRetryCount, checker.Equals, "5")
}

func (s *DockerSuite) TestUpdateRestartWithAutoRemoveFlag(c *check.C) {
	out, _ := runSleepingContainer(c, "--rm")
	id := strings.TrimSpace(out)

	// update restart policy for an AutoRemove container
	out, _, err := dockerCmdWithError("update", "--restart=always", id)
	c.Assert(err, checker.NotNil)
	c.Assert(out, checker.Contains, "Restart policy cannot be updated because AutoRemove is enabled for the container")
}
