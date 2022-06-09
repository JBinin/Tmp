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
package image

import (
	"github.com/spf13/cobra"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
)

// NewImageCommand returns a cobra command for `image` subcommands
func NewImageCommand(dockerCli *command.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage images",
		Args:  cli.NoArgs,
		RunE:  dockerCli.ShowHelp,
	}
	cmd.AddCommand(
		NewBuildCommand(dockerCli),
		NewHistoryCommand(dockerCli),
		NewImportCommand(dockerCli),
		NewLoadCommand(dockerCli),
		NewPullCommand(dockerCli),
		NewPushCommand(dockerCli),
		NewSaveCommand(dockerCli),
		NewTagCommand(dockerCli),
		newListCommand(dockerCli),
		newRemoveCommand(dockerCli),
		newInspectCommand(dockerCli),
		NewPruneCommand(dockerCli),
	)
	return cmd
}
