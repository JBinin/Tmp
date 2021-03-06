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
//+build !windows

package daemon

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func validatePSArgs(psArgs string) error {
	// NOTE: \\s does not detect unicode whitespaces.
	// So we use fieldsASCII instead of strings.Fields in parsePSOutput.
	// See https://github.com/docker/docker/pull/24358
	re := regexp.MustCompile("\\s+([^\\s]*)=\\s*(PID[^\\s]*)")
	for _, group := range re.FindAllStringSubmatch(psArgs, -1) {
		if len(group) >= 3 {
			k := group[1]
			v := group[2]
			if k != "pid" {
				return fmt.Errorf("specifying \"%s=%s\" is not allowed", k, v)
			}
		}
	}
	return nil
}

// fieldsASCII is similar to strings.Fields but only allows ASCII whitespaces
func fieldsASCII(s string) []string {
	fn := func(r rune) bool {
		switch r {
		case '\t', '\n', '\f', '\r', ' ':
			return true
		}
		return false
	}
	return strings.FieldsFunc(s, fn)
}

func parsePSOutput(output []byte, pids []int) (*container.ContainerTopOKBody, error) {
	procList := &container.ContainerTopOKBody{}

	lines := strings.Split(string(output), "\n")
	procList.Titles = fieldsASCII(lines[0])

	pidIndex := -1
	for i, name := range procList.Titles {
		if name == "PID" {
			pidIndex = i
		}
	}
	if pidIndex == -1 {
		return nil, fmt.Errorf("Couldn't find PID field in ps output")
	}

	// loop through the output and extract the PID from each line
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		fields := fieldsASCII(line)
		p, err := strconv.Atoi(fields[pidIndex])
		if err != nil {
			return nil, fmt.Errorf("Unexpected pid '%s': %s", fields[pidIndex], err)
		}

		for _, pid := range pids {
			if pid == p {
				// Make sure number of fields equals number of header titles
				// merging "overhanging" fields
				process := fields[:len(procList.Titles)-1]
				process = append(process, strings.Join(fields[len(procList.Titles)-1:], " "))
				procList.Processes = append(procList.Processes, process)
			}
		}
	}
	return procList, nil
}

// ContainerTop lists the processes running inside of the given
// container by calling ps with the given args, or with the flags
// "-ef" if no args are given.  An error is returned if the container
// is not found, or is not running, or if there are any problems
// running ps, or parsing the output.
func (daemon *Daemon) ContainerTop(name string, psArgs string) (*container.ContainerTopOKBody, error) {
	if psArgs == "" {
		psArgs = "-ef"
	}

	if err := validatePSArgs(psArgs); err != nil {
		return nil, err
	}

	container, err := daemon.GetContainer(name)
	if err != nil {
		return nil, err
	}

	if !container.IsRunning() {
		return nil, errNotRunning{container.ID}
	}

	if container.IsRestarting() {
		return nil, errContainerIsRestarting(container.ID)
	}

	pids, err := daemon.containerd.GetPidsForContainer(container.ID)
	if err != nil {
		return nil, err
	}

	output, err := exec.Command("ps", strings.Split(psArgs, " ")...).Output()
	if err != nil {
		return nil, fmt.Errorf("Error running ps: %v", err)
	}
	procList, err := parsePSOutput(output, pids)
	if err != nil {
		return nil, err
	}
	daemon.LogContainerEvent(container, "top")
	return procList, nil
}
