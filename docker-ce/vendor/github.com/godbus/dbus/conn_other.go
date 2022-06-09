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
// +build !darwin

package dbus

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func sessionBusPlatform() (*Conn, error) {
	cmd := exec.Command("dbus-launch")
	b, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	i := bytes.IndexByte(b, '=')
	j := bytes.IndexByte(b, '\n')

	if i == -1 || j == -1 {
		return nil, errors.New("dbus: couldn't determine address of session bus")
	}

	env, addr := string(b[0:i]), string(b[i+1:j])
	os.Setenv(env, addr)

	return Dial(addr)
}
