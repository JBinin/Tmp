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
package serror

import (
	"encoding/json"
	"fmt"
	"strings"
)

var kindLookupMap map[string]StorageOSErrorKind

func init() {
	kindLookupMap = make(map[string]StorageOSErrorKind)

	// Populate the lookup map with all the known constants
	for i := StorageOSErrorKind(0); !strings.HasPrefix(i.String(), "StorageOSErrorKind("); i++ {
		kindLookupMap[i.String()] = i
	}
}

func (s *StorageOSErrorKind) UnmarshalJSON(b []byte) error {
	str := ""
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	v, ok := kindLookupMap[str]
	if !ok {
		return fmt.Errorf("Failed to unmarshal ErrorKind %s", s)
	}

	*s = v
	return nil
}

func (s *StorageOSErrorKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
