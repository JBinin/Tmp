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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

// Test case for #22244
func (s *DockerSuite) TestAuthAPI(c *check.C) {
	testRequires(c, Network)
	config := types.AuthConfig{
		Username: "no-user",
		Password: "no-password",
	}

	expected := "Get https://registry-1.docker.io/v2/: unauthorized: incorrect username or password"
	status, body, err := request.SockRequest("POST", "/auth", config, daemonHost())
	c.Assert(err, check.IsNil)
	c.Assert(status, check.Equals, http.StatusUnauthorized)
	msg := getErrorMessage(c, body)
	c.Assert(msg, checker.Contains, expected, check.Commentf("Expected: %v, got: %v", expected, msg))
}
