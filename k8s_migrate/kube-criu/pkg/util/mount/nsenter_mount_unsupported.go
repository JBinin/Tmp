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

/*
Copyright 2014 The Kubernetes Authors.

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

package mount

import (
	"errors"
	"os"

	"k8s.io/kubernetes/pkg/util/nsenter"
)

type NsenterMounter struct{}

func NewNsenterMounter(rootDir string, ne *nsenter.Nsenter) *NsenterMounter {
	return &NsenterMounter{}
}

var _ = Interface(&NsenterMounter{})

func (*NsenterMounter) Mount(source string, target string, fstype string, options []string) error {
	return nil
}

func (*NsenterMounter) Unmount(target string) error {
	return nil
}

func (*NsenterMounter) List() ([]MountPoint, error) {
	return []MountPoint{}, nil
}

func (m *NsenterMounter) IsNotMountPoint(dir string) (bool, error) {
	return IsNotMountPoint(m, dir)
}

func (*NsenterMounter) IsMountPointMatch(mp MountPoint, dir string) bool {
	return (mp.Path == dir)
}

func (*NsenterMounter) IsLikelyNotMountPoint(file string) (bool, error) {
	return true, nil
}

func (*NsenterMounter) DeviceOpened(pathname string) (bool, error) {
	return false, nil
}

func (*NsenterMounter) PathIsDevice(pathname string) (bool, error) {
	return true, nil
}

func (*NsenterMounter) GetDeviceNameFromMount(mountPath, pluginDir string) (string, error) {
	return "", nil
}

func (*NsenterMounter) MakeRShared(path string) error {
	return nil
}

func (*NsenterMounter) GetFileType(_ string) (FileType, error) {
	return FileType("fake"), errors.New("not implemented")
}

func (*NsenterMounter) MakeDir(pathname string) error {
	return nil
}

func (*NsenterMounter) MakeFile(pathname string) error {
	return nil
}

func (*NsenterMounter) ExistsPath(pathname string) (bool, error) {
	return true, errors.New("not implemented")
}

func (*NsenterMounter) EvalHostSymlinks(pathname string) (string, error) {
	return "", errors.New("not implemented")
}

func (*NsenterMounter) SafeMakeDir(pathname string, base string, perm os.FileMode) error {
	return nil
}

func (*NsenterMounter) PrepareSafeSubpath(subPath Subpath) (newHostPath string, cleanupAction func(), err error) {
	return subPath.Path, nil, nil
}

func (*NsenterMounter) CleanSubPaths(podDir string, volumeName string) error {
	return nil
}

func (*NsenterMounter) GetMountRefs(pathname string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (*NsenterMounter) GetFSGroup(pathname string) (int64, error) {
	return -1, errors.New("not implemented")
}

func (*NsenterMounter) GetSELinuxSupport(pathname string) (bool, error) {
	return false, errors.New("not implemented")
}

func (*NsenterMounter) GetMode(pathname string) (os.FileMode, error) {
	return 0, errors.New("not implemented")
}
