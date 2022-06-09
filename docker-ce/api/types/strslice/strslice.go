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
package strslice

import "encoding/json"

// StrSlice represents a string or an array of strings.
// We need to override the json decoder to accept both options.
type StrSlice []string

// UnmarshalJSON decodes the byte slice whether it's a string or an array of
// strings. This method is needed to implement json.Unmarshaler.
func (e *StrSlice) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		// With no input, we preserve the existing value by returning nil and
		// leaving the target alone. This allows defining default values for
		// the type.
		return nil
	}

	p := make([]string, 0, 1)
	if err := json.Unmarshal(b, &p); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p = append(p, s)
	}

	*e = p
	return nil
}
