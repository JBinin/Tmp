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
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type topOptions struct {
	container string

	args []string
}

// NewTopCommand creates a new cobra.Command for `docker top`
func NewTopCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts topOptions

	cmd := &cobra.Command{
		Use:   "top CONTAINER [ps OPTIONS]",
		Short: "Display the running processes of a container",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.container = args[0]
			opts.args = args[1:]
			return runTop(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	return cmd
}

func runTop(dockerCli *command.DockerCli, opts *topOptions) error {
	ctx := context.Background()

	procList, err := dockerCli.Client().ContainerTop(ctx, opts.container, opts.args)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(dockerCli.Out(), 20, 1, 3, ' ', 0)
	fmt.Fprintln(w, strings.Join(procList.Titles, "\t"))

	for _, proc := range procList.Processes {
		fmt.Fprintln(w, strings.Join(proc, "\t"))
	}
	w.Flush()
	return nil
}
