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
// Copyright 2015 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"errors"
	"fmt"
)

// Scan implements sql.Scanner so UUIDs can be read from databases transparently
// Currently, database types that map to string and []byte are supported. Please
// consult database-specific driver documentation for matching types.
func (uuid *UUID) Scan(src interface{}) error {
	switch src.(type) {
	case string:
		// if an empty UUID comes from a table, we return a null UUID
		if src.(string) == "" {
			return nil
		}

		// see uuid.Parse for required string format
		parsed := Parse(src.(string))

		if parsed == nil {
			return errors.New("Scan: invalid UUID format")
		}

		*uuid = parsed
	case []byte:
		b := src.([]byte)

		// if an empty UUID comes from a table, we return a null UUID
		if len(b) == 0 {
			return nil
		}

		// assumes a simple slice of bytes if 16 bytes
		// otherwise attempts to parse
		if len(b) == 16 {
			*uuid = UUID(b)
		} else {
			u := Parse(string(b))

			if u == nil {
				return errors.New("Scan: invalid UUID format")
			}

			*uuid = u
		}

	default:
		return fmt.Errorf("Scan: unable to scan type %T into UUID", src)
	}

	return nil
}
