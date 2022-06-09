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
package objx

// Has gets whether there is something at the specified selector
// or not.
//
// If m is nil, Has will always return false.
func (m Map) Has(selector string) bool {
	if m == nil {
		return false
	}
	return !m.Get(selector).IsNil()
}

// IsNil gets whether the data is nil or not.
func (v *Value) IsNil() bool {
	return v == nil || v.data == nil
}
