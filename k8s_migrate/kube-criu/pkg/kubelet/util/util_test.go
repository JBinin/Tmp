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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		endpoint         string
		expectError      bool
		expectedProtocol string
		expectedAddr     string
	}{
		{
			endpoint:         "unix:///tmp/s1.sock",
			expectedProtocol: "unix",
			expectedAddr:     "/tmp/s1.sock",
		},
		{
			endpoint:         "tcp://localhost:15880",
			expectedProtocol: "tcp",
			expectedAddr:     "localhost:15880",
		},
		{
			endpoint:         "tcp1://abc",
			expectedProtocol: "tcp1",
			expectError:      true,
		},
		{
			endpoint:    "a b c",
			expectError: true,
		},
	}

	for _, test := range tests {
		protocol, addr, err := parseEndpoint(test.endpoint)
		assert.Equal(t, test.expectedProtocol, protocol)
		if test.expectError {
			assert.NotNil(t, err, "Expect error during parsing %q", test.endpoint)
			continue
		}
		assert.Nil(t, err, "Expect no error during parsing %q", test.endpoint)
		assert.Equal(t, test.expectedAddr, addr)
	}

}
