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
// Package local provides the default implementation for volumes. It
// is used to mount data volume containers and directories local to
// the host server.
package local

import (
	"fmt"
	"path/filepath"
	"strings"
)

type optsConfig struct{}

var validOpts map[string]bool

// scopedPath verifies that the path where the volume is located
// is under Docker's root and the valid local paths.
func (r *Root) scopedPath(realPath string) bool {
	if strings.HasPrefix(realPath, filepath.Join(r.scope, volumesPathName)) && realPath != filepath.Join(r.scope, volumesPathName) {
		return true
	}
	return false
}

func setOpts(v *localVolume, opts map[string]string) error {
	if len(opts) > 0 {
		return fmt.Errorf("options are not supported on this platform")
	}
	return nil
}

func (v *localVolume) mount() error {
	return nil
}
