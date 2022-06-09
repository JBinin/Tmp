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

package mount

import "k8s.io/utils/exec"

func NewOsExec() Exec {
	return &osExec{}
}

// Real implementation of Exec interface that uses simple util.Exec
type osExec struct{}

var _ Exec = &osExec{}

func (e *osExec) Run(cmd string, args ...string) ([]byte, error) {
	exe := exec.New()
	return exe.Command(cmd, args...).CombinedOutput()
}

func NewFakeExec(run runHook) *FakeExec {
	return &FakeExec{runHook: run}
}

// Fake for testing.
type FakeExec struct {
	runHook runHook
}
type runHook func(cmd string, args ...string) ([]byte, error)

func (f *FakeExec) Run(cmd string, args ...string) ([]byte, error) {
	if f.runHook != nil {
		return f.runHook(cmd, args...)
	}
	return nil, nil
}
