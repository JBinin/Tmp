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
	"reflect"
)

type HaveKeyMatcher struct {
	Key interface{}
}

func (matcher *HaveKeyMatcher) Match(actual interface{}) (success bool, err error) {
	if !isMap(actual) {
		return false, fmt.Errorf("HaveKey matcher expects a map.  Got:%s", format.Object(actual, 1))
	}

	keyMatcher, keyIsMatcher := matcher.Key.(omegaMatcher)
	if !keyIsMatcher {
		keyMatcher = &EqualMatcher{Expected: matcher.Key}
	}

	keys := reflect.ValueOf(actual).MapKeys()
	for i := 0; i < len(keys); i++ {
		success, err := keyMatcher.Match(keys[i].Interface())
		if err != nil {
			return false, fmt.Errorf("HaveKey's key matcher failed with:\n%s%s", format.Indent, err.Error())
		}
		if success {
			return true, nil
		}
	}

	return false, nil
}

func (matcher *HaveKeyMatcher) FailureMessage(actual interface{}) (message string) {
	switch matcher.Key.(type) {
	case omegaMatcher:
		return format.Message(actual, "to have key matching", matcher.Key)
	default:
		return format.Message(actual, "to have key", matcher.Key)
	}
}

func (matcher *HaveKeyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	switch matcher.Key.(type) {
	case omegaMatcher:
		return format.Message(actual, "not to have key matching", matcher.Key)
	default:
		return format.Message(actual, "not to have key", matcher.Key)
	}
}
