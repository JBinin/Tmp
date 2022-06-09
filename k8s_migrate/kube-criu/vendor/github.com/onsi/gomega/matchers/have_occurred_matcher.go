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

	"github.com/onsi/gomega/format"
)

type HaveOccurredMatcher struct {
}

func (matcher *HaveOccurredMatcher) Match(actual interface{}) (success bool, err error) {
	// is purely nil?
	if actual == nil {
		return false, nil
	}

	// must be an 'error' type
	if !isError(actual) {
		return false, fmt.Errorf("Expected an error-type.  Got:\n%s", format.Object(actual, 1))
	}

	// must be non-nil (or a pointer to a non-nil)
	return !isNil(actual), nil
}

func (matcher *HaveOccurredMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected an error to have occurred.  Got:\n%s", format.Object(actual, 1))
}

func (matcher *HaveOccurredMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected error:\n%s\n%s\n%s", format.Object(actual, 1), format.IndentString(actual.(error).Error(), 1), "not to have occurred")
}
