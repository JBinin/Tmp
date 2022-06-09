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
package dbus

import (
	"os"
	"sync"
)

var (
	homeDir     string
	homeDirLock sync.Mutex
)

func getHomeDir() string {
	homeDirLock.Lock()
	defer homeDirLock.Unlock()

	if homeDir != "" {
		return homeDir
	}

	homeDir = os.Getenv("HOME")
	if homeDir != "" {
		return homeDir
	}

	homeDir = lookupHomeDir()
	return homeDir
}
