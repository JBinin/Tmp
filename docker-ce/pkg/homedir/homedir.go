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
package homedir

import (
	"os"
	"runtime"

	"github.com/opencontainers/runc/libcontainer/user"
)

// Key returns the env var name for the user's home dir based on
// the platform being run on
func Key() string {
	if runtime.GOOS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}

// Get returns the home directory of the current user with the help of
// environment variables depending on the target operating system.
// Returned path should be used with "path/filepath" to form new paths.
func Get() string {
	home := os.Getenv(Key())
	if home == "" && runtime.GOOS != "windows" {
		if u, err := user.CurrentUser(); err == nil {
			return u.Home
		}
	}
	return home
}

// GetShortcutString returns the string that is shortcut to user's home directory
// in the native shell of the platform running on.
func GetShortcutString() string {
	if runtime.GOOS == "windows" {
		return "%USERPROFILE%" // be careful while using in format functions
	}
	return "~"
}
