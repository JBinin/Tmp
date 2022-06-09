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
	"io/ioutil"
	"os"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestLinksEtcHostsContentMatch(c *check.C) {
	// In a _unix file as using Unix specific files, and must be on the
	// same host as the daemon.
	testRequires(c, SameHostDaemon, NotUserNamespace)

	out, _ := dockerCmd(c, "run", "--net=host", "busybox", "cat", "/etc/hosts")
	hosts, err := ioutil.ReadFile("/etc/hosts")
	if os.IsNotExist(err) {
		c.Skip("/etc/hosts does not exist, skip this test")
	}

	c.Assert(out, checker.Equals, string(hosts), check.Commentf("container: %s\n\nhost:%s", out, hosts))

}
