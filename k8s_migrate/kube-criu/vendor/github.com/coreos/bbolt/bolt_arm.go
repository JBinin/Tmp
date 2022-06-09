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
package bolt

import "unsafe"

// maxMapSize represents the largest mmap size supported by Bolt.
const maxMapSize = 0x7FFFFFFF // 2GB

// maxAllocSize is the size used when creating array pointers.
const maxAllocSize = 0xFFFFFFF

// Are unaligned load/stores broken on this arch?
var brokenUnaligned bool

func init() {
	// Simple check to see whether this arch handles unaligned load/stores
	// correctly.

	// ARM9 and older devices require load/stores to be from/to aligned
	// addresses. If not, the lower 2 bits are cleared and that address is
	// read in a jumbled up order.

	// See http://infocenter.arm.com/help/index.jsp?topic=/com.arm.doc.faqs/ka15414.html

	raw := [6]byte{0xfe, 0xef, 0x11, 0x22, 0x22, 0x11}
	val := *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&raw)) + 2))

	brokenUnaligned = val != 0x11222211
}
