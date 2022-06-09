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
	"fmt"
	"os"
	"path/filepath"
)

// normaliseWorkdir normalises a user requested working directory in a
// platform sematically consistent way.
func normaliseWorkdir(current string, requested string) (string, error) {
	if requested == "" {
		return "", fmt.Errorf("cannot normalise nothing")
	}
	current = filepath.FromSlash(current)
	requested = filepath.FromSlash(requested)
	if !filepath.IsAbs(requested) {
		return filepath.Join(string(os.PathSeparator), current, requested), nil
	}
	return requested, nil
}

func errNotJSON(command, _ string) error {
	return fmt.Errorf("%s requires the arguments to be in JSON form", command)
}
