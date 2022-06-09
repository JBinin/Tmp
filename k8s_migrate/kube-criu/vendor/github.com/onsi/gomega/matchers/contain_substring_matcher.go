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
	"strings"
)

type ContainSubstringMatcher struct {
	Substr string
	Args   []interface{}
}

func (matcher *ContainSubstringMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return false, fmt.Errorf("ContainSubstring matcher requires a string or stringer.  Got:\n%s", format.Object(actual, 1))
	}

	return strings.Contains(actualString, matcher.stringToMatch()), nil
}

func (matcher *ContainSubstringMatcher) stringToMatch() string {
	stringToMatch := matcher.Substr
	if len(matcher.Args) > 0 {
		stringToMatch = fmt.Sprintf(matcher.Substr, matcher.Args...)
	}
	return stringToMatch
}

func (matcher *ContainSubstringMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain substring", matcher.stringToMatch())
}

func (matcher *ContainSubstringMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to contain substring", matcher.stringToMatch())
}
