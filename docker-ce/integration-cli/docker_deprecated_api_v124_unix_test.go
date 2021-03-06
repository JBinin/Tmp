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
	"fmt"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
)

// #19100 This is a deprecated feature test, it should be removed in Docker 1.12
func (s *DockerNetworkSuite) TestDeprecatedDockerNetworkStartAPIWithHostconfig(c *check.C) {
	netName := "test"
	conName := "foo"
	dockerCmd(c, "network", "create", netName)
	dockerCmd(c, "create", "--name", conName, "busybox", "top")

	config := map[string]interface{}{
		"HostConfig": map[string]interface{}{
			"NetworkMode": netName,
		},
	}
	_, _, err := request.SockRequest("POST", formatV123StartAPIURL("/containers/"+conName+"/start"), config, daemonHost())
	c.Assert(err, checker.IsNil)
	c.Assert(waitRun(conName), checker.IsNil)
	networks := inspectField(c, conName, "NetworkSettings.Networks")
	c.Assert(networks, checker.Contains, netName, check.Commentf(fmt.Sprintf("Should contain '%s' network", netName)))
	c.Assert(networks, checker.Not(checker.Contains), "bridge", check.Commentf("Should not contain 'bridge' network"))
}
