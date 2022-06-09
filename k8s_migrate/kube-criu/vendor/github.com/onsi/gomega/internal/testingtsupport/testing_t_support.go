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
package testingtsupport

import (
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/onsi/gomega/types"
)

type gomegaTestingT interface {
	Errorf(format string, args ...interface{})
}

func BuildTestingTGomegaFailHandler(t gomegaTestingT) types.GomegaFailHandler {
	return func(message string, callerSkip ...int) {
		skip := 1
		if len(callerSkip) > 0 {
			skip = callerSkip[0]
		}
		stackTrace := pruneStack(string(debug.Stack()), skip)
		t.Errorf("\n%s\n%s", stackTrace, message)
	}
}

func pruneStack(fullStackTrace string, skip int) string {
	stack := strings.Split(fullStackTrace, "\n")
	if len(stack) > 2*(skip+1) {
		stack = stack[2*(skip+1):]
	}
	prunedStack := []string{}
	re := regexp.MustCompile(`\/ginkgo\/|\/pkg\/testing\/|\/pkg\/runtime\/`)
	for i := 0; i < len(stack)/2; i++ {
		if !re.Match([]byte(stack[i*2])) {
			prunedStack = append(prunedStack, stack[i*2])
			prunedStack = append(prunedStack, stack[i*2+1])
		}
	}
	return strings.Join(prunedStack, "\n")
}
