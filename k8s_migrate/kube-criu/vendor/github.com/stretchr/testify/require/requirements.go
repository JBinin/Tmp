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
package require

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

type tHelper interface {
	Helper()
}

// ComparisonAssertionFunc is a common function prototype when comparing two values.  Can be useful
// for table driven tests.
type ComparisonAssertionFunc func(TestingT, interface{}, interface{}, ...interface{})

// ValueAssertionFunc is a common function prototype when validating a single value.  Can be useful
// for table driven tests.
type ValueAssertionFunc func(TestingT, interface{}, ...interface{})

// BoolAssertionFunc is a common function prototype when validating a bool value.  Can be useful
// for table driven tests.
type BoolAssertionFunc func(TestingT, bool, ...interface{})

// ValuesAssertionFunc is a common function prototype when validating an error value.  Can be useful
// for table driven tests.
type ErrorAssertionFunc func(TestingT, error, ...interface{})

//go:generate go run ../_codegen/main.go -output-package=require -template=require.go.tmpl -include-format-funcs
