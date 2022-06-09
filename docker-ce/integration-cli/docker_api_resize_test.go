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
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestResizeAPIResponse(c *check.C) {
	out, _ := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)

	endpoint := "/containers/" + cleanedContainerID + "/resize?h=40&w=40"
	status, _, err := request.SockRequest("POST", endpoint, nil, daemonHost())
	c.Assert(status, check.Equals, http.StatusOK)
	c.Assert(err, check.IsNil)
}

func (s *DockerSuite) TestResizeAPIHeightWidthNoInt(c *check.C) {
	out, _ := runSleepingContainer(c, "-d")
	cleanedContainerID := strings.TrimSpace(out)

	endpoint := "/containers/" + cleanedContainerID + "/resize?h=foo&w=bar"
	status, _, err := request.SockRequest("POST", endpoint, nil, daemonHost())
	c.Assert(status, check.Equals, http.StatusInternalServerError)
	c.Assert(err, check.IsNil)
}

func (s *DockerSuite) TestResizeAPIResponseWhenContainerNotStarted(c *check.C) {
	out, _ := dockerCmd(c, "run", "-d", "busybox", "true")
	cleanedContainerID := strings.TrimSpace(out)

	// make sure the exited container is not running
	dockerCmd(c, "wait", cleanedContainerID)

	endpoint := "/containers/" + cleanedContainerID + "/resize?h=40&w=40"
	status, body, err := request.SockRequest("POST", endpoint, nil, daemonHost())
	c.Assert(status, check.Equals, http.StatusInternalServerError)
	c.Assert(err, check.IsNil)

	c.Assert(getErrorMessage(c, body), checker.Contains, "is not running", check.Commentf("resize should fail with message 'Container is not running'"))
}
