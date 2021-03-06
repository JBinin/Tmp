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
// +build linux

/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package empty_dir

import (
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/sys/unix"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/util/mount"
)

// Defined by Linux - the type number for tmpfs mounts.
const (
	linuxTmpfsMagic     = 0x01021994
	linuxHugetlbfsMagic = 0x958458f6
)

// realMountDetector implements mountDetector in terms of syscalls.
type realMountDetector struct {
	mounter mount.Interface
}

func (m *realMountDetector) GetMountMedium(path string) (v1.StorageMedium, bool, error) {
	glog.V(5).Infof("Determining mount medium of %v", path)
	notMnt, err := m.mounter.IsLikelyNotMountPoint(path)
	if err != nil {
		return v1.StorageMediumDefault, false, fmt.Errorf("IsLikelyNotMountPoint(%q): %v", path, err)
	}
	buf := unix.Statfs_t{}
	if err := unix.Statfs(path, &buf); err != nil {
		return v1.StorageMediumDefault, false, fmt.Errorf("statfs(%q): %v", path, err)
	}

	glog.V(5).Infof("Statfs_t of %v: %+v", path, buf)
	if buf.Type == linuxTmpfsMagic {
		return v1.StorageMediumMemory, !notMnt, nil
	} else if int64(buf.Type) == linuxHugetlbfsMagic {
		return v1.StorageMediumHugePages, !notMnt, nil
	}
	return v1.StorageMediumDefault, !notMnt, nil
}
