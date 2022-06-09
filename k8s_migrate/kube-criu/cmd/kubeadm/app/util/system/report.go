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
/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package system

import (
	"fmt"
	"io"
	"os"
)

// ValidationResultType is type of the validation result. Different validation results
// corresponds to different colors.
type ValidationResultType int32

const (
	good ValidationResultType = iota
	bad
	warn
)

// color is the color of the message.
type color int32

const (
	red    color = 31
	green        = 32
	yellow       = 33
	white        = 37
)

func colorize(s string, c color) string {
	return fmt.Sprintf("\033[0;%dm%s\033[0m", c, s)
}

// The default reporter for the system verification test
type StreamReporter struct {
	// The stream that this reporter is writing to
	WriteStream io.Writer
}

func (dr *StreamReporter) Report(key, value string, resultType ValidationResultType) error {
	var c color
	switch resultType {
	case good:
		c = green
	case bad:
		c = red
	case warn:
		c = yellow
	default:
		c = white
	}
	if dr.WriteStream == nil {
		return fmt.Errorf("WriteStream has to be defined for this reporter")
	}

	fmt.Fprintf(dr.WriteStream, "%s: %s\n", colorize(key, white), colorize(value, c))
	return nil
}

// DefaultReporter is the default Reporter
var DefaultReporter = &StreamReporter{
	WriteStream: os.Stdout,
}
