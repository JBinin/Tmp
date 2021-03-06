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
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) TestStackRemoveUnknown(c *check.C) {
	d := s.AddDaemon(c, true, true)

	stackArgs := append([]string{"stack", "remove", "UNKNOWN_STACK"})

	out, err := d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "Nothing found in stack: UNKNOWN_STACK\n")
}

func (s *DockerSwarmSuite) TestStackPSUnknown(c *check.C) {
	d := s.AddDaemon(c, true, true)

	stackArgs := append([]string{"stack", "ps", "UNKNOWN_STACK"})

	out, err := d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "Nothing found in stack: UNKNOWN_STACK\n")
}

func (s *DockerSwarmSuite) TestStackServicesUnknown(c *check.C) {
	d := s.AddDaemon(c, true, true)

	stackArgs := append([]string{"stack", "services", "UNKNOWN_STACK"})

	out, err := d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "Nothing found in stack: UNKNOWN_STACK\n")
}

func (s *DockerSwarmSuite) TestStackDeployComposeFile(c *check.C) {
	d := s.AddDaemon(c, true, true)

	testStackName := "testdeploy"
	stackArgs := []string{
		"stack", "deploy",
		"--compose-file", "fixtures/deploy/default.yaml",
		testStackName,
	}
	out, err := d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = d.Cmd("stack", "ls")
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "NAME        SERVICES\n"+"testdeploy  2\n")

	out, err = d.Cmd("stack", "rm", testStackName)
	c.Assert(err, checker.IsNil)
	out, err = d.Cmd("stack", "ls")
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "NAME  SERVICES\n")
}

func (s *DockerSwarmSuite) TestStackDeployWithSecretsTwice(c *check.C) {
	d := s.AddDaemon(c, true, true)

	out, err := d.Cmd("secret", "create", "outside", "fixtures/secrets/default")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	testStackName := "testdeploy"
	stackArgs := []string{
		"stack", "deploy",
		"--compose-file", "fixtures/deploy/secrets.yaml",
		testStackName,
	}
	out, err = d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = d.Cmd("service", "inspect", "--format", "{{ json .Spec.TaskTemplate.ContainerSpec.Secrets }}", "testdeploy_web")
	c.Assert(err, checker.IsNil)

	var refs []swarm.SecretReference
	c.Assert(json.Unmarshal([]byte(out), &refs), checker.IsNil)
	c.Assert(refs, checker.HasLen, 3)

	sort.Sort(sortSecrets(refs))
	c.Assert(refs[0].SecretName, checker.Equals, "outside")
	c.Assert(refs[1].SecretName, checker.Equals, "testdeploy_special")
	c.Assert(refs[1].File.Name, checker.Equals, "special")
	c.Assert(refs[2].SecretName, checker.Equals, "testdeploy_super")
	c.Assert(refs[2].File.Name, checker.Equals, "foo.txt")
	c.Assert(refs[2].File.Mode, checker.Equals, os.FileMode(0400))

	// Deploy again to ensure there are no errors when secret hasn't changed
	out, err = d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil, check.Commentf(out))
}

func (s *DockerSwarmSuite) TestStackRemove(c *check.C) {
	d := s.AddDaemon(c, true, true)

	stackName := "testdeploy"
	stackArgs := []string{
		"stack", "deploy",
		"--compose-file", "fixtures/deploy/remove.yaml",
		stackName,
	}
	out, err := d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = d.Cmd("stack", "ps", stackName)
	c.Assert(err, checker.IsNil)
	c.Assert(strings.Split(strings.TrimSpace(out), "\n"), checker.HasLen, 2)

	out, err = d.Cmd("stack", "rm", stackName)
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "Removing service testdeploy_web")
	c.Assert(out, checker.Contains, "Removing network testdeploy_default")
	c.Assert(out, checker.Contains, "Removing secret testdeploy_special")
}

type sortSecrets []swarm.SecretReference

func (s sortSecrets) Len() int           { return len(s) }
func (s sortSecrets) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortSecrets) Less(i, j int) bool { return s[i].SecretName < s[j].SecretName }

// testDAB is the DAB JSON used for testing.
// TODO: Use template/text and substitute "Image" with the result of
// `docker inspect --format '{{index .RepoDigests 0}}' busybox:latest`
const testDAB = `{
    "Version": "0.1",
    "Services": {
	"srv1": {
	    "Image": "busybox@sha256:e4f93f6ed15a0cdd342f5aae387886fba0ab98af0a102da6276eaf24d6e6ade0",
	    "Command": ["top"]
	},
	"srv2": {
	    "Image": "busybox@sha256:e4f93f6ed15a0cdd342f5aae387886fba0ab98af0a102da6276eaf24d6e6ade0",
	    "Command": ["tail"],
	    "Args": ["-f", "/dev/null"]
	}
    }
}`

func (s *DockerSwarmSuite) TestStackDeployWithDAB(c *check.C) {
	testRequires(c, ExperimentalDaemon)
	// setup
	testStackName := "test"
	testDABFileName := testStackName + ".dab"
	defer os.RemoveAll(testDABFileName)
	err := ioutil.WriteFile(testDABFileName, []byte(testDAB), 0444)
	c.Assert(err, checker.IsNil)
	d := s.AddDaemon(c, true, true)
	// deploy
	stackArgs := []string{
		"stack", "deploy",
		"--bundle-file", testDABFileName,
		testStackName,
	}
	out, err := d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, "Loading bundle from test.dab\n")
	c.Assert(out, checker.Contains, "Creating service test_srv1\n")
	c.Assert(out, checker.Contains, "Creating service test_srv2\n")
	// ls
	stackArgs = []string{"stack", "ls"}
	out, err = d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "NAME  SERVICES\n"+"test  2\n")
	// rm
	stackArgs = []string{"stack", "rm", testStackName}
	out, err = d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, checker.Contains, "Removing service test_srv1\n")
	c.Assert(out, checker.Contains, "Removing service test_srv2\n")
	// ls (empty)
	stackArgs = []string{"stack", "ls"}
	out, err = d.Cmd(stackArgs...)
	c.Assert(err, checker.IsNil)
	c.Assert(out, check.Equals, "NAME  SERVICES\n")
}
