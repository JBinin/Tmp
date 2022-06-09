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
package time

import (
	"testing"
	"time"
)

func TestDurationToSecondsString(t *testing.T) {
	cases := []struct {
		in       time.Duration
		expected string
	}{
		{0 * time.Second, "0"},
		{1 * time.Second, "1"},
		{1 * time.Minute, "60"},
		{24 * time.Hour, "86400"},
	}

	for _, c := range cases {
		s := DurationToSecondsString(c.in)
		if s != c.expected {
			t.Errorf("wrong value for input `%v`: expected `%s`, got `%s`", c.in, c.expected, s)
			t.Fail()
		}
	}
}
