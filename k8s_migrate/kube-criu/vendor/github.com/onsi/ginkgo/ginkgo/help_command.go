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
)

func BuildHelpCommand() *Command {
	return &Command{
		Name:         "help",
		FlagSet:      flag.NewFlagSet("help", flag.ExitOnError),
		UsageCommand: "ginkgo help <COMMAND>",
		Usage: []string{
			"Print usage information.  If a command is passed in, print usage information just for that command.",
		},
		Command: printHelp,
	}
}

func printHelp(args []string, additionalArgs []string) {
	if len(args) == 0 {
		usage()
	} else {
		command, found := commandMatching(args[0])
		if !found {
			complainAndQuit(fmt.Sprintf("Unknown command: %s", args[0]))
		}

		usageForCommand(command, true)
	}
}
