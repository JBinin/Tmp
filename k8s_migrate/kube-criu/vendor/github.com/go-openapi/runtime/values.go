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
package runtime

// Values typically represent parameters on a http request.
type Values map[string][]string

// GetOK returns the values collection for the given key.
// When the key is present in the map it will return true for hasKey.
// When the value is not empty it will return true for hasValue.
func (v Values) GetOK(key string) (value []string, hasKey bool, hasValue bool) {
	value, hasKey = v[key]
	if !hasKey {
		return
	}
	if len(value) == 0 {
		return
	}
	hasValue = true
	return
}
