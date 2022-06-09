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

type BeFalseMatcher struct {
}

func (matcher *BeFalseMatcher) Match(actual interface{}) (success bool, err error) {
	if !isBool(actual) {
		return false, fmt.Errorf("Expected a boolean.  Got:\n%s", format.Object(actual, 1))
	}

	return actual == false, nil
}

func (matcher *BeFalseMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be false")
}

func (matcher *BeFalseMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be false")
}
