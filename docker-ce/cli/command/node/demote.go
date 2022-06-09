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
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

func newDemoteCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "demote NODE [NODE...]",
		Short: "Demote one or more nodes from manager in the swarm",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDemote(dockerCli, args)
		},
	}
}

func runDemote(dockerCli command.Cli, nodes []string) error {
	demote := func(node *swarm.Node) error {
		if node.Spec.Role == swarm.NodeRoleWorker {
			fmt.Fprintf(dockerCli.Out(), "Node %s is already a worker.\n", node.ID)
			return errNoRoleChange
		}
		node.Spec.Role = swarm.NodeRoleWorker
		return nil
	}
	success := func(nodeID string) {
		fmt.Fprintf(dockerCli.Out(), "Manager %s demoted in the swarm.\n", nodeID)
	}
	return updateNodes(dockerCli, nodes, demote, success)
}
