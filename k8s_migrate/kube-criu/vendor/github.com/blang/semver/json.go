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
package semver

import (
	"encoding/json"
)

// MarshalJSON implements the encoding/json.Marshaler interface.
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (v *Version) UnmarshalJSON(data []byte) (err error) {
	var versionString string

	if err = json.Unmarshal(data, &versionString); err != nil {
		return
	}

	*v, err = Parse(versionString)

	return
}
