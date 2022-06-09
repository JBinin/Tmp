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
	"github.com/onsi/ginkgo/ginkgo/convert"
	"os"
)

func BuildConvertCommand() *Command {
	return &Command{
		Name:         "convert",
		FlagSet:      flag.NewFlagSet("convert", flag.ExitOnError),
		UsageCommand: "ginkgo convert /path/to/package",
		Usage: []string{
			"Convert the package at the passed in path from an XUnit-style test to a Ginkgo-style test",
		},
		Command: convertPackage,
	}
}

func convertPackage(args []string, additionalArgs []string) {
	if len(args) != 1 {
		println(fmt.Sprintf("usage: ginkgo convert /path/to/your/package"))
		os.Exit(1)
	}

	defer func() {
		err := recover()
		if err != nil {
			switch err := err.(type) {
			case error:
				println(err.Error())
			case string:
				println(err)
			default:
				println(fmt.Sprintf("unexpected error: %#v", err))
			}
			os.Exit(1)
		}
	}()

	convert.RewritePackage(args[0])
}
