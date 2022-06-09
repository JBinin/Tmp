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
// +build !windows

package idtools

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func resolveBinary(binname string) (string, error) {
	binaryPath, err := exec.LookPath(binname)
	if err != nil {
		return "", err
	}
	resolvedPath, err := filepath.EvalSymlinks(binaryPath)
	if err != nil {
		return "", err
	}
	//only return no error if the final resolved binary basename
	//matches what was searched for
	if filepath.Base(resolvedPath) == binname {
		return resolvedPath, nil
	}
	return "", fmt.Errorf("Binary %q does not resolve to a binary of that name in $PATH (%q)", binname, resolvedPath)
}

func execCmd(cmd, args string) ([]byte, error) {
	execCmd := exec.Command(cmd, strings.Split(args, " ")...)
	return execCmd.CombinedOutput()
}
