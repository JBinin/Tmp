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
	"github.com/onsi/gomega/internal/oraclematcher"
	"github.com/onsi/gomega/types"
)

type AndMatcher struct {
	Matchers []types.GomegaMatcher

	// state
	firstFailedMatcher types.GomegaMatcher
}

func (m *AndMatcher) Match(actual interface{}) (success bool, err error) {
	m.firstFailedMatcher = nil
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if !success || err != nil {
			m.firstFailedMatcher = matcher
			return false, err
		}
	}
	return true, nil
}

func (m *AndMatcher) FailureMessage(actual interface{}) (message string) {
	return m.firstFailedMatcher.FailureMessage(actual)
}

func (m *AndMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	// not the most beautiful list of matchers, but not bad either...
	return format.Message(actual, fmt.Sprintf("To not satisfy all of these matchers: %s", m.Matchers))
}

func (m *AndMatcher) MatchMayChangeInTheFuture(actual interface{}) bool {
	/*
		Example with 3 matchers: A, B, C

		Match evaluates them: T, F, <?>  => F
		So match is currently F, what should MatchMayChangeInTheFuture() return?
		Seems like it only depends on B, since currently B MUST change to allow the result to become T

		Match eval: T, T, T  => T
		So match is currently T, what should MatchMayChangeInTheFuture() return?
		Seems to depend on ANY of them being able to change to F.
	*/

	if m.firstFailedMatcher == nil {
		// so all matchers succeeded.. Any one of them changing would change the result.
		for _, matcher := range m.Matchers {
			if oraclematcher.MatchMayChangeInTheFuture(matcher, actual) {
				return true
			}
		}
		return false // none of were going to change
	} else {
		// one of the matchers failed.. it must be able to change in order to affect the result
		return oraclematcher.MatchMayChangeInTheFuture(m.firstFailedMatcher, actual)
	}
}
