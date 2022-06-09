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
package leafnodes

import (
	"github.com/onsi/ginkgo/internal/failer"
	"github.com/onsi/ginkgo/types"
	"time"
)

type SuiteNode interface {
	Run(parallelNode int, parallelTotal int, syncHost string) bool
	Passed() bool
	Summary() *types.SetupSummary
}

type simpleSuiteNode struct {
	runner  *runner
	outcome types.SpecState
	failure types.SpecFailure
	runTime time.Duration
}

func (node *simpleSuiteNode) Run(parallelNode int, parallelTotal int, syncHost string) bool {
	t := time.Now()
	node.outcome, node.failure = node.runner.run()
	node.runTime = time.Since(t)

	return node.outcome == types.SpecStatePassed
}

func (node *simpleSuiteNode) Passed() bool {
	return node.outcome == types.SpecStatePassed
}

func (node *simpleSuiteNode) Summary() *types.SetupSummary {
	return &types.SetupSummary{
		ComponentType: node.runner.nodeType,
		CodeLocation:  node.runner.codeLocation,
		State:         node.outcome,
		RunTime:       node.runTime,
		Failure:       node.failure,
	}
}

func NewBeforeSuiteNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer) SuiteNode {
	return &simpleSuiteNode{
		runner: newRunner(body, codeLocation, timeout, failer, types.SpecComponentTypeBeforeSuite, 0),
	}
}

func NewAfterSuiteNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer) SuiteNode {
	return &simpleSuiteNode{
		runner: newRunner(body, codeLocation, timeout, failer, types.SpecComponentTypeAfterSuite, 0),
	}
}
