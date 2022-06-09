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
// +build !windows

package main

import (
	"strings"

	"github.com/go-check/check"
)

// Test case for #30166 (target was not validated)
func (s *DockerSuite) TestCreateTmpfsMountsTarget(c *check.C) {
	testRequires(c, DaemonIsLinux)
	type testCase struct {
		target        string
		expectedError string
	}
	cases := []testCase{
		{
			target:        ".",
			expectedError: "mount path must be absolute",
		},
		{
			target:        "foo",
			expectedError: "mount path must be absolute",
		},
		{
			target:        "/",
			expectedError: "destination can't be '/'",
		},
		{
			target:        "//",
			expectedError: "destination can't be '/'",
		},
	}
	for _, x := range cases {
		out, _, _ := dockerCmdWithError("create", "--tmpfs", x.target, "busybox", "sh")
		if x.expectedError != "" && !strings.Contains(out, x.expectedError) {
			c.Fatalf("mounting tmpfs over %q should fail with %q, but got %q",
				x.target, x.expectedError, out)
		}
	}
}
