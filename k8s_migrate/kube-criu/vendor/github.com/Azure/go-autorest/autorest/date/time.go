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
	"regexp"
	"time"
)

// Azure reports time in UTC but it doesn't include the 'Z' time zone suffix in some cases.
const (
	azureUtcFormatJSON = `"2006-01-02T15:04:05.999999999"`
	azureUtcFormat     = "2006-01-02T15:04:05.999999999"
	rfc3339JSON        = `"` + time.RFC3339Nano + `"`
	rfc3339            = time.RFC3339Nano
	tzOffsetRegex      = `(Z|z|\+|-)(\d+:\d+)*"*$`
)

// Time defines a type similar to time.Time but assumes a layout of RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
type Time struct {
	time.Time
}

// MarshalBinary preserves the Time as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (t Time) MarshalBinary() ([]byte, error) {
	return t.Time.MarshalText()
}

// UnmarshalBinary reconstitutes a Time saved as a byte array conforming to RFC3339 date-time
// (i.e., 2006-01-02T15:04:05Z).
func (t *Time) UnmarshalBinary(data []byte) error {
	return t.UnmarshalText(data)
}

// MarshalJSON preserves the Time as a JSON string conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (t Time) MarshalJSON() (json []byte, err error) {
	return t.Time.MarshalJSON()
}

// UnmarshalJSON reconstitutes the Time from a JSON string conforming to RFC3339 date-time
// (i.e., 2006-01-02T15:04:05Z).
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	timeFormat := azureUtcFormatJSON
	match, err := regexp.Match(tzOffsetRegex, data)
	if err != nil {
		return err
	} else if match {
		timeFormat = rfc3339JSON
	}
	t.Time, err = ParseTime(timeFormat, string(data))
	return err
}

// MarshalText preserves the Time as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (t Time) MarshalText() (text []byte, err error) {
	return t.Time.MarshalText()
}

// UnmarshalText reconstitutes a Time saved as a byte array conforming to RFC3339 date-time
// (i.e., 2006-01-02T15:04:05Z).
func (t *Time) UnmarshalText(data []byte) (err error) {
	timeFormat := azureUtcFormat
	match, err := regexp.Match(tzOffsetRegex, data)
	if err != nil {
		return err
	} else if match {
		timeFormat = rfc3339
	}
	t.Time, err = ParseTime(timeFormat, string(data))
	return err
}

// String returns the Time formatted as an RFC3339 date-time string (i.e.,
// 2006-01-02T15:04:05Z).
func (t Time) String() string {
	// Note: time.Time.String does not return an RFC3339 compliant string, time.Time.MarshalText does.
	b, err := t.MarshalText()
	if err != nil {
		return ""
	}
	return string(b)
}

// ToTime returns a Time as a time.Time
func (t Time) ToTime() time.Time {
	return t.Time
}
