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
package ansiterm

type csiParamState struct {
	baseState
}

func (csiState csiParamState) Handle(b byte) (s state, e error) {
	logger.Infof("CsiParam::Handle %#x", b)

	nextState, err := csiState.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case sliceContains(alphabetics, b):
		return csiState.parser.ground, nil
	case sliceContains(csiCollectables, b):
		csiState.parser.collectParam()
		return csiState, nil
	case sliceContains(executors, b):
		return csiState, csiState.parser.execute()
	}

	return csiState, nil
}

func (csiState csiParamState) Transition(s state) error {
	logger.Infof("CsiParam::Transition %s --> %s", csiState.Name(), s.Name())
	csiState.baseState.Transition(s)

	switch s {
	case csiState.parser.ground:
		return csiState.parser.csiDispatch()
	}

	return nil
}
