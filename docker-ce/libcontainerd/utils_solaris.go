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
package libcontainerd

import (
	"syscall"

	containerd "github.com/docker/containerd/api/grpc/types"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func getRootIDs(s specs.Spec) (int, int, error) {
	return 0, 0, nil
}

func systemPid(ctr *containerd.Container) uint32 {
	var pid uint32
	for _, p := range ctr.Processes {
		if p.Pid == InitFriendlyName {
			pid = p.SystemPid
		}
	}
	return pid
}

// setPDeathSig sets the parent death signal to SIGKILL
func setSysProcAttr(sid bool) *syscall.SysProcAttr {
	return nil
}
