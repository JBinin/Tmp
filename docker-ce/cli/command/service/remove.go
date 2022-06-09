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
package service

import (
	"fmt"
	"strings"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func newRemoveCommand(dockerCli *command.DockerCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "rm SERVICE [SERVICE...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more services",
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(dockerCli, args)
		},
	}
	cmd.Flags()

	return cmd
}

func runRemove(dockerCli *command.DockerCli, sids []string) error {
	client := dockerCli.Client()

	ctx := context.Background()

	var errs []string
	for _, sid := range sids {
		err := client.ServiceRemove(ctx, sid)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintf(dockerCli.Out(), "%s\n", sid)
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}
