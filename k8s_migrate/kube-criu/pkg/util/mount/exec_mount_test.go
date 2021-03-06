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
Copyright 2017 The Kubernetes Authors.

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
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var (
	sourcePath      = "/mnt/srv"
	destinationPath = "/mnt/dst"
	fsType          = "xfs"
	mountOptions    = []string{"vers=1", "foo=bar"}
)

func TestMount(t *testing.T) {
	exec := NewFakeExec(func(cmd string, args ...string) ([]byte, error) {
		if cmd != "mount" {
			t.Errorf("expected mount command, got %q", cmd)
		}
		// mount -t fstype -o options source target
		expectedArgs := []string{"-t", fsType, "-o", strings.Join(mountOptions, ","), sourcePath, destinationPath}
		if !reflect.DeepEqual(expectedArgs, args) {
			t.Errorf("expected arguments %q, got %q", strings.Join(expectedArgs, " "), strings.Join(args, " "))
		}
		return nil, nil
	})

	wrappedMounter := &fakeMounter{FakeMounter: &FakeMounter{}, t: t}
	mounter := NewExecMounter(exec, wrappedMounter)

	mounter.Mount(sourcePath, destinationPath, fsType, mountOptions)
}

func TestBindMount(t *testing.T) {
	cmdCount := 0
	exec := NewFakeExec(func(cmd string, args ...string) ([]byte, error) {
		cmdCount++
		if cmd != "mount" {
			t.Errorf("expected mount command, got %q", cmd)
		}
		var expectedArgs []string
		switch cmdCount {
		case 1:
			// mount -t fstype -o "bind" source target
			expectedArgs = []string{"-t", fsType, "-o", "bind", sourcePath, destinationPath}
		case 2:
			// mount -t fstype -o "remount,opts" source target
			expectedArgs = []string{"-t", fsType, "-o", "bind,remount," + strings.Join(mountOptions, ","), sourcePath, destinationPath}
		}
		if !reflect.DeepEqual(expectedArgs, args) {
			t.Errorf("expected arguments %q, got %q", strings.Join(expectedArgs, " "), strings.Join(args, " "))
		}
		return nil, nil
	})

	wrappedMounter := &fakeMounter{FakeMounter: &FakeMounter{}, t: t}
	mounter := NewExecMounter(exec, wrappedMounter)
	bindOptions := append(mountOptions, "bind")
	mounter.Mount(sourcePath, destinationPath, fsType, bindOptions)
}

func TestUnmount(t *testing.T) {
	exec := NewFakeExec(func(cmd string, args ...string) ([]byte, error) {
		if cmd != "umount" {
			t.Errorf("expected unmount command, got %q", cmd)
		}
		// unmount $target
		expectedArgs := []string{destinationPath}
		if !reflect.DeepEqual(expectedArgs, args) {
			t.Errorf("expected arguments %q, got %q", strings.Join(expectedArgs, " "), strings.Join(args, " "))
		}
		return nil, nil
	})

	wrappedMounter := &fakeMounter{&FakeMounter{}, t}
	mounter := NewExecMounter(exec, wrappedMounter)

	mounter.Unmount(destinationPath)
}

/* Fake wrapped mounter */
type fakeMounter struct {
	*FakeMounter
	t *testing.T
}

func (fm *fakeMounter) Mount(source string, target string, fstype string, options []string) error {
	// Mount() of wrapped mounter should never be called. We call exec instead.
	fm.t.Errorf("Unexpected wrapped mount call")
	return fmt.Errorf("Unexpected wrapped mount call")
}

func (fm *fakeMounter) Unmount(target string) error {
	// umount() of wrapped mounter should never be called. We call exec instead.
	fm.t.Errorf("Unexpected wrapped mount call")
	return fmt.Errorf("Unexpected wrapped mount call")
}
