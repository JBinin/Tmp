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
package gstruct

import (
	"github.com/onsi/gomega/types"
)

//Ignore ignores the actual value and always succeeds.
//  Expect(nil).To(Ignore())
//  Expect(true).To(Ignore())
func Ignore() types.GomegaMatcher {
	return &IgnoreMatcher{true}
}

//Reject ignores the actual value and always fails. It can be used in conjunction with IgnoreMissing
//to catch problematic elements, or to verify tests are running.
//  Expect(nil).NotTo(Reject())
//  Expect(true).NotTo(Reject())
func Reject() types.GomegaMatcher {
	return &IgnoreMatcher{false}
}

// A matcher that either always succeeds or always fails.
type IgnoreMatcher struct {
	Succeed bool
}

func (m *IgnoreMatcher) Match(actual interface{}) (bool, error) {
	return m.Succeed, nil
}

func (m *IgnoreMatcher) FailureMessage(_ interface{}) (message string) {
	return "Unconditional failure"
}

func (m *IgnoreMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return "Unconditional success"
}
