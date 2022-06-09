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
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type leaveOptions struct {
	force bool
}

func newLeaveCommand(dockerCli command.Cli) *cobra.Command {
	opts := leaveOptions{}

	cmd := &cobra.Command{
		Use:   "leave [OPTIONS]",
		Short: "Leave the swarm",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLeave(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.force, "force", "f", false, "Force this node to leave the swarm, ignoring warnings")
	return cmd
}

func runLeave(dockerCli command.Cli, opts leaveOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	if err := client.SwarmLeave(ctx, opts.force); err != nil {
		return err
	}

	fmt.Fprintln(dockerCli.Out(), "Node left the swarm.")
	return nil
}
