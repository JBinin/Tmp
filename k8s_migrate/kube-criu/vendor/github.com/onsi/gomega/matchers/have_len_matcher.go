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

type HaveLenMatcher struct {
	Count int
}

func (matcher *HaveLenMatcher) Match(actual interface{}) (success bool, err error) {
	length, ok := lengthOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}

	return length == matcher.Count, nil
}

func (matcher *HaveLenMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto have length %d", format.Object(actual, 1), matcher.Count)
}

func (matcher *HaveLenMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to have length %d", format.Object(actual, 1), matcher.Count)
}
