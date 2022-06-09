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

type ItNode struct {
	runner *runner

	flag types.FlagType
	text string
}

func NewItNode(text string, body interface{}, flag types.FlagType, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer, componentIndex int) *ItNode {
	return &ItNode{
		runner: newRunner(body, codeLocation, timeout, failer, types.SpecComponentTypeIt, componentIndex),
		flag:   flag,
		text:   text,
	}
}

func (node *ItNode) Run() (outcome types.SpecState, failure types.SpecFailure) {
	return node.runner.run()
}

func (node *ItNode) Type() types.SpecComponentType {
	return types.SpecComponentTypeIt
}

func (node *ItNode) Text() string {
	return node.text
}

func (node *ItNode) Flag() types.FlagType {
	return node.flag
}

func (node *ItNode) CodeLocation() types.CodeLocation {
	return node.runner.codeLocation
}

func (node *ItNode) Samples() int {
	return 1
}
