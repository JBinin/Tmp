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
package libcontainerd

import (
	"github.com/Microsoft/hcsshim"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// Process contains information to start a specific application inside the container.
type Process specs.Process

// Summary contains a ProcessList item from HCS to support `top`
type Summary hcsshim.ProcessListItem

// StateInfo contains description about the new state container has entered.
type StateInfo struct {
	CommonStateInfo

	// Platform specific StateInfo
	UpdatePending bool // Indicates that there are some update operations pending that should be completed by a servicing container.
}

// Stats contains statistics from HCS
type Stats hcsshim.Statistics

// Resources defines updatable container resource values.
type Resources struct{}

// ServicingOption is a CreateOption with a no-op application that signifies
// the container needs to be used for a Windows servicing operation.
type ServicingOption struct {
	IsServicing bool
}

// FlushOption is a CreateOption that signifies if the container should be
// started with flushes ignored until boot has completed. This is an optimisation
// for first boot of a container.
type FlushOption struct {
	IgnoreFlushesDuringBoot bool
}

// HyperVIsolationOption is a CreateOption that indicates whether the runtime
// should start the container as a Hyper-V container, and if so, the sandbox path.
type HyperVIsolationOption struct {
	IsHyperV    bool
	SandboxPath string `json:",omitempty"`
}

// LayerOption is a CreateOption that indicates to the runtime the layer folder
// and layer paths for a container.
type LayerOption struct {
	// LayerFolderPath is the path to the current layer folder. Empty for Hyper-V containers.
	LayerFolderPath string `json:",omitempty"`
	// Layer paths of the parent layers
	LayerPaths []string
}

// NetworkEndpointsOption is a CreateOption that provides the runtime list
// of network endpoints to which a container should be attached during its creation.
type NetworkEndpointsOption struct {
	Endpoints                []string
	AllowUnqualifiedDNSQuery bool
	DNSSearchList            []string
}

// CredentialsOption is a CreateOption that indicates the credentials from
// a credential spec to be used to the runtime
type CredentialsOption struct {
	Credentials string
}

// Checkpoint holds the details of a checkpoint (not supported in windows)
type Checkpoint struct {
	Name string
}

// Checkpoints contains the details of a checkpoint
type Checkpoints struct {
	Checkpoints []*Checkpoint
}
