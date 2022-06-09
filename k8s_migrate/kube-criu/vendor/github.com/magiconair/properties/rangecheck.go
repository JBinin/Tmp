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
// Copyright 2016 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"fmt"
	"math"
)

// make this a var to overwrite it in a test
var is32Bit = ^uint(0) == math.MaxUint32

// intRangeCheck checks if the value fits into the int type and
// panics if it does not.
func intRangeCheck(key string, v int64) int {
	if is32Bit && (v < math.MinInt32 || v > math.MaxInt32) {
		panic(fmt.Sprintf("Value %d for key %s out of range", v, key))
	}
	return int(v)
}

// uintRangeCheck checks if the value fits into the uint type and
// panics if it does not.
func uintRangeCheck(key string, v uint64) uint {
	if is32Bit && v > math.MaxUint32 {
		panic(fmt.Sprintf("Value %d for key %s out of range", v, key))
	}
	return uint(v)
}
