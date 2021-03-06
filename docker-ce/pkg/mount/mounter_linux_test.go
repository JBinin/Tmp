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

package mount

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMount(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("not root tests would fail")
	}

	source, err := ioutil.TempDir("", "mount-test-source-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(source)

	// Ensure we have a known start point by mounting tmpfs with given options
	if err := Mount("tmpfs", source, "tmpfs", "private"); err != nil {
		t.Fatal(err)
	}
	defer ensureUnmount(t, source)
	validateMount(t, source, "", "")
	if t.Failed() {
		t.FailNow()
	}

	target, err := ioutil.TempDir("", "mount-test-target-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(target)

	tests := []struct {
		source           string
		ftype            string
		options          string
		expectedOpts     string
		expectedOptional string
	}{
		// No options
		{"tmpfs", "tmpfs", "", "", ""},
		// Default rw / ro test
		{source, "", "bind", "", ""},
		{source, "", "bind,private", "", ""},
		{source, "", "bind,shared", "", "shared"},
		{source, "", "bind,slave", "", "master"},
		{source, "", "bind,unbindable", "", "unbindable"},
		// Read Write tests
		{source, "", "bind,rw", "rw", ""},
		{source, "", "bind,rw,private", "rw", ""},
		{source, "", "bind,rw,shared", "rw", "shared"},
		{source, "", "bind,rw,slave", "rw", "master"},
		{source, "", "bind,rw,unbindable", "rw", "unbindable"},
		// Read Only tests
		{source, "", "bind,ro", "ro", ""},
		{source, "", "bind,ro,private", "ro", ""},
		{source, "", "bind,ro,shared", "ro", "shared"},
		{source, "", "bind,ro,slave", "ro", "master"},
		{source, "", "bind,ro,unbindable", "ro", "unbindable"},
	}

	for _, tc := range tests {
		ftype, options := tc.ftype, tc.options
		if tc.ftype == "" {
			ftype = "none"
		}
		if tc.options == "" {
			options = "none"
		}

		t.Run(fmt.Sprintf("%v-%v", ftype, options), func(t *testing.T) {
			if strings.Contains(tc.options, "slave") {
				// Slave requires a shared source
				if err := MakeShared(source); err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := MakePrivate(source); err != nil {
						t.Fatal(err)
					}
				}()
			}
			if err := Mount(tc.source, target, tc.ftype, tc.options); err != nil {
				t.Fatal(err)
			}
			defer ensureUnmount(t, target)
			validateMount(t, target, tc.expectedOpts, tc.expectedOptional)
		})
	}
}

// ensureUnmount umounts mnt checking for errors
func ensureUnmount(t *testing.T, mnt string) {
	if err := Unmount(mnt); err != nil {
		t.Error(err)
	}
}

// validateMount checks that mnt has the given options
func validateMount(t *testing.T, mnt string, opts, optional string) {
	info, err := GetMounts()
	if err != nil {
		t.Fatal(err)
	}

	wantedOpts := make(map[string]struct{})
	if opts != "" {
		for _, opt := range strings.Split(opts, ",") {
			wantedOpts[opt] = struct{}{}
		}
	}

	wantedOptional := make(map[string]struct{})
	if optional != "" {
		for _, opt := range strings.Split(optional, ",") {
			wantedOptional[opt] = struct{}{}
		}
	}

	mnts := make(map[int]*Info, len(info))
	for _, mi := range info {
		mnts[mi.ID] = mi
	}

	for _, mi := range info {
		if mi.Mountpoint != mnt {
			continue
		}

		// Use parent info as the defaults
		p := mnts[mi.Parent]
		pOpts := make(map[string]struct{})
		if p.Opts != "" {
			for _, opt := range strings.Split(p.Opts, ",") {
				pOpts[clean(opt)] = struct{}{}
			}
		}
		pOptional := make(map[string]struct{})
		if p.Optional != "" {
			for _, field := range strings.Split(p.Optional, ",") {
				pOptional[clean(field)] = struct{}{}
			}
		}

		// Validate Opts
		if mi.Opts != "" {
			for _, opt := range strings.Split(mi.Opts, ",") {
				opt = clean(opt)
				if !has(wantedOpts, opt) && !has(pOpts, opt) {
					t.Errorf("unexpected mount option %q expected %q", opt, opts)
				}
				delete(wantedOpts, opt)
			}
		}
		for opt := range wantedOpts {
			t.Errorf("missing mount option %q found %q", opt, mi.Opts)
		}

		// Validate Optional
		if mi.Optional != "" {
			for _, field := range strings.Split(mi.Optional, ",") {
				field = clean(field)
				if !has(wantedOptional, field) && !has(pOptional, field) {
					t.Errorf("unexpected optional failed %q expected %q", field, optional)
				}
				delete(wantedOptional, field)
			}
		}
		for field := range wantedOptional {
			t.Errorf("missing optional field %q found %q", field, mi.Optional)
		}

		return
	}

	t.Errorf("failed to find mount %q", mnt)
}

// clean strips off any value param after the colon
func clean(v string) string {
	return strings.SplitN(v, ":", 2)[0]
}

// has returns true if key is a member of m
func has(m map[string]struct{}, key string) bool {
	_, ok := m[key]
	return ok
}
