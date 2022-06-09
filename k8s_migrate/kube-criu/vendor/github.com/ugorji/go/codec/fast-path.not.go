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
// +build notfastpath

package codec

import "reflect"

const fastpathEnabled = false

// The generated fast-path code is very large, and adds a few seconds to the build time.
// This causes test execution, execution of small tools which use codec, etc
// to take a long time.
//
// To mitigate, we now support the notfastpath tag.
// This tag disables fastpath during build, allowing for faster build, test execution,
// short-program runs, etc.

func fastpathDecodeTypeSwitch(iv interface{}, d *Decoder) bool      { return false }
func fastpathEncodeTypeSwitch(iv interface{}, e *Encoder) bool      { return false }
func fastpathEncodeTypeSwitchSlice(iv interface{}, e *Encoder) bool { return false }
func fastpathEncodeTypeSwitchMap(iv interface{}, e *Encoder) bool   { return false }

type fastpathT struct{}
type fastpathE struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encFnInfo, reflect.Value)
	decfn func(*decFnInfo, reflect.Value)
}
type fastpathA [0]fastpathE

func (x fastpathA) index(rtid uintptr) int { return -1 }

var fastpathAV fastpathA
var fastpathTV fastpathT
