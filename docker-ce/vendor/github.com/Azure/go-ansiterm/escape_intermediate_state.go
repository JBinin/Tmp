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

type escapeIntermediateState struct {
	baseState
}

func (escState escapeIntermediateState) Handle(b byte) (s state, e error) {
	logger.Infof("escapeIntermediateState::Handle %#x", b)
	nextState, err := escState.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case sliceContains(intermeds, b):
		return escState, escState.parser.collectInter()
	case sliceContains(executors, b):
		return escState, escState.parser.execute()
	case sliceContains(escapeIntermediateToGroundBytes, b):
		return escState.parser.ground, nil
	}

	return escState, nil
}

func (escState escapeIntermediateState) Transition(s state) error {
	logger.Infof("escapeIntermediateState::Transition %s --> %s", escState.Name(), s.Name())
	escState.baseState.Transition(s)

	switch s {
	case escState.parser.ground:
		return escState.parser.escDispatch()
	}

	return nil
}
