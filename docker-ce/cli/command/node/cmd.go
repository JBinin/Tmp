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
package node

import (
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	apiclient "github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewNodeCommand returns a cobra command for `node` subcommands
func NewNodeCommand(dockerCli *command.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage Swarm nodes",
		Args:  cli.NoArgs,
		RunE:  dockerCli.ShowHelp,
	}
	cmd.AddCommand(
		newDemoteCommand(dockerCli),
		newInspectCommand(dockerCli),
		newListCommand(dockerCli),
		newPromoteCommand(dockerCli),
		newRemoveCommand(dockerCli),
		newPsCommand(dockerCli),
		newUpdateCommand(dockerCli),
	)
	return cmd
}

// Reference returns the reference of a node. The special value "self" for a node
// reference is mapped to the current node, hence the node ID is retrieved using
// the `/info` endpoint.
func Reference(ctx context.Context, client apiclient.APIClient, ref string) (string, error) {
	if ref == "self" {
		info, err := client.Info(ctx)
		if err != nil {
			return "", err
		}
		return info.Swarm.NodeID, nil
	}
	return ref, nil
}
