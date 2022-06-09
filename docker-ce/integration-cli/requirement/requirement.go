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
package requirement

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

type skipT interface {
	Skip(reason string)
}

// Test represent a function that can be used as a requirement validation.
type Test func() bool

// Is checks if the environment satisfies the requirements
// for the test to run or skips the tests.
func Is(s skipT, requirements ...Test) {
	for _, r := range requirements {
		isValid := r()
		if !isValid {
			requirementFunc := runtime.FuncForPC(reflect.ValueOf(r).Pointer()).Name()
			s.Skip(fmt.Sprintf("unmatched requirement %s", extractRequirement(requirementFunc)))
		}
	}
}

func extractRequirement(requirementFunc string) string {
	requirement := path.Base(requirementFunc)
	return strings.SplitN(requirement, ".", 2)[1]
}
