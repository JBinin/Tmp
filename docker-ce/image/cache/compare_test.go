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
package cache

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

// Just to make life easier
func newPortNoError(proto, port string) nat.Port {
	p, _ := nat.NewPort(proto, port)
	return p
}

func TestCompare(t *testing.T) {
	ports1 := make(nat.PortSet)
	ports1[newPortNoError("tcp", "1111")] = struct{}{}
	ports1[newPortNoError("tcp", "2222")] = struct{}{}
	ports2 := make(nat.PortSet)
	ports2[newPortNoError("tcp", "3333")] = struct{}{}
	ports2[newPortNoError("tcp", "4444")] = struct{}{}
	ports3 := make(nat.PortSet)
	ports3[newPortNoError("tcp", "1111")] = struct{}{}
	ports3[newPortNoError("tcp", "2222")] = struct{}{}
	ports3[newPortNoError("tcp", "5555")] = struct{}{}
	volumes1 := make(map[string]struct{})
	volumes1["/test1"] = struct{}{}
	volumes2 := make(map[string]struct{})
	volumes2["/test2"] = struct{}{}
	volumes3 := make(map[string]struct{})
	volumes3["/test1"] = struct{}{}
	volumes3["/test3"] = struct{}{}
	envs1 := []string{"ENV1=value1", "ENV2=value2"}
	envs2 := []string{"ENV1=value1", "ENV3=value3"}
	entrypoint1 := strslice.StrSlice{"/bin/sh", "-c"}
	entrypoint2 := strslice.StrSlice{"/bin/sh", "-d"}
	entrypoint3 := strslice.StrSlice{"/bin/sh", "-c", "echo"}
	cmd1 := strslice.StrSlice{"/bin/sh", "-c"}
	cmd2 := strslice.StrSlice{"/bin/sh", "-d"}
	cmd3 := strslice.StrSlice{"/bin/sh", "-c", "echo"}
	labels1 := map[string]string{"LABEL1": "value1", "LABEL2": "value2"}
	labels2 := map[string]string{"LABEL1": "value1", "LABEL2": "value3"}
	labels3 := map[string]string{"LABEL1": "value1", "LABEL2": "value2", "LABEL3": "value3"}

	sameConfigs := map[*container.Config]*container.Config{
		// Empty config
		&container.Config{}: {},
		// Does not compare hostname, domainname & image
		&container.Config{
			Hostname:   "host1",
			Domainname: "domain1",
			Image:      "image1",
			User:       "user",
		}: {
			Hostname:   "host2",
			Domainname: "domain2",
			Image:      "image2",
			User:       "user",
		},
		// only OpenStdin
		&container.Config{OpenStdin: false}: {OpenStdin: false},
		// only env
		&container.Config{Env: envs1}: {Env: envs1},
		// only cmd
		&container.Config{Cmd: cmd1}: {Cmd: cmd1},
		// only labels
		&container.Config{Labels: labels1}: {Labels: labels1},
		// only exposedPorts
		&container.Config{ExposedPorts: ports1}: {ExposedPorts: ports1},
		// only entrypoints
		&container.Config{Entrypoint: entrypoint1}: {Entrypoint: entrypoint1},
		// only volumes
		&container.Config{Volumes: volumes1}: {Volumes: volumes1},
	}
	differentConfigs := map[*container.Config]*container.Config{
		nil: nil,
		&container.Config{
			Hostname:   "host1",
			Domainname: "domain1",
			Image:      "image1",
			User:       "user1",
		}: {
			Hostname:   "host1",
			Domainname: "domain1",
			Image:      "image1",
			User:       "user2",
		},
		// only OpenStdin
		&container.Config{OpenStdin: false}: {OpenStdin: true},
		&container.Config{OpenStdin: true}:  {OpenStdin: false},
		// only env
		&container.Config{Env: envs1}: {Env: envs2},
		// only cmd
		&container.Config{Cmd: cmd1}: {Cmd: cmd2},
		// not the same number of parts
		&container.Config{Cmd: cmd1}: {Cmd: cmd3},
		// only labels
		&container.Config{Labels: labels1}: {Labels: labels2},
		// not the same number of labels
		&container.Config{Labels: labels1}: {Labels: labels3},
		// only exposedPorts
		&container.Config{ExposedPorts: ports1}: {ExposedPorts: ports2},
		// not the same number of ports
		&container.Config{ExposedPorts: ports1}: {ExposedPorts: ports3},
		// only entrypoints
		&container.Config{Entrypoint: entrypoint1}: {Entrypoint: entrypoint2},
		// not the same number of parts
		&container.Config{Entrypoint: entrypoint1}: {Entrypoint: entrypoint3},
		// only volumes
		&container.Config{Volumes: volumes1}: {Volumes: volumes2},
		// not the same number of labels
		&container.Config{Volumes: volumes1}: {Volumes: volumes3},
	}
	for config1, config2 := range sameConfigs {
		if !compare(config1, config2) {
			t.Fatalf("Compare should be true for [%v] and [%v]", config1, config2)
		}
	}
	for config1, config2 := range differentConfigs {
		if compare(config1, config2) {
			t.Fatalf("Compare should be false for [%v] and [%v]", config1, config2)
		}
	}
}
