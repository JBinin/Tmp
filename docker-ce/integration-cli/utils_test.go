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
package main

import (
	"os/exec"

	"github.com/docker/docker/pkg/testutil/cmd"
)

func getPrefixAndSlashFromDaemonPlatform() (prefix, slash string) {
	if testEnv.DaemonPlatform() == "windows" {
		return "c:", `\`
	}
	return "", "/"
}

// TODO: update code to call cmd.RunCmd directly, and remove this function
// Deprecated: use pkg/testutil/cmd instead
func runCommandWithOutput(execCmd *exec.Cmd) (string, int, error) {
	result := cmd.RunCmd(transformCmd(execCmd))
	return result.Combined(), result.ExitCode, result.Error
}

// Temporary shim for migrating commands to the new function
func transformCmd(execCmd *exec.Cmd) cmd.Cmd {
	return cmd.Cmd{
		Command: execCmd.Args,
		Env:     execCmd.Env,
		Dir:     execCmd.Dir,
		Stdin:   execCmd.Stdin,
		Stdout:  execCmd.Stdout,
	}
}
