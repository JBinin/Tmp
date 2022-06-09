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
package controlapi

import (
	"errors"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/manager/state/raft"
	"github.com/docker/swarmkit/manager/state/store"
)

var (
	errNotImplemented  = errors.New("not implemented")
	errInvalidArgument = errors.New("invalid argument")
)

// Server is the Cluster API gRPC server.
type Server struct {
	store  *store.MemoryStore
	raft   *raft.Node
	rootCA *ca.RootCA
	pg     plugingetter.PluginGetter
}

// NewServer creates a Cluster API server.
func NewServer(store *store.MemoryStore, raft *raft.Node, rootCA *ca.RootCA, pg plugingetter.PluginGetter) *Server {
	return &Server{
		store:  store,
		raft:   raft,
		rootCA: rootCA,
		pg:     pg,
	}
}
