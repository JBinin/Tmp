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
package registry

import "testing"

func TestLookupV1Endpoints(t *testing.T) {
	s := NewService(ServiceOptions{})

	cases := []struct {
		hostname    string
		expectedLen int
	}{
		{"example.com", 1},
		{DefaultNamespace, 0},
		{DefaultV2Registry.Host, 0},
		{IndexHostname, 0},
	}

	for _, c := range cases {
		if ret, err := s.lookupV1Endpoints(c.hostname); err != nil || len(ret) != c.expectedLen {
			t.Errorf("lookupV1Endpoints(`"+c.hostname+"`) returned %+v and %+v", ret, err)
		}
	}
}
