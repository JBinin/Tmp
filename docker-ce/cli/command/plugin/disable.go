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
package plugin

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func newDisableCommand(dockerCli *command.DockerCli) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "disable [OPTIONS] PLUGIN",
		Short: "Disable a plugin",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisable(dockerCli, args[0], force)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Force the disable of an active plugin")
	return cmd
}

func runDisable(dockerCli *command.DockerCli, name string, force bool) error {
	if err := dockerCli.Client().PluginDisable(context.Background(), name, types.PluginDisableOptions{Force: force}); err != nil {
		return err
	}
	fmt.Fprintln(dockerCli.Out(), name)
	return nil
}
