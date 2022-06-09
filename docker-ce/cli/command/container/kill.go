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
package container

import (
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type killOptions struct {
	signal string

	containers []string
}

// NewKillCommand creates a new cobra.Command for `docker kill`
func NewKillCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts killOptions

	cmd := &cobra.Command{
		Use:   "kill [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Kill one or more running containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runKill(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.signal, "signal", "s", "KILL", "Signal to send to the container")
	return cmd
}

func runKill(dockerCli *command.DockerCli, opts *killOptions) error {
	var errs []string
	ctx := context.Background()
	errChan := parallelOperation(ctx, opts.containers, func(ctx context.Context, container string) error {
		return dockerCli.Client().ContainerKill(ctx, container, opts.signal)
	})
	for _, name := range opts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
		} else {
			fmt.Fprintln(dockerCli.Out(), name)
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
