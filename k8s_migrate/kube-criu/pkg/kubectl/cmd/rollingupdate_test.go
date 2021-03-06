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

package cmd

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdtesting "k8s.io/kubernetes/pkg/kubectl/cmd/testing"
)

func TestValidateArgs(t *testing.T) {
	f := cmdtesting.NewTestFactory()
	defer f.Cleanup()

	tests := []struct {
		testName  string
		flags     map[string]string
		filenames []string
		args      []string
		expectErr bool
	}{
		{
			testName:  "nothing",
			expectErr: true,
		},
		{
			testName:  "no file, no image",
			flags:     map[string]string{},
			args:      []string{"foo"},
			expectErr: true,
		},
		{
			testName:  "valid file example",
			filenames: []string{"bar.yaml"},
			args:      []string{"foo"},
		},
		{
			testName: "missing second image name",
			flags: map[string]string{
				"image": "foo:v2",
			},
			args: []string{"foo"},
		},
		{
			testName: "valid image example",
			flags: map[string]string{
				"image": "foo:v2",
			},
			args: []string{"foo", "foo-v2"},
		},
		{
			testName: "both filename and image example",
			flags: map[string]string{
				"image": "foo:v2",
			},
			filenames: []string{"bar.yaml"},
			args:      []string{"foo", "foo-v2"},
			expectErr: true,
		},
	}
	for _, test := range tests {
		cmd := NewCmdRollingUpdate(f, genericclioptions.NewTestIOStreamsDiscard())

		if test.flags != nil {
			for key, val := range test.flags {
				cmd.Flags().Set(key, val)
			}
		}
		err := validateArguments(cmd, test.filenames, test.args)
		if err != nil && !test.expectErr {
			t.Errorf("%s: unexpected error: %v", test.testName, err)
		}
		if err == nil && test.expectErr {
			t.Errorf("%s: unexpected non-error", test.testName)
		}
	}
}
