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
	"io"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type exportOptions struct {
	container string
	output    string
}

// NewExportCommand creates a new `docker export` command
func NewExportCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts exportOptions

	cmd := &cobra.Command{
		Use:   "export [OPTIONS] CONTAINER",
		Short: "Export a container's filesystem as a tar archive",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.container = args[0]
			return runExport(dockerCli, opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.output, "output", "o", "", "Write to a file, instead of STDOUT")

	return cmd
}

func runExport(dockerCli *command.DockerCli, opts exportOptions) error {
	if opts.output == "" && dockerCli.Out().IsTerminal() {
		return errors.New("Cowardly refusing to save to a terminal. Use the -o flag or redirect.")
	}

	clnt := dockerCli.Client()

	responseBody, err := clnt.ContainerExport(context.Background(), opts.container)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	if opts.output == "" {
		_, err := io.Copy(dockerCli.Out(), responseBody)
		return err
	}

	return command.CopyToFile(opts.output, responseBody)
}
