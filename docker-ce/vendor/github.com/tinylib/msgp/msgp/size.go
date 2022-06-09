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
package msgp

// The sizes provided
// are the worst-case
// encoded sizes for
// each type. For variable-
// length types ([]byte, string),
// the total encoded size is
// the prefix size plus the
// length of the object.
const (
	Int64Size      = 9
	IntSize        = Int64Size
	UintSize       = Int64Size
	Int8Size       = 2
	Int16Size      = 3
	Int32Size      = 5
	Uint8Size      = 2
	ByteSize       = Uint8Size
	Uint16Size     = 3
	Uint32Size     = 5
	Uint64Size     = Int64Size
	Float64Size    = 9
	Float32Size    = 5
	Complex64Size  = 10
	Complex128Size = 18

	TimeSize = 15
	BoolSize = 1
	NilSize  = 1

	MapHeaderSize   = 5
	ArrayHeaderSize = 5

	BytesPrefixSize     = 5
	StringPrefixSize    = 5
	ExtensionPrefixSize = 6
)
