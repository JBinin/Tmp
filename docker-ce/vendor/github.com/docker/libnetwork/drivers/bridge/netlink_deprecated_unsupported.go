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
// +build !linux

package bridge

import (
	"errors"
	"net"
)

// Add a slave to a bridge device.  This is more backward-compatible than
// netlink.NetworkSetMaster and works on RHEL 6.
func ioctlAddToBridge(iface, master *net.Interface) error {
	return errors.New("not implemented")
}

func ioctlCreateBridge(name string, setMacAddr bool) error {
	return errors.New("not implemented")
}
