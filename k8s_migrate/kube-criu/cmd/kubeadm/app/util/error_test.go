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

package util

import (
	"fmt"
	"testing"
)

type pferror struct{}

func (p *pferror) Preflight() bool { return true }
func (p *pferror) Error() string   { return "" }
func TestCheckErr(t *testing.T) {
	var codeReturned int
	errHandle := func(err string, code int) {
		codeReturned = code
	}

	var tokenTest = []struct {
		e        error
		expected int
	}{
		{nil, 0},
		{fmt.Errorf(""), DefaultErrorExitCode},
		{&pferror{}, PreFlightExitCode},
	}

	for _, rt := range tokenTest {
		codeReturned = 0
		checkErr(rt.e, errHandle)
		if codeReturned != rt.expected {
			t.Errorf(
				"failed checkErr:\n\texpected: %d\n\t  actual: %d",
				rt.expected,
				codeReturned,
			)
		}
	}
}

func TestFormatErrMsg(t *testing.T) {
	errMsg1 := "specified version to upgrade to v1.9.0-alpha.3 is equal to or lower than the cluster version v1.10.0-alpha.0.69+638add6ddfb6d2. Downgrades are not supported yet"
	errMsg2 := "specified version to upgrade to v1.9.0-alpha.3 is higher than the kubeadm version v1.9.0-alpha.1.3121+84178212527295-dirty. Upgrade kubeadm first using the tool you used to install kubeadm"

	testCases := []struct {
		errs   []error
		expect string
	}{
		{
			errs: []error{
				fmt.Errorf(errMsg1),
				fmt.Errorf(errMsg2),
			},
			expect: "\t- " + errMsg1 + "\n" + "\t- " + errMsg2 + "\n",
		},
		{
			errs: []error{
				fmt.Errorf(errMsg1),
			},
			expect: "\t- " + errMsg1 + "\n",
		},
	}

	for _, testCase := range testCases {
		got := FormatErrMsg(testCase.errs)
		if got != testCase.expect {
			t.Errorf("FormatErrMsg error, expect: %v, got: %v", testCase.expect, got)
		}
	}
}
