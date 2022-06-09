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
package network

import (
	"github.com/spf13/cobra"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
)

// NewNetworkCommand returns a cobra command for `network` subcommands
func NewNetworkCommand(dockerCli *command.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage networks",
		Args:  cli.NoArgs,
		RunE:  dockerCli.ShowHelp,
	}
	cmd.AddCommand(
		newConnectCommand(dockerCli),
		newCreateCommand(dockerCli),
		newDisconnectCommand(dockerCli),
		newInspectCommand(dockerCli),
		newListCommand(dockerCli),
		newRemoveCommand(dockerCli),
		NewPruneCommand(dockerCli),
	)
	return cmd
}
