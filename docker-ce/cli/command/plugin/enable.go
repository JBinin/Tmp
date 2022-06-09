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

type enableOpts struct {
	timeout int
	name    string
}

func newEnableCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts enableOpts

	cmd := &cobra.Command{
		Use:   "enable [OPTIONS] PLUGIN",
		Short: "Enable a plugin",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runEnable(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runEnable(dockerCli *command.DockerCli, opts *enableOpts) error {
	name := opts.name
	if opts.timeout < 0 {
		return fmt.Errorf("negative timeout %d is invalid", opts.timeout)
	}

	if err := dockerCli.Client().PluginEnable(context.Background(), name, types.PluginEnableOptions{Timeout: opts.timeout}); err != nil {
		return err
	}
	fmt.Fprintln(dockerCli.Out(), name)
	return nil
}
