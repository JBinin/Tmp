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
package secret

import (
	"fmt"
	"strings"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type removeOptions struct {
	names []string
}

func newSecretRemoveCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:     "rm SECRET [SECRET...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more secrets",
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := removeOptions{
				names: args,
			}
			return runSecretRemove(dockerCli, opts)
		},
	}
}

func runSecretRemove(dockerCli *command.DockerCli, opts removeOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	var errs []string

	for _, name := range opts.names {
		if err := client.SecretRemove(ctx, name); err != nil {
			errs = append(errs, err.Error())
			continue
		}

		fmt.Fprintln(dockerCli.Out(), name)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	return nil
}
