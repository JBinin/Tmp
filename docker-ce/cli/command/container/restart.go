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
	"time"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type restartOptions struct {
	nSeconds        int
	nSecondsChanged bool

	containers []string
}

// NewRestartCommand creates a new cobra.Command for `docker restart`
func NewRestartCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts restartOptions

	cmd := &cobra.Command{
		Use:   "restart [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Restart one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			opts.nSecondsChanged = cmd.Flags().Changed("time")
			return runRestart(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.nSeconds, "time", "t", 10, "Seconds to wait for stop before killing the container")
	return cmd
}

func runRestart(dockerCli *command.DockerCli, opts *restartOptions) error {
	ctx := context.Background()
	var errs []string
	var timeout *time.Duration
	if opts.nSecondsChanged {
		timeoutValue := time.Duration(opts.nSeconds) * time.Second
		timeout = &timeoutValue
	}

	for _, name := range opts.containers {
		if err := dockerCli.Client().ContainerRestart(ctx, name, timeout); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintln(dockerCli.Out(), name)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
