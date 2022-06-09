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
package shakers

import (
	"github.com/go-check/check"
)

// True checker verifies the obtained value is true
//
//    c.Assert(myBool, True)
//
var True check.Checker = &boolChecker{
	&check.CheckerInfo{
		Name:   "True",
		Params: []string{"obtained"},
	},
	true,
}

// False checker verifies the obtained value is false
//
//    c.Assert(myBool, False)
//
var False check.Checker = &boolChecker{
	&check.CheckerInfo{
		Name:   "False",
		Params: []string{"obtained"},
	},
	false,
}

type boolChecker struct {
	*check.CheckerInfo
	expected bool
}

func (checker *boolChecker) Check(params []interface{}, names []string) (bool, string) {
	return is(checker.expected, params[0])
}

func is(expected bool, obtained interface{}) (bool, string) {
	obtainedBool, ok := obtained.(bool)
	if !ok {
		return false, "obtained value must be a bool."
	}
	return obtainedBool == expected, ""
}
