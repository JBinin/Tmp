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
package agent

import (
	"github.com/boltdb/bolt"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/connectionbroker"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
)

// Config provides values for an Agent.
type Config struct {
	// Hostname the name of host for agent instance.
	Hostname string

	// ConnBroker provides a connection broker for retrieving gRPC
	// connections to managers.
	ConnBroker *connectionbroker.Broker

	// Executor specifies the executor to use for the agent.
	Executor exec.Executor

	// DB used for task storage. Must be open for the lifetime of the agent.
	DB *bolt.DB

	// NotifyNodeChange channel receives new node changes from session messages.
	NotifyNodeChange chan<- *api.Node

	// Credentials is credentials for grpc connection to manager.
	Credentials credentials.TransportCredentials
}

func (c *Config) validate() error {
	if c.Credentials == nil {
		return errors.New("agent: Credentials is required")
	}

	if c.Executor == nil {
		return errors.New("agent: executor required")
	}

	if c.DB == nil {
		return errors.New("agent: database required")
	}

	return nil
}
