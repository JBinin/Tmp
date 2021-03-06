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
package convert

import (
	"testing"

	"github.com/docker/docker/api/types/mount"
	composetypes "github.com/docker/docker/cli/compose/types"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestIsReadOnly(t *testing.T) {
	assert.Equal(t, isReadOnly([]string{"foo", "bar", "ro"}), true)
	assert.Equal(t, isReadOnly([]string{"ro"}), true)
	assert.Equal(t, isReadOnly([]string{}), false)
	assert.Equal(t, isReadOnly([]string{"foo", "rw"}), false)
	assert.Equal(t, isReadOnly([]string{"foo"}), false)
}

func TestIsNoCopy(t *testing.T) {
	assert.Equal(t, isNoCopy([]string{"foo", "bar", "nocopy"}), true)
	assert.Equal(t, isNoCopy([]string{"nocopy"}), true)
	assert.Equal(t, isNoCopy([]string{}), false)
	assert.Equal(t, isNoCopy([]string{"foo", "rw"}), false)
}

func TestGetBindOptions(t *testing.T) {
	opts := getBindOptions([]string{"slave"})
	expected := mount.BindOptions{Propagation: mount.PropagationSlave}
	assert.Equal(t, *opts, expected)
}

func TestGetBindOptionsNone(t *testing.T) {
	opts := getBindOptions([]string{"ro"})
	assert.Equal(t, opts, (*mount.BindOptions)(nil))
}

func TestConvertVolumeToMountAnonymousVolume(t *testing.T) {
	stackVolumes := volumes{}
	namespace := NewNamespace("foo")
	expected := mount.Mount{
		Type:   mount.TypeVolume,
		Target: "/foo/bar",
	}
	mount, err := convertVolumeToMount("/foo/bar", stackVolumes, namespace)
	assert.NilError(t, err)
	assert.DeepEqual(t, mount, expected)
}

func TestConvertVolumeToMountInvalidFormat(t *testing.T) {
	namespace := NewNamespace("foo")
	invalids := []string{"::", "::cc", ":bb:", "aa::", "aa::cc", "aa:bb:", " : : ", " : :cc", " :bb: ", "aa: : ", "aa: :cc", "aa:bb: "}
	for _, vol := range invalids {
		_, err := convertVolumeToMount(vol, volumes{}, namespace)
		assert.Error(t, err, "invalid volume: "+vol)
	}
}

func TestConvertVolumeToMountNamedVolume(t *testing.T) {
	stackVolumes := volumes{
		"normal": composetypes.VolumeConfig{
			Driver: "glusterfs",
			DriverOpts: map[string]string{
				"opt": "value",
			},
			Labels: map[string]string{
				"something": "labeled",
			},
		},
	}
	namespace := NewNamespace("foo")
	expected := mount.Mount{
		Type:     mount.TypeVolume,
		Source:   "foo_normal",
		Target:   "/foo",
		ReadOnly: true,
		VolumeOptions: &mount.VolumeOptions{
			Labels: map[string]string{
				LabelNamespace: "foo",
				"something":    "labeled",
			},
			DriverConfig: &mount.Driver{
				Name: "glusterfs",
				Options: map[string]string{
					"opt": "value",
				},
			},
		},
	}
	mount, err := convertVolumeToMount("normal:/foo:ro", stackVolumes, namespace)
	assert.NilError(t, err)
	assert.DeepEqual(t, mount, expected)
}

func TestConvertVolumeToMountNamedVolumeExternal(t *testing.T) {
	stackVolumes := volumes{
		"outside": composetypes.VolumeConfig{
			External: composetypes.External{
				External: true,
				Name:     "special",
			},
		},
	}
	namespace := NewNamespace("foo")
	expected := mount.Mount{
		Type:   mount.TypeVolume,
		Source: "special",
		Target: "/foo",
		VolumeOptions: &mount.VolumeOptions{
			NoCopy: false,
		},
	}
	mount, err := convertVolumeToMount("outside:/foo", stackVolumes, namespace)
	assert.NilError(t, err)
	assert.DeepEqual(t, mount, expected)
}

func TestConvertVolumeToMountNamedVolumeExternalNoCopy(t *testing.T) {
	stackVolumes := volumes{
		"outside": composetypes.VolumeConfig{
			External: composetypes.External{
				External: true,
				Name:     "special",
			},
		},
	}
	namespace := NewNamespace("foo")
	expected := mount.Mount{
		Type:   mount.TypeVolume,
		Source: "special",
		Target: "/foo",
		VolumeOptions: &mount.VolumeOptions{
			NoCopy: true,
		},
	}
	mount, err := convertVolumeToMount("outside:/foo:nocopy", stackVolumes, namespace)
	assert.NilError(t, err)
	assert.DeepEqual(t, mount, expected)
}

func TestConvertVolumeToMountBind(t *testing.T) {
	stackVolumes := volumes{}
	namespace := NewNamespace("foo")
	expected := mount.Mount{
		Type:        mount.TypeBind,
		Source:      "/bar",
		Target:      "/foo",
		ReadOnly:    true,
		BindOptions: &mount.BindOptions{Propagation: mount.PropagationShared},
	}
	mount, err := convertVolumeToMount("/bar:/foo:ro,shared", stackVolumes, namespace)
	assert.NilError(t, err)
	assert.DeepEqual(t, mount, expected)
}

func TestConvertVolumeToMountVolumeDoesNotExist(t *testing.T) {
	namespace := NewNamespace("foo")
	_, err := convertVolumeToMount("unknown:/foo:ro", volumes{}, namespace)
	assert.Error(t, err, "undefined volume: unknown")
}
