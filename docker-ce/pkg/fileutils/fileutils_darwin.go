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
package fileutils

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GetTotalUsedFds returns the number of used File Descriptors by
// executing `lsof -p PID`
func GetTotalUsedFds() int {
	pid := os.Getpid()

	cmd := exec.Command("lsof", "-p", strconv.Itoa(pid))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1
	}

	outputStr := strings.TrimSpace(string(output))

	fds := strings.Split(outputStr, "\n")

	return len(fds) - 1
}
