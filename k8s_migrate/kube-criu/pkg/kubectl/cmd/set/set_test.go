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
Copyright 2016 The Kubernetes Authors.

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

package set

import (
	"testing"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientcmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func TestLocalAndDryRunFlags(t *testing.T) {
	f := clientcmdutil.NewFactory(genericclioptions.NewTestConfigFlags())
	setCmd := NewCmdSet(f, genericclioptions.NewTestIOStreamsDiscard())
	ensureLocalAndDryRunFlagsOnChildren(t, setCmd, "")
}

func ensureLocalAndDryRunFlagsOnChildren(t *testing.T, c *cobra.Command, prefix string) {
	for _, cmd := range c.Commands() {
		name := prefix + cmd.Name()
		if localFlag := cmd.Flag("local"); localFlag == nil {
			t.Errorf("Command %s does not implement the --local flag", name)
		}
		if dryRunFlag := cmd.Flag("dry-run"); dryRunFlag == nil {
			t.Errorf("Command %s does not implement the --dry-run flag", name)
		}
		ensureLocalAndDryRunFlagsOnChildren(t, cmd, name+".")
	}
}
