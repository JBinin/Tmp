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
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

func newRemoveCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:     "rm NETWORK [NETWORK...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more networks",
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(dockerCli, args)
		},
	}
}

func runRemove(dockerCli *command.DockerCli, networks []string) error {
	client := dockerCli.Client()
	ctx := context.Background()
	status := 0

	for _, name := range networks {
		if err := client.NetworkRemove(ctx, name); err != nil {
			fmt.Fprintf(dockerCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		fmt.Fprintf(dockerCli.Out(), "%s\n", name)
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}
