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
// +build !windows

package dockerfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/system"
)

// normaliseDest normalises the destination of a COPY/ADD command in a
// platform semantically consistent way.
func normaliseDest(cmdName, workingDir, requested string) (string, error) {
	dest := filepath.FromSlash(requested)
	endsInSlash := strings.HasSuffix(requested, string(os.PathSeparator))
	if !system.IsAbs(requested) {
		dest = filepath.Join(string(os.PathSeparator), filepath.FromSlash(workingDir), dest)
		// Make sure we preserve any trailing slash
		if endsInSlash {
			dest += string(os.PathSeparator)
		}
	}
	return dest, nil
}

func containsWildcards(name string) bool {
	for i := 0; i < len(name); i++ {
		ch := name[i]
		if ch == '\\' {
			i++
		} else if ch == '*' || ch == '?' || ch == '[' {
			return true
		}
	}
	return false
}
