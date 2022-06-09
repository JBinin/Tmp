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
package libcontainer

import (
	"github.com/opencontainers/runc/libcontainer/configs"
)

type Factory interface {
	// Creates a new container with the given id and starts the initial process inside it.
	// id must be a string containing only letters, digits and underscores and must contain
	// between 1 and 1024 characters, inclusive.
	//
	// The id must not already be in use by an existing container. Containers created using
	// a factory with the same path (and filesystem) must have distinct ids.
	//
	// Returns the new container with a running process.
	//
	// errors:
	// IdInUse - id is already in use by a container
	// InvalidIdFormat - id has incorrect format
	// ConfigInvalid - config is invalid
	// Systemerror - System error
	//
	// On error, any partially created container parts are cleaned up (the operation is atomic).
	Create(id string, config *configs.Config) (Container, error)

	// Load takes an ID for an existing container and returns the container information
	// from the state.  This presents a read only view of the container.
	//
	// errors:
	// Path does not exist
	// System error
	Load(id string) (Container, error)

	// StartInitialization is an internal API to libcontainer used during the reexec of the
	// container.
	//
	// Errors:
	// Pipe connection error
	// System error
	StartInitialization() error

	// Type returns info string about factory type (e.g. lxc, libcontainer...)
	Type() string
}
