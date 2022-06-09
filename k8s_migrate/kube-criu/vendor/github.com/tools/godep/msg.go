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
	"log"

	"github.com/kr/pretty"
)

func debugln(a ...interface{}) (int, error) {
	if debug {
		return fmt.Println(a...)
	}
	return 0, nil
}

func verboseln(a ...interface{}) {
	if verbose {
		log.Println(a...)
	}
}

func debugf(format string, a ...interface{}) (int, error) {
	if debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func ppln(a ...interface{}) (int, error) {
	if debug {
		return pretty.Println(a...)
	}
	return 0, nil
}
