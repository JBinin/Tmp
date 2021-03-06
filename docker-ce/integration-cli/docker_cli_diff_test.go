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

// ensure that an added file shows up in docker diff
func (s *DockerSuite) TestDiffFilenameShownInOutput(c *check.C) {
	containerCmd := `mkdir /foo; echo xyzzy > /foo/bar`
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", containerCmd)

	// Wait for it to exit as cannot diff a running container on Windows, and
	// it will take a few seconds to exit. Also there's no way in Windows to
	// differentiate between an Add or a Modify, and all files are under
	// a "Files/" prefix.
	containerID := strings.TrimSpace(out)
	lookingFor := "A /foo/bar"
	if testEnv.DaemonPlatform() == "windows" {
		err := waitExited(containerID, 60*time.Second)
		c.Assert(err, check.IsNil)
		lookingFor = "C Files/foo/bar"
	}

	cleanCID := strings.TrimSpace(out)
	out, _ = dockerCmd(c, "diff", cleanCID)

	found := false
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, lookingFor) {
			found = true
			break
		}
	}
	c.Assert(found, checker.True)
}

// test to ensure GH #3840 doesn't occur any more
func (s *DockerSuite) TestDiffEnsureInitLayerFilesAreIgnored(c *check.C) {
	testRequires(c, DaemonIsLinux)
	// this is a list of files which shouldn't show up in `docker diff`
	initLayerFiles := []string{"/etc/resolv.conf", "/etc/hostname", "/etc/hosts", "/.dockerenv"}
	containerCount := 5

	// we might not run into this problem from the first run, so start a few containers
	for i := 0; i < containerCount; i++ {
		containerCmd := `echo foo > /root/bar`
		out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", containerCmd)

		cleanCID := strings.TrimSpace(out)
		out, _ = dockerCmd(c, "diff", cleanCID)

		for _, filename := range initLayerFiles {
			c.Assert(out, checker.Not(checker.Contains), filename)
		}
	}
}

func (s *DockerSuite) TestDiffEnsureDefaultDevs(c *check.C) {
	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sleep", "0")

	cleanCID := strings.TrimSpace(out)
	out, _ = dockerCmd(c, "diff", cleanCID)

	expected := map[string]bool{
		"C /dev":         true,
		"A /dev/full":    true, // busybox
		"C /dev/ptmx":    true, // libcontainer
		"A /dev/mqueue":  true,
		"A /dev/kmsg":    true,
		"A /dev/fd":      true,
		"A /dev/ptmx":    true,
		"A /dev/null":    true,
		"A /dev/random":  true,
		"A /dev/stdout":  true,
		"A /dev/stderr":  true,
		"A /dev/tty1":    true,
		"A /dev/stdin":   true,
		"A /dev/tty":     true,
		"A /dev/urandom": true,
		"A /dev/zero":    true,
	}

	for _, line := range strings.Split(out, "\n") {
		c.Assert(line == "" || expected[line], checker.True, check.Commentf(line))
	}
}

// https://github.com/docker/docker/pull/14381#discussion_r33859347
func (s *DockerSuite) TestDiffEmptyArgClientError(c *check.C) {
	out, _, err := dockerCmdWithError("diff", "")
	c.Assert(err, checker.NotNil)
	c.Assert(strings.TrimSpace(out), checker.Contains, "Container name cannot be empty")
}
