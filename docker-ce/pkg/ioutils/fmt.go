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
package ioutils

import (
	"fmt"
	"io"
)

// FprintfIfNotEmpty prints the string value if it's not empty
func FprintfIfNotEmpty(w io.Writer, format, value string) (int, error) {
	if value != "" {
		return fmt.Fprintf(w, format, value)
	}
	return 0, nil
}

// FprintfIfTrue prints the boolean value if it's true
func FprintfIfTrue(w io.Writer, format string, ok bool) (int, error) {
	if ok {
		return fmt.Fprintf(w, format, ok)
	}
	return 0, nil
}
