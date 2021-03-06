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
Copyright 2016 The Kubernetes Authors.

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

package sysctl

import (
	"strings"
)

// Namespace represents a kernel namespace name.
type Namespace string

const (
	// the Linux IPC namespace
	IpcNamespace = Namespace("ipc")

	// the network namespace
	NetNamespace = Namespace("net")

	// the zero value if no namespace is known
	UnknownNamespace = Namespace("")
)

var namespaces = map[string]Namespace{
	"kernel.sem": IpcNamespace,
}

var prefixNamespaces = map[string]Namespace{
	"kernel.shm": IpcNamespace,
	"kernel.msg": IpcNamespace,
	"fs.mqueue.": IpcNamespace,
	"net.":       NetNamespace,
}

// NamespacedBy returns the namespace of the Linux kernel for a sysctl, or
// UnknownNamespace if the sysctl is not known to be namespaced.
func NamespacedBy(val string) Namespace {
	if ns, found := namespaces[val]; found {
		return ns
	}
	for p, ns := range prefixNamespaces {
		if strings.HasPrefix(val, p) {
			return ns
		}
	}
	return UnknownNamespace
}
