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
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Extensions to the standard "os" package.
package osext

import "path/filepath"

// Executable returns an absolute path that can be used to
// re-invoke the current program.
// It may not be valid after the current program exits.
func Executable() (string, error) {
	p, err := executable()
	return filepath.Clean(p), err
}

// Returns same path as Executable, returns just the folder
// path. Excludes the executable name.
func ExecutableFolder() (string, error) {
	p, err := Executable()
	if err != nil {
		return "", err
	}
	folder, _ := filepath.Split(p)
	return folder, nil
}
