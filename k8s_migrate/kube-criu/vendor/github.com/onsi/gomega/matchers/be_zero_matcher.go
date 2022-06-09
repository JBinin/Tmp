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
	"github.com/onsi/gomega/format"
	"reflect"
)

type BeZeroMatcher struct {
}

func (matcher *BeZeroMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return true, nil
	}
	zeroValue := reflect.Zero(reflect.TypeOf(actual)).Interface()

	return reflect.DeepEqual(zeroValue, actual), nil

}

func (matcher *BeZeroMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be zero-valued")
}

func (matcher *BeZeroMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be zero-valued")
}
