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
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/ginkgo/interrupthandler"
	"github.com/onsi/ginkgo/ginkgo/testrunner"
)

func BuildBuildCommand() *Command {
	commandFlags := NewBuildCommandFlags(flag.NewFlagSet("build", flag.ExitOnError))
	interruptHandler := interrupthandler.NewInterruptHandler()
	builder := &SpecBuilder{
		commandFlags:     commandFlags,
		interruptHandler: interruptHandler,
	}

	return &Command{
		Name:         "build",
		FlagSet:      commandFlags.FlagSet,
		UsageCommand: "ginkgo build <FLAGS> <PACKAGES>",
		Usage: []string{
			"Build the passed in <PACKAGES> (or the package in the current directory if left blank).",
			"Accepts the following flags:",
		},
		Command: builder.BuildSpecs,
	}
}

type SpecBuilder struct {
	commandFlags     *RunWatchAndBuildCommandFlags
	interruptHandler *interrupthandler.InterruptHandler
}

func (r *SpecBuilder) BuildSpecs(args []string, additionalArgs []string) {
	r.commandFlags.computeNodes()

	suites, _ := findSuites(args, r.commandFlags.Recurse, r.commandFlags.SkipPackage, false)

	if len(suites) == 0 {
		complainAndQuit("Found no test suites")
	}

	passed := true
	for _, suite := range suites {
		runner := testrunner.New(suite, 1, false, r.commandFlags.GoOpts, nil)
		fmt.Printf("Compiling %s...\n", suite.PackageName)

		path, _ := filepath.Abs(filepath.Join(suite.Path, fmt.Sprintf("%s.test", suite.PackageName)))
		err := runner.CompileTo(path)
		if err != nil {
			fmt.Println(err.Error())
			passed = false
		} else {
			fmt.Printf("    compiled %s.test\n", suite.PackageName)
		}

		runner.CleanUp()
	}

	if passed {
		os.Exit(0)
	}
	os.Exit(1)
}
