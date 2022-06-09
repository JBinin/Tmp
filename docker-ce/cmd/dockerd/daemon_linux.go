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

package main

import systemdDaemon "github.com/coreos/go-systemd/daemon"

// preNotifySystem sends a message to the host when the API is active, but before the daemon is
func preNotifySystem() {
}

// notifySystem sends a message to the host when the server is ready to be used
func notifySystem() {
	// Tell the init daemon we are accepting requests
	go systemdDaemon.SdNotify("READY=1")
}
