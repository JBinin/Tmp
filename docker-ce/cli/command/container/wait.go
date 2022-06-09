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

type waitOptions struct {
	containers []string
}

// NewWaitCommand creates a new cobra.Command for `docker wait`
func NewWaitCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts waitOptions

	cmd := &cobra.Command{
		Use:   "wait CONTAINER [CONTAINER...]",
		Short: "Block until one or more containers stop, then print their exit codes",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runWait(dockerCli, &opts)
		},
	}
	return cmd
}

func runWait(dockerCli *command.DockerCli, opts *waitOptions) error {
	ctx := context.Background()

	var errs []string
	for _, container := range opts.containers {
		status, err := dockerCli.Client().ContainerWait(ctx, container)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintf(dockerCli.Out(), "%d\n", status)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
