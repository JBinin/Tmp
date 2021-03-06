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
/*
Package date provides time.Time derivatives that conform to the Swagger.io (https://swagger.io/)
defined date   formats: Date and DateTime. Both types may, in most cases, be used in lieu of
time.Time types. And both convert to time.Time through a ToTime method.
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
	"fmt"
	"time"
)

const (
	fullDate     = "2006-01-02"
	fullDateJSON = `"2006-01-02"`
	dateFormat   = "%04d-%02d-%02d"
	jsonFormat   = `"%04d-%02d-%02d"`
)

// Date defines a type similar to time.Time but assumes a layout of RFC3339 full-date (i.e.,
// 2006-01-02).
type Date struct {
	time.Time
}

// ParseDate create a new Date from the passed string.
func ParseDate(date string) (d Date, err error) {
	return parseDate(date, fullDate)
}

func parseDate(date string, format string) (Date, error) {
	d, err := time.Parse(format, date)
	return Date{Time: d}, err
}

// MarshalBinary preserves the Date as a byte array conforming to RFC3339 full-date (i.e.,
// 2006-01-02).
func (d Date) MarshalBinary() ([]byte, error) {
	return d.MarshalText()
}

// UnmarshalBinary reconstitutes a Date saved as a byte array conforming to RFC3339 full-date (i.e.,
// 2006-01-02).
func (d *Date) UnmarshalBinary(data []byte) error {
	return d.UnmarshalText(data)
}

// MarshalJSON preserves the Date as a JSON string conforming to RFC3339 full-date (i.e.,
// 2006-01-02).
func (d Date) MarshalJSON() (json []byte, err error) {
	return []byte(fmt.Sprintf(jsonFormat, d.Year(), d.Month(), d.Day())), nil
}

// UnmarshalJSON reconstitutes the Date from a JSON string conforming to RFC3339 full-date (i.e.,
// 2006-01-02).
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	d.Time, err = time.Parse(fullDateJSON, string(data))
	return err
}

// MarshalText preserves the Date as a byte array conforming to RFC3339 full-date (i.e.,
// 2006-01-02).
func (d Date) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf(dateFormat, d.Year(), d.Month(), d.Day())), nil
}

// UnmarshalText reconstitutes a Date saved as a byte array conforming to RFC3339 full-date (i.e.,
// 2006-01-02).
func (d *Date) UnmarshalText(data []byte) (err error) {
	d.Time, err = time.Parse(fullDate, string(data))
	return err
}

// String returns the Date formatted as an RFC3339 full-date string (i.e., 2006-01-02).
func (d Date) String() string {
	return fmt.Sprintf(dateFormat, d.Year(), d.Month(), d.Day())
}

// ToTime returns a Date as a time.Time
func (d Date) ToTime() time.Time {
	return d.Time
}
