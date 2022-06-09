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
package dockerfile

import "strings"

// handleJSONArgs parses command passed to CMD, ENTRYPOINT, RUN and SHELL instruction in Dockerfile
// for exec form it returns untouched args slice
// for shell form it returns concatenated args as the first element of a slice
func handleJSONArgs(args []string, attributes map[string]bool) []string {
	if len(args) == 0 {
		return []string{}
	}

	if attributes != nil && attributes["json"] {
		return args
	}

	// literal string command, not an exec array
	return []string{strings.Join(args, " ")}
}
