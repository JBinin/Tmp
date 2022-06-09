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

type HaveSuffixMatcher struct {
	Suffix string
	Args   []interface{}
}

func (matcher *HaveSuffixMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return false, fmt.Errorf("HaveSuffix matcher requires a string or stringer.  Got:\n%s", format.Object(actual, 1))
	}
	suffix := matcher.suffix()
	return len(actualString) >= len(suffix) && actualString[len(actualString)-len(suffix):] == suffix, nil
}

func (matcher *HaveSuffixMatcher) suffix() string {
	if len(matcher.Args) > 0 {
		return fmt.Sprintf(matcher.Suffix, matcher.Args...)
	}
	return matcher.Suffix
}

func (matcher *HaveSuffixMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to have suffix", matcher.suffix())
}

func (matcher *HaveSuffixMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to have suffix", matcher.suffix())
}
