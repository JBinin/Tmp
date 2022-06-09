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
	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/command/inspect"
	"github.com/spf13/cobra"
)

type inspectOptions struct {
	format string
	names  []string
}

func newInspectCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] NETWORK [NETWORK...]",
		Short: "Display detailed information on one or more networks",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.names = args
			return runInspect(dockerCli, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runInspect(dockerCli *command.DockerCli, opts inspectOptions) error {
	client := dockerCli.Client()

	ctx := context.Background()

	getNetFunc := func(name string) (interface{}, []byte, error) {
		return client.NetworkInspectWithRaw(ctx, name)
	}

	return inspect.Inspect(dockerCli.Out(), opts.names, opts.format, getNetFunc)
}
