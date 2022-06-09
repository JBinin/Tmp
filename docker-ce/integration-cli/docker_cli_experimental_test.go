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

func (s *DockerSuite) TestExperimentalVersionTrue(c *check.C) {
	testExperimentalInVersion(c, ExperimentalDaemon, "*true")
}

func (s *DockerSuite) TestExperimentalVersionFalse(c *check.C) {
	testExperimentalInVersion(c, NotExperimentalDaemon, "*false")
}

func testExperimentalInVersion(c *check.C, requirement func() bool, expectedValue string) {
	testRequires(c, requirement)
	out, _ := dockerCmd(c, "version")
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Experimental:") {
			c.Assert(line, checker.Matches, expectedValue)
			return
		}
	}

	c.Fatal(`"Experimental" not found in version output`)
}
