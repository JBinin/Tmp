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
package volume

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/mount"
)

func TestValidateMount(t *testing.T) {
	testDir, err := ioutil.TempDir("", "test-validate-mount")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	cases := []struct {
		input    mount.Mount
		expected error
	}{
		{mount.Mount{Type: mount.TypeVolume}, errMissingField("Target")},
		{mount.Mount{Type: mount.TypeVolume, Target: testDestinationPath, Source: "hello"}, nil},
		{mount.Mount{Type: mount.TypeVolume, Target: testDestinationPath}, nil},
		{mount.Mount{Type: mount.TypeBind}, errMissingField("Target")},
		{mount.Mount{Type: mount.TypeBind, Target: testDestinationPath}, errMissingField("Source")},
		{mount.Mount{Type: mount.TypeBind, Target: testDestinationPath, Source: testSourcePath, VolumeOptions: &mount.VolumeOptions{}}, errExtraField("VolumeOptions")},
		{mount.Mount{Type: mount.TypeBind, Source: testSourcePath, Target: testDestinationPath}, errBindNotExist},
		{mount.Mount{Type: mount.TypeBind, Source: testDir, Target: testDestinationPath}, nil},
		{mount.Mount{Type: "invalid", Target: testDestinationPath}, errors.New("mount type unknown")},
	}
	for i, x := range cases {
		err := validateMountConfig(&x.input)
		if err == nil && x.expected == nil {
			continue
		}
		if (err == nil && x.expected != nil) || (x.expected == nil && err != nil) || !strings.Contains(err.Error(), x.expected.Error()) {
			t.Fatalf("expected %q, got %q, case: %d", x.expected, err, i)
		}
	}
}
