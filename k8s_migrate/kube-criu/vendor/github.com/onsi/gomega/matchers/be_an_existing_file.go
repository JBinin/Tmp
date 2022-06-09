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
package matchers

import (
	"fmt"
	"os"

	"github.com/onsi/gomega/format"
)

type BeAnExistingFileMatcher struct {
	expected interface{}
}

func (matcher *BeAnExistingFileMatcher) Match(actual interface{}) (success bool, err error) {
	actualFilename, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeAnExistingFileMatcher matcher expects a file path")
	}

	if _, err = os.Stat(actualFilename); err != nil {
		switch {
		case os.IsNotExist(err):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (matcher *BeAnExistingFileMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to exist"))
}

func (matcher *BeAnExistingFileMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to exist"))
}
