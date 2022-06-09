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

type SetupNode struct {
	runner *runner
}

func (node *SetupNode) Run() (outcome types.SpecState, failure types.SpecFailure) {
	return node.runner.run()
}

func (node *SetupNode) Type() types.SpecComponentType {
	return node.runner.nodeType
}

func (node *SetupNode) CodeLocation() types.CodeLocation {
	return node.runner.codeLocation
}

func NewBeforeEachNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer, componentIndex int) *SetupNode {
	return &SetupNode{
		runner: newRunner(body, codeLocation, timeout, failer, types.SpecComponentTypeBeforeEach, componentIndex),
	}
}

func NewAfterEachNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer, componentIndex int) *SetupNode {
	return &SetupNode{
		runner: newRunner(body, codeLocation, timeout, failer, types.SpecComponentTypeAfterEach, componentIndex),
	}
}

func NewJustBeforeEachNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer, componentIndex int) *SetupNode {
	return &SetupNode{
		runner: newRunner(body, codeLocation, timeout, failer, types.SpecComponentTypeJustBeforeEach, componentIndex),
	}
}
