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

import (
	"strconv"
)

func parseParams(bytes []byte) ([]string, error) {
	paramBuff := make([]byte, 0, 0)
	params := []string{}

	for _, v := range bytes {
		if v == ';' {
			if len(paramBuff) > 0 {
				// Completed parameter, append it to the list
				s := string(paramBuff)
				params = append(params, s)
				paramBuff = make([]byte, 0, 0)
			}
		} else {
			paramBuff = append(paramBuff, v)
		}
	}

	// Last parameter may not be terminated with ';'
	if len(paramBuff) > 0 {
		s := string(paramBuff)
		params = append(params, s)
	}

	logger.Infof("Parsed params: %v with length: %d", params, len(params))
	return params, nil
}

func parseCmd(context ansiContext) (string, error) {
	return string(context.currentChar), nil
}

func getInt(params []string, dflt int) int {
	i := getInts(params, 1, dflt)[0]
	logger.Infof("getInt: %v", i)
	return i
}

func getInts(params []string, minCount int, dflt int) []int {
	ints := []int{}

	for _, v := range params {
		i, _ := strconv.Atoi(v)
		// Zero is mapped to the default value in VT100.
		if i == 0 {
			i = dflt
		}
		ints = append(ints, i)
	}

	if len(ints) < minCount {
		remaining := minCount - len(ints)
		for i := 0; i < remaining; i++ {
			ints = append(ints, dflt)
		}
	}

	logger.Infof("getInts: %v", ints)

	return ints
}

func (ap *AnsiParser) modeDispatch(param string, set bool) error {
	switch param {
	case "?3":
		return ap.eventHandler.DECCOLM(set)
	case "?6":
		return ap.eventHandler.DECOM(set)
	case "?25":
		return ap.eventHandler.DECTCEM(set)
	}
	return nil
}

func (ap *AnsiParser) hDispatch(params []string) error {
	if len(params) == 1 {
		return ap.modeDispatch(params[0], true)
	}

	return nil
}

func (ap *AnsiParser) lDispatch(params []string) error {
	if len(params) == 1 {
		return ap.modeDispatch(params[0], false)
	}

	return nil
}

func getEraseParam(params []string) int {
	param := getInt(params, 0)
	if param < 0 || 3 < param {
		param = 0
	}

	return param
}
