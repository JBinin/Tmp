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

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestPause(c *check.C) {
	testRequires(c, IsPausable)
	defer unpauseAllContainers(c)

	name := "testeventpause"
	runSleepingContainer(c, "-d", "--name", name)

	dockerCmd(c, "pause", name)
	pausedContainers := getPausedContainers(c)
	c.Assert(len(pausedContainers), checker.Equals, 1)

	dockerCmd(c, "unpause", name)

	out, _ := dockerCmd(c, "events", "--since=0", "--until", daemonUnixTime(c))
	events := strings.Split(strings.TrimSpace(out), "\n")
	actions := eventActionsByIDAndType(c, events, name, "container")

	c.Assert(actions[len(actions)-2], checker.Equals, "pause")
	c.Assert(actions[len(actions)-1], checker.Equals, "unpause")
}

func (s *DockerSuite) TestPauseMultipleContainers(c *check.C) {
	testRequires(c, IsPausable)
	defer unpauseAllContainers(c)

	containers := []string{
		"testpausewithmorecontainers1",
		"testpausewithmorecontainers2",
	}
	for _, name := range containers {
		runSleepingContainer(c, "-d", "--name", name)
	}
	dockerCmd(c, append([]string{"pause"}, containers...)...)
	pausedContainers := getPausedContainers(c)
	c.Assert(len(pausedContainers), checker.Equals, len(containers))

	dockerCmd(c, append([]string{"unpause"}, containers...)...)

	out, _ := dockerCmd(c, "events", "--since=0", "--until", daemonUnixTime(c))
	events := strings.Split(strings.TrimSpace(out), "\n")

	for _, name := range containers {
		actions := eventActionsByIDAndType(c, events, name, "container")

		c.Assert(actions[len(actions)-2], checker.Equals, "pause")
		c.Assert(actions[len(actions)-1], checker.Equals, "unpause")
	}
}

func (s *DockerSuite) TestPauseFailsOnWindowsServerContainers(c *check.C) {
	testRequires(c, DaemonIsWindows, NotPausable)
	runSleepingContainer(c, "-d", "--name=test")
	out, _, _ := dockerCmdWithError("pause", "test")
	c.Assert(out, checker.Contains, "cannot pause Windows Server Containers")
}
