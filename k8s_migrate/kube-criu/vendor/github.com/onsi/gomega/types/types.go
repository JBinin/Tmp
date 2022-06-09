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
package types

type GomegaFailHandler func(message string, callerSkip ...int)

//A simple *testing.T interface wrapper
type GomegaTestingT interface {
	Errorf(format string, args ...interface{})
}

//All Gomega matchers must implement the GomegaMatcher interface
//
//For details on writing custom matchers, check out: http://onsi.github.io/gomega/#adding_your_own_matchers
type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}
