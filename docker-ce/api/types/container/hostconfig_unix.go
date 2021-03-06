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
// +build !windows

package container

import "strings"

// IsValid indicates if an isolation technology is valid
func (i Isolation) IsValid() bool {
	return i.IsDefault()
}

// IsPrivate indicates whether container uses its private network stack.
func (n NetworkMode) IsPrivate() bool {
	return !(n.IsHost() || n.IsContainer())
}

// IsDefault indicates whether container uses the default network stack.
func (n NetworkMode) IsDefault() bool {
	return n == "default"
}

// NetworkName returns the name of the network stack.
func (n NetworkMode) NetworkName() string {
	if n.IsBridge() {
		return "bridge"
	} else if n.IsHost() {
		return "host"
	} else if n.IsContainer() {
		return "container"
	} else if n.IsNone() {
		return "none"
	} else if n.IsDefault() {
		return "default"
	} else if n.IsUserDefined() {
		return n.UserDefined()
	}
	return ""
}

// IsBridge indicates whether container uses the bridge network stack
func (n NetworkMode) IsBridge() bool {
	return n == "bridge"
}

// IsHost indicates whether container uses the host network stack.
func (n NetworkMode) IsHost() bool {
	return n == "host"
}

// IsContainer indicates whether container uses a container network stack.
func (n NetworkMode) IsContainer() bool {
	parts := strings.SplitN(string(n), ":", 2)
	return len(parts) > 1 && parts[0] == "container"
}

// IsNone indicates whether container isn't using a network stack.
func (n NetworkMode) IsNone() bool {
	return n == "none"
}

// ConnectedContainer is the id of the container which network this container is connected to.
func (n NetworkMode) ConnectedContainer() string {
	parts := strings.SplitN(string(n), ":", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// IsUserDefined indicates user-created network
func (n NetworkMode) IsUserDefined() bool {
	return !n.IsDefault() && !n.IsBridge() && !n.IsHost() && !n.IsNone() && !n.IsContainer()
}

//UserDefined indicates user-created network
func (n NetworkMode) UserDefined() string {
	if n.IsUserDefined() {
		return string(n)
	}
	return ""
}
