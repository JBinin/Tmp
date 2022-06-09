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
	"reflect"
)

type MeasureNode struct {
	runner *runner

	text        string
	flag        types.FlagType
	samples     int
	benchmarker *benchmarker
}

func NewMeasureNode(text string, body interface{}, flag types.FlagType, codeLocation types.CodeLocation, samples int, failer *failer.Failer, componentIndex int) *MeasureNode {
	benchmarker := newBenchmarker()

	wrappedBody := func() {
		reflect.ValueOf(body).Call([]reflect.Value{reflect.ValueOf(benchmarker)})
	}

	return &MeasureNode{
		runner: newRunner(wrappedBody, codeLocation, 0, failer, types.SpecComponentTypeMeasure, componentIndex),

		text:        text,
		flag:        flag,
		samples:     samples,
		benchmarker: benchmarker,
	}
}

func (node *MeasureNode) Run() (outcome types.SpecState, failure types.SpecFailure) {
	return node.runner.run()
}

func (node *MeasureNode) MeasurementsReport() map[string]*types.SpecMeasurement {
	return node.benchmarker.measurementsReport()
}

func (node *MeasureNode) Type() types.SpecComponentType {
	return types.SpecComponentTypeMeasure
}

func (node *MeasureNode) Text() string {
	return node.text
}

func (node *MeasureNode) Flag() types.FlagType {
	return node.flag
}

func (node *MeasureNode) CodeLocation() types.CodeLocation {
	return node.runner.codeLocation
}

func (node *MeasureNode) Samples() int {
	return node.samples
}
