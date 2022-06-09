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
package checkpoint

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type createOptions struct {
	container     string
	checkpoint    string
	checkpointDir string
	parentPath    string
	leaveRunning  bool
	preDump       bool
}

func newCreateCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts createOptions

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] CONTAINER CHECKPOINT",
		Short: "Create a checkpoint from a running container",
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.container = args[0]
			opts.checkpoint = args[1]
			return runCreate(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.leaveRunning, "leave-running", false, "Leave the container running after checkpoint")
	flags.BoolVar(&opts.preDump, "pre-dump", false, "Pre-dump is used to pre-copy live migration")
	flags.StringVarP(&opts.parentPath, "parent-path", "","", "Parent-Path is the path of last iteration dump image")
	flags.StringVarP(&opts.checkpointDir, "checkpoint-dir", "", "", "Use a custom checkpoint storage directory")

	return cmd
}

func runCreate(dockerCli *command.DockerCli, opts createOptions) error {
	client := dockerCli.Client()

	checkpointOpts := types.CheckpointCreateOptions{
		CheckpointID:  opts.checkpoint,
		CheckpointDir: opts.checkpointDir,
		PreDump:       opts.preDump,
		ParentPath:    opts.parentPath,
		Exit:          !opts.leaveRunning,
	}

	err := client.CheckpointCreate(context.Background(), opts.container, checkpointOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", opts.checkpoint)
	return nil
}
