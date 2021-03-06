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
/*
Copyright 2018 The Kubernetes Authors.

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

package util

import (
	"testing"
)

func TestGetCgroupDriverDocker(t *testing.T) {
	testCases := []struct {
		name          string
		info          string
		expectedError bool
	}{
		{
			name:          "valid: value is 'cgroupfs'",
			info:          `Cgroup Driver: cgroupfs`,
			expectedError: false,
		},
		{
			name:          "valid: value is 'systemd'",
			info:          `Cgroup Driver: systemd`,
			expectedError: false,
		},
		{
			name:          "invalid: missing 'Cgroup Driver' key and value",
			info:          "",
			expectedError: true,
		},
		{
			name:          "invalid: only a 'Cgroup Driver' key is present",
			info:          `Cgroup Driver`,
			expectedError: true,
		},
		{
			name:          "invalid: empty 'Cgroup Driver' value",
			info:          `Cgroup Driver: `,
			expectedError: true,
		},
		{
			name:          "invalid: unknown 'Cgroup Driver' value",
			info:          `Cgroup Driver: invalid-value`,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := getCgroupDriverFromDockerInfo(tc.info); (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, saw: %v, error: %v", tc.expectedError, (err != nil), err)
			}
		})
	}
}
