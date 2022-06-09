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
// +build daemon

package main

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/spf13/cobra"
)

func stubRun(cmd *cobra.Command, args []string) error {
	return nil
}

func TestDaemonCommandHelp(t *testing.T) {
	cmd := newDaemonCommand()
	cmd.RunE = stubRun
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NilError(t, err)
}

func TestDaemonCommand(t *testing.T) {
	cmd := newDaemonCommand()
	cmd.RunE = stubRun
	cmd.SetArgs([]string{"--containerd", "/foo"})
	err := cmd.Execute()
	assert.NilError(t, err)
}
