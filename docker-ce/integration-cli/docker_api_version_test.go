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
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestGetVersion(c *check.C) {
	status, body, err := request.SockRequest("GET", "/version", nil, daemonHost())
	c.Assert(status, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)

	var v types.Version

	c.Assert(json.Unmarshal(body, &v), checker.IsNil)

	c.Assert(v.Version, checker.Equals, dockerversion.Version, check.Commentf("Version mismatch"))
}
