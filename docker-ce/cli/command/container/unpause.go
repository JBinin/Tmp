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

type unpauseOptions struct {
	containers []string
}

// NewUnpauseCommand creates a new cobra.Command for `docker unpause`
func NewUnpauseCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts unpauseOptions

	cmd := &cobra.Command{
		Use:   "unpause CONTAINER [CONTAINER...]",
		Short: "Unpause all processes within one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runUnpause(dockerCli, &opts)
		},
	}
	return cmd
}

func runUnpause(dockerCli *command.DockerCli, opts *unpauseOptions) error {
	ctx := context.Background()

	var errs []string
	errChan := parallelOperation(ctx, opts.containers, dockerCli.Client().ContainerUnpause)
	for _, container := range opts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintln(dockerCli.Out(), container)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
