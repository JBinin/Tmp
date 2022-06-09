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
	"bufio"
	"flag"
	"github.com/onsi/ginkgo/ginkgo/nodot"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func BuildNodotCommand() *Command {
	return &Command{
		Name:         "nodot",
		FlagSet:      flag.NewFlagSet("bootstrap", flag.ExitOnError),
		UsageCommand: "ginkgo nodot",
		Usage: []string{
			"Update the nodot declarations in your test suite",
			"Any missing declarations (from, say, a recently added matcher) will be added to your bootstrap file.",
			"If you've renamed a declaration, that name will be honored and not overwritten.",
		},
		Command: updateNodot,
	}
}

func updateNodot(args []string, additionalArgs []string) {
	suiteFile, perm := findSuiteFile()

	data, err := ioutil.ReadFile(suiteFile)
	if err != nil {
		complainAndQuit("Failed to update nodot declarations: " + err.Error())
	}

	content, err := nodot.ApplyNoDot(data)
	if err != nil {
		complainAndQuit("Failed to update nodot declarations: " + err.Error())
	}
	ioutil.WriteFile(suiteFile, content, perm)

	goFmt(suiteFile)
}

func findSuiteFile() (string, os.FileMode) {
	workingDir, err := os.Getwd()
	if err != nil {
		complainAndQuit("Could not find suite file for nodot: " + err.Error())
	}

	files, err := ioutil.ReadDir(workingDir)
	if err != nil {
		complainAndQuit("Could not find suite file for nodot: " + err.Error())
	}

	re := regexp.MustCompile(`RunSpecs\(|RunSpecsWithDefaultAndCustomReporters\(|RunSpecsWithCustomReporters\(`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(workingDir, file.Name())
		f, err := os.Open(path)
		if err != nil {
			complainAndQuit("Could not find suite file for nodot: " + err.Error())
		}
		defer f.Close()

		if re.MatchReader(bufio.NewReader(f)) {
			return path, file.Mode()
		}
	}

	complainAndQuit("Could not find a suite file for nodot: you need a bootstrap file that call's Ginkgo's RunSpecs() command.\nTry running ginkgo bootstrap first.")

	return "", 0
}
