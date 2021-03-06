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

package cm

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/kubernetes/pkg/util/mount"
)

func fakeContainerMgrMountInt() mount.Interface {
	return &mount.FakeMounter{
		MountPoints: []mount.MountPoint{
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuset"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpu"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuacct"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "memory"},
			},
		},
	}
}

func TestCgroupMountValidationSuccess(t *testing.T) {
	f, err := validateSystemRequirements(fakeContainerMgrMountInt())
	assert.Nil(t, err)
	assert.False(t, f.cpuHardcapping, "cpu hardcapping is expected to be disabled")
}

func TestCgroupMountValidationMemoryMissing(t *testing.T) {
	mountInt := &mount.FakeMounter{
		MountPoints: []mount.MountPoint{
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuset"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpu"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuacct"},
			},
		},
	}
	_, err := validateSystemRequirements(mountInt)
	assert.Error(t, err)
}

func TestCgroupMountValidationMultipleSubsystem(t *testing.T) {
	mountInt := &mount.FakeMounter{
		MountPoints: []mount.MountPoint{
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuset", "memory"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpu"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuacct"},
			},
		},
	}
	_, err := validateSystemRequirements(mountInt)
	assert.Nil(t, err)
}

func TestSoftRequirementsValidationSuccess(t *testing.T) {
	req := require.New(t)
	tempDir, err := ioutil.TempDir("", "")
	req.NoError(err)
	defer os.RemoveAll(tempDir)
	req.NoError(ioutil.WriteFile(path.Join(tempDir, "cpu.cfs_period_us"), []byte("0"), os.ModePerm))
	req.NoError(ioutil.WriteFile(path.Join(tempDir, "cpu.cfs_quota_us"), []byte("0"), os.ModePerm))
	mountInt := &mount.FakeMounter{
		MountPoints: []mount.MountPoint{
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuset"},
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpu"},
				Path:   tempDir,
			},
			{
				Device: "cgroup",
				Type:   "cgroup",
				Opts:   []string{"rw", "relatime", "cpuacct", "memory"},
			},
		},
	}
	f, err := validateSystemRequirements(mountInt)
	assert.NoError(t, err)
	assert.True(t, f.cpuHardcapping, "cpu hardcapping is expected to be enabled")
}
