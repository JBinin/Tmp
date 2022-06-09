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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type rmOptions struct {
	rmVolumes bool
	rmLink    bool
	force     bool

	containers []string
}

// NewRmCommand creates a new cobra.Command for `docker rm`
func NewRmCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts rmOptions

	cmd := &cobra.Command{
		Use:   "rm [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Remove one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runRm(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.rmVolumes, "volumes", "v", false, "Remove the volumes associated with the container")
	flags.BoolVarP(&opts.rmLink, "link", "l", false, "Remove the specified link")
	flags.BoolVarP(&opts.force, "force", "f", false, "Force the removal of a running container (uses SIGKILL)")
	return cmd
}

func runRm(dockerCli *command.DockerCli, opts *rmOptions) error {
	ctx := context.Background()

	var errs []string
	options := types.ContainerRemoveOptions{
		RemoveVolumes: opts.rmVolumes,
		RemoveLinks:   opts.rmLink,
		Force:         opts.force,
	}

	errChan := parallelOperation(ctx, opts.containers, func(ctx context.Context, container string) error {
		container = strings.Trim(container, "/")
		if container == "" {
			return errors.New("Container name cannot be empty")
		}
		return dockerCli.Client().ContainerRemove(ctx, container, options)
	})

	for _, name := range opts.containers {
		if err := <-errChan; err != nil {
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
