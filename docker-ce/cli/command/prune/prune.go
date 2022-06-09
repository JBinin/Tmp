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
package prune

import (
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/command/container"
	"github.com/docker/docker/cli/command/image"
	"github.com/docker/docker/cli/command/network"
	"github.com/docker/docker/cli/command/volume"
	"github.com/docker/docker/opts"
	"github.com/spf13/cobra"
)

// NewContainerPruneCommand returns a cobra prune command for containers
func NewContainerPruneCommand(dockerCli *command.DockerCli) *cobra.Command {
	return container.NewPruneCommand(dockerCli)
}

// NewVolumePruneCommand returns a cobra prune command for volumes
func NewVolumePruneCommand(dockerCli *command.DockerCli) *cobra.Command {
	return volume.NewPruneCommand(dockerCli)
}

// NewImagePruneCommand returns a cobra prune command for images
func NewImagePruneCommand(dockerCli *command.DockerCli) *cobra.Command {
	return image.NewPruneCommand(dockerCli)
}

// NewNetworkPruneCommand returns a cobra prune command for Networks
func NewNetworkPruneCommand(dockerCli *command.DockerCli) *cobra.Command {
	return network.NewPruneCommand(dockerCli)
}

// RunContainerPrune executes a prune command for containers
func RunContainerPrune(dockerCli *command.DockerCli, filter opts.FilterOpt) (uint64, string, error) {
	return container.RunPrune(dockerCli, filter)
}

// RunVolumePrune executes a prune command for volumes
func RunVolumePrune(dockerCli *command.DockerCli, filter opts.FilterOpt) (uint64, string, error) {
	return volume.RunPrune(dockerCli)
}

// RunImagePrune executes a prune command for images
func RunImagePrune(dockerCli *command.DockerCli, all bool, filter opts.FilterOpt) (uint64, string, error) {
	return image.RunPrune(dockerCli, all, filter)
}

// RunNetworkPrune executes a prune command for networks
func RunNetworkPrune(dockerCli *command.DockerCli, filter opts.FilterOpt) (uint64, string, error) {
	return network.RunPrune(dockerCli, filter)
}
