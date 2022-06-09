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
// +build windows

package system

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultPathEnv is deliberately empty on Windows as the default path will be set by
// the container. Docker has no context of what the default path should be.
const DefaultPathEnv = ""

// CheckSystemDriveAndRemoveDriveLetter verifies and manipulates a Windows path.
// This is used, for example, when validating a user provided path in docker cp.
// If a drive letter is supplied, it must be the system drive. The drive letter
// is always removed. Also, it translates it to OS semantics (IOW / to \). We
// need the path in this syntax so that it can ultimately be contatenated with
// a Windows long-path which doesn't support drive-letters. Examples:
// C:			--> Fail
// C:\			--> \
// a			--> a
// /a			--> \a
// d:\			--> Fail
func CheckSystemDriveAndRemoveDriveLetter(path string) (string, error) {
	if len(path) == 2 && string(path[1]) == ":" {
		return "", fmt.Errorf("No relative path specified in %q", path)
	}
	if !filepath.IsAbs(path) || len(path) < 2 {
		return filepath.FromSlash(path), nil
	}
	if string(path[1]) == ":" && !strings.EqualFold(string(path[0]), "c") {
		return "", fmt.Errorf("The specified path is not on the system drive (C:)")
	}
	return filepath.FromSlash(path[2:]), nil
}
