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
package dbus

import (
	"errors"
	"os/exec"
)

func sessionBusPlatform() (*Conn, error) {
	cmd := exec.Command("launchctl", "getenv", "DBUS_LAUNCHD_SESSION_BUS_SOCKET")
	b, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, errors.New("dbus: couldn't determine address of session bus")
	}

	return Dial("unix:path=" + string(b[:len(b)-1]))
}
