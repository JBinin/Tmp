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
package oraclematcher

import "github.com/onsi/gomega/types"

/*
GomegaMatchers that also match the OracleMatcher interface can convey information about
whether or not their result will change upon future attempts.

This allows `Eventually` and `Consistently` to short circuit if success becomes impossible.

For example, a process' exit code can never change.  So, gexec's Exit matcher returns `true`
for `MatchMayChangeInTheFuture` until the process exits, at which point it returns `false` forevermore.
*/
type OracleMatcher interface {
	MatchMayChangeInTheFuture(actual interface{}) bool
}

func MatchMayChangeInTheFuture(matcher types.GomegaMatcher, value interface{}) bool {
	oracleMatcher, ok := matcher.(OracleMatcher)
	if !ok {
		return true
	}

	return oracleMatcher.MatchMayChangeInTheFuture(value)
}
