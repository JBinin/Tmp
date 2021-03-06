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
	"io"

	containerd "github.com/docker/containerd/api/grpc/types"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/net/context"
)

// State constants used in state change reporting.
const (
	StateStart       = "start-container"
	StatePause       = "pause"
	StateResume      = "resume"
	StateExit        = "exit"
	StateRestore     = "restore"
	StateExitProcess = "exit-process"
	StateOOM         = "oom" // fake state
)

// CommonStateInfo contains the state info common to all platforms.
type CommonStateInfo struct { // FIXME: event?
	State     string
	Pid       uint32
	ExitCode  uint32
	ProcessID string
}

// Backend defines callbacks that the client of the library needs to implement.
type Backend interface {
	StateChanged(containerID string, state StateInfo) error
}

// Client provides access to containerd features.
type Client interface {
	GetServerVersion(ctx context.Context) (*ServerVersion, error)
	Create(containerID string, checkpoint string, checkpointDir string, spec specs.Spec, attachStdio StdioCallback, options ...CreateOption) error
	Signal(containerID string, sig int) error
	SignalProcess(containerID string, processFriendlyName string, sig int) error
	AddProcess(ctx context.Context, containerID, processFriendlyName string, process Process, attachStdio StdioCallback) (int, error)
	Resize(containerID, processFriendlyName string, width, height int) error
	Pause(containerID string) error
	Resume(containerID string) error
	Restore(containerID string, attachStdio StdioCallback, options ...CreateOption) error
	Stats(containerID string) (*Stats, error)
	GetPidsForContainer(containerID string) ([]int, error)
	Summary(containerID string) ([]Summary, error)
	UpdateResources(containerID string, resources Resources) error
	CreateCheckpoint(containerID string, checkpointID string, checkpointDir string,preDump bool,parentPath string, exit bool) error
	DeleteCheckpoint(containerID string, checkpointID string, checkpointDir string) error
	ListCheckpoints(containerID string, checkpointDir string) (*Checkpoints, error)
}

// CreateOption allows to configure parameters of container creation.
type CreateOption interface {
	Apply(interface{}) error
}

// StdioCallback is called to connect a container or process stdio.
type StdioCallback func(IOPipe) error

// IOPipe contains the stdio streams.
type IOPipe struct {
	Stdin    io.WriteCloser
	Stdout   io.ReadCloser
	Stderr   io.ReadCloser
	Terminal bool // Whether stderr is connected on Windows
}

// ServerVersion contains version information as retrieved from the
// server
type ServerVersion struct {
	containerd.GetServerVersionResponse
}
