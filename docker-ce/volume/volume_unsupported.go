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
// +build !linux

package volume

import (
	"fmt"
	"runtime"

	mounttypes "github.com/docker/docker/api/types/mount"
)

// ConvertTmpfsOptions converts *mounttypes.TmpfsOptions to the raw option string
// for mount(2).
func ConvertTmpfsOptions(opt *mounttypes.TmpfsOptions, readOnly bool) (string, error) {
	return "", fmt.Errorf("%s does not support tmpfs", runtime.GOOS)
}
