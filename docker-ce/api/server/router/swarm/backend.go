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
package swarm

import (
	basictypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	types "github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// Backend abstracts a swarm manager.
type Backend interface {
	Init(req types.InitRequest) (string, error)
	Join(req types.JoinRequest) error
	Leave(force bool) error
	Inspect() (types.Swarm, error)
	Update(uint64, types.Spec, types.UpdateFlags) error
	GetUnlockKey() (string, error)
	UnlockSwarm(req types.UnlockRequest) error
	GetServices(basictypes.ServiceListOptions) ([]types.Service, error)
	GetService(string) (types.Service, error)
	CreateService(types.ServiceSpec, string) (*basictypes.ServiceCreateResponse, error)
	UpdateService(string, uint64, types.ServiceSpec, string, string) (*basictypes.ServiceUpdateResponse, error)
	RemoveService(string) error
	ServiceLogs(context.Context, string, *backend.ContainerLogsConfig, chan struct{}) error
	GetNodes(basictypes.NodeListOptions) ([]types.Node, error)
	GetNode(string) (types.Node, error)
	UpdateNode(string, uint64, types.NodeSpec) error
	RemoveNode(string, bool) error
	GetTasks(basictypes.TaskListOptions) ([]types.Task, error)
	GetTask(string) (types.Task, error)
	GetSecrets(opts basictypes.SecretListOptions) ([]types.Secret, error)
	CreateSecret(s types.SecretSpec) (string, error)
	RemoveSecret(id string) error
	GetSecret(id string) (types.Secret, error)
	UpdateSecret(id string, version uint64, spec types.SecretSpec) error
}
