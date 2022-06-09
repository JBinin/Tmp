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
// +build !daemon

package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func newDaemonCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "daemon",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDaemon()
		},
	}
}

func runDaemon() error {
	return fmt.Errorf(
		"`docker daemon` is not supported on %s. Please run `dockerd` directly",
		strings.Title(runtime.GOOS))
}
