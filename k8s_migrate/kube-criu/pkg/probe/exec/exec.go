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
Copyright 2015 The Kubernetes Authors.

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

package exec

import (
	"k8s.io/kubernetes/pkg/probe"
	"k8s.io/utils/exec"

	"github.com/golang/glog"
)

func New() ExecProber {
	return execProber{}
}

type ExecProber interface {
	Probe(e exec.Cmd) (probe.Result, string, error)
}

type execProber struct{}

func (pr execProber) Probe(e exec.Cmd) (probe.Result, string, error) {
	data, err := e.CombinedOutput()
	glog.V(4).Infof("Exec probe response: %q", string(data))
	if err != nil {
		exit, ok := err.(exec.ExitError)
		if ok {
			if exit.ExitStatus() == 0 {
				return probe.Success, string(data), nil
			} else {
				return probe.Failure, string(data), nil
			}
		}
		return probe.Unknown, "", err
	}
	return probe.Success, string(data), nil
}
