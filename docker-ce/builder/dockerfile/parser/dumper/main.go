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

	"github.com/docker/docker/builder/dockerfile/parser"
)

func main() {
	var f *os.File
	var err error

	if len(os.Args) < 2 {
		fmt.Println("please supply filename(s)")
		os.Exit(1)
	}

	for _, fn := range os.Args[1:] {
		f, err = os.Open(fn)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		d := parser.Directive{LookingForDirectives: true}
		parser.SetEscapeToken(parser.DefaultEscapeToken, &d)

		ast, err := parser.Parse(f, &d)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(ast.Dump())
		}
	}
}
