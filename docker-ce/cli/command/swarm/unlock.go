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
package swarm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"golang.org/x/net/context"
)

type unlockOptions struct{}

func newUnlockCommand(dockerCli command.Cli) *cobra.Command {
	opts := unlockOptions{}

	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Unlock swarm",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnlock(dockerCli, opts)
		},
	}

	return cmd
}

func runUnlock(dockerCli command.Cli, opts unlockOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	// First see if the node is actually part of a swarm, and if it is actually locked first.
	// If it's in any other state than locked, don't ask for the key.
	info, err := client.Info(ctx)
	if err != nil {
		return err
	}

	switch info.Swarm.LocalNodeState {
	case swarm.LocalNodeStateInactive:
		return errors.New("Error: This node is not part of a swarm")
	case swarm.LocalNodeStateLocked:
		break
	default:
		return errors.New("Error: swarm is not locked")
	}

	key, err := readKey(dockerCli.In(), "Please enter unlock key: ")
	if err != nil {
		return err
	}
	req := swarm.UnlockRequest{
		UnlockKey: key,
	}

	return client.SwarmUnlock(ctx, req)
}

func readKey(in *command.InStream, prompt string) (string, error) {
	if in.IsTerminal() {
		fmt.Print(prompt)
		dt, err := terminal.ReadPassword(int(in.FD()))
		fmt.Println()
		return string(dt), err
	}
	key, err := bufio.NewReader(in).ReadString('\n')
	if err == io.EOF {
		err = nil
	}
	return strings.TrimSpace(key), err
}
