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
package system

import "runtime"

const defaultUnixPathEnv = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

// DefaultPathEnv is unix style list of directories to search for
// executables. Each directory is separated from the next by a colon
// ':' character .
func DefaultPathEnv(platform string) string {
	if runtime.GOOS == "windows" {
		if platform != runtime.GOOS && LCOWSupported() {
			return defaultUnixPathEnv
		}
		// Deliberately empty on Windows containers on Windows as the default path will be set by
		// the container. Docker has no context of what the default path should be.
		return ""
	}
	return defaultUnixPathEnv

}
