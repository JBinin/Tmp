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
package date

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"errors"
	"time"
)

const (
	rfc1123JSON = `"` + time.RFC1123 + `"`
	rfc1123     = time.RFC1123
)

// TimeRFC1123 defines a type similar to time.Time but assumes a layout of RFC1123 date-time (i.e.,
// Mon, 02 Jan 2006 15:04:05 MST).
type TimeRFC1123 struct {
	time.Time
}

// UnmarshalJSON reconstitutes the Time from a JSON string conforming to RFC1123 date-time
// (i.e., Mon, 02 Jan 2006 15:04:05 MST).
func (t *TimeRFC1123) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = ParseTime(rfc1123JSON, string(data))
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON preserves the Time as a JSON string conforming to RFC1123 date-time (i.e.,
// Mon, 02 Jan 2006 15:04:05 MST).
func (t TimeRFC1123) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	b := []byte(t.Format(rfc1123JSON))
	return b, nil
}

// MarshalText preserves the Time as a byte array conforming to RFC1123 date-time (i.e.,
// Mon, 02 Jan 2006 15:04:05 MST).
func (t TimeRFC1123) MarshalText() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalText: year outside of range [0,9999]")
	}

	b := []byte(t.Format(rfc1123))
	return b, nil
}

// UnmarshalText reconstitutes a Time saved as a byte array conforming to RFC1123 date-time
// (i.e., Mon, 02 Jan 2006 15:04:05 MST).
func (t *TimeRFC1123) UnmarshalText(data []byte) (err error) {
	t.Time, err = ParseTime(rfc1123, string(data))
	if err != nil {
		return err
	}
	return nil
}

// MarshalBinary preserves the Time as a byte array conforming to RFC1123 date-time (i.e.,
// Mon, 02 Jan 2006 15:04:05 MST).
func (t TimeRFC1123) MarshalBinary() ([]byte, error) {
	return t.MarshalText()
}

// UnmarshalBinary reconstitutes a Time saved as a byte array conforming to RFC1123 date-time
// (i.e., Mon, 02 Jan 2006 15:04:05 MST).
func (t *TimeRFC1123) UnmarshalBinary(data []byte) error {
	return t.UnmarshalText(data)
}

// ToTime returns a Time as a time.Time
func (t TimeRFC1123) ToTime() time.Time {
	return t.Time
}

// String returns the Time formatted as an RFC1123 date-time string (i.e.,
// Mon, 02 Jan 2006 15:04:05 MST).
func (t TimeRFC1123) String() string {
	// Note: time.Time.String does not return an RFC1123 compliant string, time.Time.MarshalText does.
	b, err := t.MarshalText()
	if err != nil {
		return ""
	}
	return string(b)
}
