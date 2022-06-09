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
package gcfg

import (
	"gopkg.in/warnings.v0"
)

// FatalOnly filters the results of a Read*Into invocation and returns only
// fatal errors. That is, errors (warnings) indicating data for unknown
// sections / variables is ignored. Example invocation:
//
//  err := gcfg.FatalOnly(gcfg.ReadFileInto(&cfg, configFile))
//  if err != nil {
//      ...
//
func FatalOnly(err error) error {
	return warnings.FatalOnly(err)
}

func isFatal(err error) bool {
	_, ok := err.(extraData)
	return !ok
}

type extraData struct {
	section    string
	subsection *string
	variable   *string
}

func (e extraData) Error() string {
	s := "can't store data at section \"" + e.section + "\""
	if e.subsection != nil {
		s += ", subsection \"" + *e.subsection + "\""
	}
	if e.variable != nil {
		s += ", variable \"" + *e.variable + "\""
	}
	return s
}

var _ error = extraData{}
