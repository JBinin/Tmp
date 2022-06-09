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

import "github.com/onsi/gomega/format"

type BeNilMatcher struct {
}

func (matcher *BeNilMatcher) Match(actual interface{}) (success bool, err error) {
	return isNil(actual), nil
}

func (matcher *BeNilMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be nil")
}

func (matcher *BeNilMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be nil")
}
