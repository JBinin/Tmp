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

type groundState struct {
	baseState
}

func (gs groundState) Handle(b byte) (s state, e error) {
	gs.parser.context.currentChar = b

	nextState, err := gs.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case sliceContains(printables, b):
		return gs, gs.parser.print()

	case sliceContains(executors, b):
		return gs, gs.parser.execute()
	}

	return gs, nil
}
