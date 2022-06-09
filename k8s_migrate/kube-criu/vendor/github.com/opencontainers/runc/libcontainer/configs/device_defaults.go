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

package configs

var (
	// DefaultSimpleDevices are devices that are to be both allowed and created.
	DefaultSimpleDevices = []*Device{
		// /dev/null and zero
		{
			Path:        "/dev/null",
			Type:        'c',
			Major:       1,
			Minor:       3,
			Permissions: "rwm",
			FileMode:    0666,
		},
		{
			Path:        "/dev/zero",
			Type:        'c',
			Major:       1,
			Minor:       5,
			Permissions: "rwm",
			FileMode:    0666,
		},

		{
			Path:        "/dev/full",
			Type:        'c',
			Major:       1,
			Minor:       7,
			Permissions: "rwm",
			FileMode:    0666,
		},

		// consoles and ttys
		{
			Path:        "/dev/tty",
			Type:        'c',
			Major:       5,
			Minor:       0,
			Permissions: "rwm",
			FileMode:    0666,
		},

		// /dev/urandom,/dev/random
		{
			Path:        "/dev/urandom",
			Type:        'c',
			Major:       1,
			Minor:       9,
			Permissions: "rwm",
			FileMode:    0666,
		},
		{
			Path:        "/dev/random",
			Type:        'c',
			Major:       1,
			Minor:       8,
			Permissions: "rwm",
			FileMode:    0666,
		},
	}
	DefaultAllowedDevices = append([]*Device{
		// allow mknod for any device
		{
			Type:        'c',
			Major:       Wildcard,
			Minor:       Wildcard,
			Permissions: "m",
		},
		{
			Type:        'b',
			Major:       Wildcard,
			Minor:       Wildcard,
			Permissions: "m",
		},

		{
			Path:        "/dev/console",
			Type:        'c',
			Major:       5,
			Minor:       1,
			Permissions: "rwm",
		},
		// /dev/pts/ - pts namespaces are "coming soon"
		{
			Path:        "",
			Type:        'c',
			Major:       136,
			Minor:       Wildcard,
			Permissions: "rwm",
		},
		{
			Path:        "",
			Type:        'c',
			Major:       5,
			Minor:       2,
			Permissions: "rwm",
		},

		// tuntap
		{
			Path:        "",
			Type:        'c',
			Major:       10,
			Minor:       200,
			Permissions: "rwm",
		},
	}, DefaultSimpleDevices...)
	DefaultAutoCreatedDevices = append([]*Device{}, DefaultSimpleDevices...)
)
