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
package testrunner

type RunResult struct {
	Passed               bool
	HasProgrammaticFocus bool
}

func PassingRunResult() RunResult {
	return RunResult{
		Passed:               true,
		HasProgrammaticFocus: false,
	}
}

func FailingRunResult() RunResult {
	return RunResult{
		Passed:               false,
		HasProgrammaticFocus: false,
	}
}

func (r RunResult) Merge(o RunResult) RunResult {
	return RunResult{
		Passed:               r.Passed && o.Passed,
		HasProgrammaticFocus: r.HasProgrammaticFocus || o.HasProgrammaticFocus,
	}
}
