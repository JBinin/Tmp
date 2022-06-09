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

package system

// DefaultPathEnv is unix style list of directories to search for
// executables. Each directory is separated from the next by a colon
// ':' character .
const DefaultPathEnv = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

// CheckSystemDriveAndRemoveDriveLetter verifies that a path, if it includes a drive letter,
// is the system drive. This is a no-op on Linux.
func CheckSystemDriveAndRemoveDriveLetter(path string) (string, error) {
	return path, nil
}
