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
package stenographer

import (
	"fmt"
	"strings"
)

func (s *consoleStenographer) colorize(colorCode string, format string, args ...interface{}) string {
	var out string

	if len(args) > 0 {
		out = fmt.Sprintf(format, args...)
	} else {
		out = format
	}

	if s.color {
		return fmt.Sprintf("%s%s%s", colorCode, out, defaultStyle)
	} else {
		return out
	}
}

func (s *consoleStenographer) printBanner(text string, bannerCharacter string) {
	fmt.Fprintln(s.w, text)
	fmt.Fprintln(s.w, strings.Repeat(bannerCharacter, len(text)))
}

func (s *consoleStenographer) printNewLine() {
	fmt.Fprintln(s.w, "")
}

func (s *consoleStenographer) printDelimiter() {
	fmt.Fprintln(s.w, s.colorize(grayColor, "%s", strings.Repeat("-", 30)))
}

func (s *consoleStenographer) print(indentation int, format string, args ...interface{}) {
	fmt.Fprint(s.w, s.indent(indentation, format, args...))
}

func (s *consoleStenographer) println(indentation int, format string, args ...interface{}) {
	fmt.Fprintln(s.w, s.indent(indentation, format, args...))
}

func (s *consoleStenographer) indent(indentation int, format string, args ...interface{}) string {
	var text string

	if len(args) > 0 {
		text = fmt.Sprintf(format, args...)
	} else {
		text = format
	}

	stringArray := strings.Split(text, "\n")
	padding := ""
	if indentation >= 0 {
		padding = strings.Repeat("  ", indentation)
	}
	for i, s := range stringArray {
		stringArray[i] = fmt.Sprintf("%s%s", padding, s)
	}

	return strings.Join(stringArray, "\n")
}
