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

package configs

import "syscall"

func (n *Namespace) Syscall() int {
	return namespaceInfo[n.Type]
}

var namespaceInfo = map[NamespaceType]int{
	NEWNET:  syscall.CLONE_NEWNET,
	NEWNS:   syscall.CLONE_NEWNS,
	NEWUSER: syscall.CLONE_NEWUSER,
	NEWIPC:  syscall.CLONE_NEWIPC,
	NEWUTS:  syscall.CLONE_NEWUTS,
	NEWPID:  syscall.CLONE_NEWPID,
}

// CloneFlags parses the container's Namespaces options to set the correct
// flags on clone, unshare. This function returns flags only for new namespaces.
func (n *Namespaces) CloneFlags() uintptr {
	var flag int
	for _, v := range *n {
		if v.Path != "" {
			continue
		}
		flag |= namespaceInfo[v.Type]
	}
	return uintptr(flag)
}
