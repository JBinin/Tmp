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
	"reflect"

	"github.com/onsi/gomega/format"
)

type EqualMatcher struct {
	Expected interface{}
}

func (matcher *EqualMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	return reflect.DeepEqual(actual, matcher.Expected), nil
}

func (matcher *EqualMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *EqualMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}
