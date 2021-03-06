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
Copyright 2017 The Kubernetes Authors.

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

package ipam

import (
	"errors"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	time10s := 10 * time.Second
	time5s := 5 * time.Second
	timeout := &Timeout{
		Resync:       time10s,
		MaxBackoff:   time5s,
		InitialRetry: time.Second,
	}

	for _, testStep := range []struct {
		err  error
		want time.Duration
	}{
		{nil, time10s},
		{nil, time10s},
		{errors.New("x"), time.Second},
		{errors.New("x"), 2 * time.Second},
		{errors.New("x"), 4 * time.Second},
		{errors.New("x"), 5 * time.Second},
		{errors.New("x"), 5 * time.Second},
		{nil, time10s},
		{nil, time10s},
		{errors.New("x"), time.Second},
		{errors.New("x"), 2 * time.Second},
		{nil, time10s},
	} {
		timeout.Update(testStep.err == nil)
		next := timeout.Next()
		if next != testStep.want {
			t.Errorf("timeout.next(%v) = %v, want %v", testStep.err, next, testStep.want)
		}
	}
}
