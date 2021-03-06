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
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcd3

import (
	"testing"
)

func TestGetReuseDurationSeconds(t *testing.T) {
	testCases := []struct {
		ttl      int64
		duration int64
	}{
		{
			ttl:      3600,
			duration: 60,
		},
		{
			ttl:      1000,
			duration: 50,
		},
	}
	lm := newDefaultLeaseManager(nil)
	for i := 0; i < len(testCases); i++ {
		dur := lm.getReuseDurationSecondsLocked(testCases[i].ttl)
		if dur != testCases[i].duration {
			t.Errorf("Duration error: ttl %v, expected duration %v, get %v\n", testCases[i].ttl, testCases[i].duration, dur)
		}
	}
}
