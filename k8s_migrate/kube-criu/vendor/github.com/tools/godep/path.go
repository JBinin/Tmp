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
	"fmt"
	"os"
)

var cmdPath = &Command{
	Name:  "path",
	Short: "print GOPATH for dependency code",
	Long: `
Command path prints a path for use in env var GOPATH
that makes available the specified version of each dependency.

The printed path does not include any GOPATH value from
the environment.

For more about how GOPATH works, see 'go help gopath'.
`,
	Run:          runPath,
	OnlyInGOPATH: true,
}

// Print the gopath that points to
// the included dependency code.
func runPath(cmd *Command, args []string) {
	if len(args) != 0 {
		cmd.UsageExit()
	}
	if VendorExperiment {
		fmt.Fprintln(os.Stderr, "Error: GO15VENDOREXPERIMENT is enabled and the vendor/ directory is not a valid Go workspace.")
		os.Exit(1)
	}
	gopath := prepareGopath()
	fmt.Println(gopath)
}
