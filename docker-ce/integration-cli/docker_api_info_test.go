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
	"net/http"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestInfoAPI(c *check.C) {
	endpoint := "/info"

	status, body, err := request.SockRequest("GET", endpoint, nil, daemonHost())
	c.Assert(status, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)

	// always shown fields
	stringsToCheck := []string{
		"ID",
		"Containers",
		"ContainersRunning",
		"ContainersPaused",
		"ContainersStopped",
		"Images",
		"LoggingDriver",
		"OperatingSystem",
		"NCPU",
		"OSType",
		"Architecture",
		"MemTotal",
		"KernelVersion",
		"Driver",
		"ServerVersion",
		"SecurityOptions"}

	out := string(body)
	for _, linePrefix := range stringsToCheck {
		c.Assert(out, checker.Contains, linePrefix)
	}
}

func (s *DockerSuite) TestInfoAPIVersioned(c *check.C) {
	testRequires(c, DaemonIsLinux) // Windows only supports 1.25 or later
	endpoint := "/v1.20/info"

	status, body, err := request.SockRequest("GET", endpoint, nil, daemonHost())
	c.Assert(status, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)

	out := string(body)
	c.Assert(out, checker.Contains, "ExecutionDriver")
	c.Assert(out, checker.Contains, "not supported")
}
