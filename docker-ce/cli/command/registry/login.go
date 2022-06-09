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
package registry

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	serverAddress string
	user          string
	password      string
	email         string
}

// NewLoginCommand creates a new `docker login` command
func NewLoginCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts loginOptions

	cmd := &cobra.Command{
		Use:   "login [OPTIONS] [SERVER]",
		Short: "Log in to a Docker registry",
		Long:  "Log in to a Docker registry.\nIf no server is specified, the default is defined by the daemon.",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.serverAddress = args[0]
			}
			return runLogin(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.user, "username", "u", "", "Username")
	flags.StringVarP(&opts.password, "password", "p", "", "Password")

	return cmd
}

func runLogin(dockerCli *command.DockerCli, opts loginOptions) error {
	ctx := context.Background()
	clnt := dockerCli.Client()

	var (
		serverAddress string
		authServer    = command.ElectAuthServer(ctx, dockerCli)
	)
	if opts.serverAddress != "" {
		serverAddress = opts.serverAddress
	} else {
		serverAddress = authServer
	}

	isDefaultRegistry := serverAddress == authServer

	authConfig, err := command.ConfigureAuth(dockerCli, opts.user, opts.password, serverAddress, isDefaultRegistry)
	if err != nil {
		return err
	}
	response, err := clnt.RegistryLogin(ctx, authConfig)
	if err != nil {
		return err
	}
	if response.IdentityToken != "" {
		authConfig.Password = ""
		authConfig.IdentityToken = response.IdentityToken
	}
	if err := dockerCli.CredentialsStore(serverAddress).Store(authConfig); err != nil {
		return fmt.Errorf("Error saving credentials: %v", err)
	}

	if response.Status != "" {
		fmt.Fprintln(dockerCli.Out(), response.Status)
	}
	return nil
}
