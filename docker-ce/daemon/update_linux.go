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

package daemon

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/libcontainerd"
)

func toContainerdResources(resources container.Resources) libcontainerd.Resources {
	var r libcontainerd.Resources
	r.BlkioWeight = uint64(resources.BlkioWeight)
	r.CpuShares = uint64(resources.CPUShares)
	r.CpuPeriod = uint64(resources.CPUPeriod)
	r.CpuQuota = uint64(resources.CPUQuota)
	r.CpusetCpus = resources.CpusetCpus
	r.CpusetMems = resources.CpusetMems
	r.MemoryLimit = uint64(resources.Memory)
	if resources.MemorySwap > 0 {
		r.MemorySwap = uint64(resources.MemorySwap)
	}
	r.MemoryReservation = uint64(resources.MemoryReservation)
	r.KernelMemoryLimit = uint64(resources.KernelMemory)
	return r
}
