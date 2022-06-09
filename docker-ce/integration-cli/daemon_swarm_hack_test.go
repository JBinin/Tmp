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
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/go-check/check"
)

func (s *DockerSwarmSuite) getDaemon(c *check.C, nodeID string) *daemon.Swarm {
	s.daemonsLock.Lock()
	defer s.daemonsLock.Unlock()
	for _, d := range s.daemons {
		if d.NodeID == nodeID {
			return d
		}
	}
	c.Fatalf("could not find node with id: %s", nodeID)
	return nil
}

// nodeCmd executes a command on a given node via the normal docker socket
func (s *DockerSwarmSuite) nodeCmd(c *check.C, id string, args ...string) (string, error) {
	return s.getDaemon(c, id).Cmd(args...)
}
