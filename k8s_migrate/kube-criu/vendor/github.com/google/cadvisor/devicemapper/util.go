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
// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package devicemapper

import (
	"fmt"
	"os"
	"path/filepath"
)

// ThinLsBinaryPresent returns the location of the thin_ls binary in the mount
// namespace cadvisor is running in or an error.  The locations checked are:
//
// - /sbin/
// - /bin/
// - /usr/sbin/
// - /usr/bin/
//
// The thin_ls binary is provided by the device-mapper-persistent-data
// package.
func ThinLsBinaryPresent() (string, error) {
	var (
		thinLsPath string
		err        error
	)

	for _, path := range []string{"/sbin", "/bin", "/usr/sbin/", "/usr/bin"} {
		// try paths for non-containerized operation
		// note: thin_ls is most likely a symlink to pdata_tools
		thinLsPath = filepath.Join(path, "thin_ls")
		_, err = os.Stat(thinLsPath)
		if err == nil {
			return thinLsPath, nil
		}
	}

	return "", fmt.Errorf("unable to find thin_ls binary")
}
