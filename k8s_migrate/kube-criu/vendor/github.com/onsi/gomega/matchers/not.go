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
	"github.com/onsi/gomega/internal/oraclematcher"
	"github.com/onsi/gomega/types"
)

type NotMatcher struct {
	Matcher types.GomegaMatcher
}

func (m *NotMatcher) Match(actual interface{}) (bool, error) {
	success, err := m.Matcher.Match(actual)
	if err != nil {
		return false, err
	}
	return !success, nil
}

func (m *NotMatcher) FailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(actual) // works beautifully
}

func (m *NotMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.FailureMessage(actual) // works beautifully
}

func (m *NotMatcher) MatchMayChangeInTheFuture(actual interface{}) bool {
	return oraclematcher.MatchMayChangeInTheFuture(m.Matcher, actual) // just return m.Matcher's value
}
