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

package phases

import (
	"fmt"
	"os"
	"testing"

	testutil "k8s.io/kubernetes/cmd/kubeadm/test"
	cmdtestutil "k8s.io/kubernetes/cmd/kubeadm/test/cmd"
)

func TestEtcdSubCommandsHasFlags(t *testing.T) {

	subCmds := getEtcdSubCommands("", phaseTestK8sVersion)

	commonFlags := []string{
		"cert-dir",
		"config",
	}

	var tests = []struct {
		command         string
		additionalFlags []string
	}{
		{
			command: "local",
		},
	}

	for _, test := range tests {
		expectedFlags := append(commonFlags, test.additionalFlags...)
		cmdtestutil.AssertSubCommandHasFlags(t, subCmds, test.command, expectedFlags...)
	}
}

func TestEtcdCreateFilesWithFlags(t *testing.T) {

	var tests = []struct {
		command         string
		additionalFlags []string
		expectedFiles   []string
	}{
		{
			command:         "local",
			expectedFiles:   []string{"etcd.yaml"},
			additionalFlags: []string{},
		},
	}

	for _, test := range tests {

		// Create temp folder for the test case
		tmpdir := testutil.SetupTempDir(t)
		defer os.RemoveAll(tmpdir)

		// Get subcommands working in the temporary directory
		subCmds := getEtcdSubCommands(tmpdir, phaseTestK8sVersion)

		// Execute the subcommand
		certDirFlag := fmt.Sprintf("--cert-dir=%s", tmpdir)
		allFlags := append(test.additionalFlags, certDirFlag)
		cmdtestutil.RunSubCommand(t, subCmds, test.command, allFlags...)

		// Checks that requested files are there
		testutil.AssertFileExists(t, tmpdir, test.expectedFiles...)
	}
}
