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
	"flag"
	"fmt"
	"os/exec"
)

func BuildUnfocusCommand() *Command {
	return &Command{
		Name:         "unfocus",
		AltName:      "blur",
		FlagSet:      flag.NewFlagSet("unfocus", flag.ExitOnError),
		UsageCommand: "ginkgo unfocus (or ginkgo blur)",
		Usage: []string{
			"Recursively unfocuses any focused tests under the current directory",
		},
		Command: unfocusSpecs,
	}
}

func unfocusSpecs([]string, []string) {
	unfocus("Describe")
	unfocus("Context")
	unfocus("It")
	unfocus("Measure")
	unfocus("DescribeTable")
	unfocus("Entry")
}

func unfocus(component string) {
	fmt.Printf("Removing F%s...\n", component)
	cmd := exec.Command("gofmt", fmt.Sprintf("-r=F%s -> %s", component, component), "-w", ".")
	out, _ := cmd.CombinedOutput()
	if string(out) != "" {
		println(string(out))
	}
}
