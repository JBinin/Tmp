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

type BeEquivalentToMatcher struct {
	Expected interface{}
}

func (matcher *BeEquivalentToMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Both actual and expected must not be nil.")
	}

	convertedActual := actual

	if actual != nil && matcher.Expected != nil && reflect.TypeOf(actual).ConvertibleTo(reflect.TypeOf(matcher.Expected)) {
		convertedActual = reflect.ValueOf(actual).Convert(reflect.TypeOf(matcher.Expected)).Interface()
	}

	return reflect.DeepEqual(convertedActual, matcher.Expected), nil
}

func (matcher *BeEquivalentToMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be equivalent to", matcher.Expected)
}

func (matcher *BeEquivalentToMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be equivalent to", matcher.Expected)
}
