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
package jsoniter

// ReadArray read array element, tells if the array has more element to read.
func (iter *Iterator) ReadArray() (ret bool) {
	c := iter.nextToken()
	switch c {
	case 'n':
		iter.skipThreeBytes('u', 'l', 'l')
		return false // null
	case '[':
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			return true
		}
		return false
	case ']':
		return false
	case ',':
		return true
	default:
		iter.ReportError("ReadArray", "expect [ or , or ] or n, but found "+string([]byte{c}))
		return
	}
}

// ReadArrayCB read array with callback
func (iter *Iterator) ReadArrayCB(callback func(*Iterator) bool) (ret bool) {
	c := iter.nextToken()
	if c == '[' {
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			if !callback(iter) {
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				if !callback(iter) {
					return false
				}
				c = iter.nextToken()
			}
			if c != ']' {
				iter.ReportError("ReadArrayCB", "expect ] in the end, but found "+string([]byte{c}))
				return false
			}
			return true
		}
		return true
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	iter.ReportError("ReadArrayCB", "expect [ or n, but found "+string([]byte{c}))
	return false
}
