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
//+build go1.9

package reflect2

import (
	"unsafe"
)

//go:linkname makemap reflect.makemap
func makemap(rtype unsafe.Pointer, cap int) (m unsafe.Pointer)

func makeMapWithSize(rtype unsafe.Pointer, cap int) unsafe.Pointer {
	return makemap(rtype, cap)
}