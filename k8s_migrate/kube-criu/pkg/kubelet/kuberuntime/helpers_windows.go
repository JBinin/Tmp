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
// +build windows

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

package kuberuntime

import (
	"github.com/docker/docker/pkg/sysinfo"
)

const (
	// Taken from https://docs.microsoft.com/en-us/virtualization/windowscontainers/manage-containers/resource-controls
	minSharesProcess = 5000
	minSharesHyperV  = 10
	maxShares        = 10000
	milliCPUToCPU    = 1000
)

// milliCPUToShares converts milliCPU to CPU shares
func milliCPUToShares(milliCPU int64, hyperv bool) int64 {
	var minShares int64 = minSharesProcess
	if hyperv {
		minShares = minSharesHyperV
	}

	if milliCPU == 0 {
		// Return here to really match kernel default for zero milliCPU.
		return minShares
	}

	// Conceptually (milliCPU / milliCPUToCPU) * sharesPerCPU, but factored to improve rounding.
	totalCPU := sysinfo.NumCPU()
	shares := (milliCPU * (maxShares - minShares)) / int64(totalCPU) / milliCPUToCPU
	if shares < minShares {
		return minShares
	}
	if shares > maxShares {
		return maxShares
	}
	return shares
}
