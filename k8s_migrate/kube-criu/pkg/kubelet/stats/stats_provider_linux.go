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
// +build linux

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

package stats

import (
	"io/ioutil"
	"strconv"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

func (p *StatsProvider) RlimitStats() (*statsapi.RlimitStats, error) {
	rlimit := &statsapi.RlimitStats{}

	if content, err := ioutil.ReadFile("/proc/sys/kernel/pid_max"); err == nil {
		if maxPid, err := strconv.ParseInt(string(content[:len(content)-1]), 10, 64); err == nil {
			rlimit.MaxPID = &maxPid
		}
	}

	var info syscall.Sysinfo_t
	syscall.Sysinfo(&info)
	procs := int64(info.Procs)
	rlimit.NumOfRunningProcesses = &procs

	rlimit.Time = v1.NewTime(time.Now())

	return rlimit, nil
}
