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
package image

import (
	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type tagOptions struct {
	image string
	name  string
}

// NewTagCommand creates a new `docker tag` command
func NewTagCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts tagOptions

	cmd := &cobra.Command{
		Use:   "tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG]",
		Short: "Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE",
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			opts.name = args[1]
			return runTag(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	return cmd
}

func runTag(dockerCli *command.DockerCli, opts tagOptions) error {
	ctx := context.Background()

	return dockerCli.Client().ImageTag(ctx, opts.image, opts.name)
}
