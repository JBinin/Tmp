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
// Package jsonlog provides helper functions to parse and print time (time.Time) as JSON.
package jsonlog

import (
	"errors"
	"time"
)

const (
	// RFC3339NanoFixed is our own version of RFC339Nano because we want one
	// that pads the nano seconds part with zeros to ensure
	// the timestamps are aligned in the logs.
	RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"
	// JSONFormat is the format used by FastMarshalJSON
	JSONFormat = `"` + time.RFC3339Nano + `"`
)

// FastTimeMarshalJSON avoids one of the extra allocations that
// time.MarshalJSON is making.
func FastTimeMarshalJSON(t time.Time) (string, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return "", errors.New("time.MarshalJSON: year outside of range [0,9999]")
	}
	return t.Format(JSONFormat), nil
}
