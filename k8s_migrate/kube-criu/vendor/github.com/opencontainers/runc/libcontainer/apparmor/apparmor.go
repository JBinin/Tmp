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
// +build apparmor,linux

package apparmor

import (
	"fmt"
	"io/ioutil"
	"os"
)

// IsEnabled returns true if apparmor is enabled for the host.
func IsEnabled() bool {
	if _, err := os.Stat("/sys/kernel/security/apparmor"); err == nil && os.Getenv("container") == "" {
		if _, err = os.Stat("/sbin/apparmor_parser"); err == nil {
			buf, err := ioutil.ReadFile("/sys/module/apparmor/parameters/enabled")
			return err == nil && len(buf) > 1 && buf[0] == 'Y'
		}
	}
	return false
}

func setprocattr(attr, value string) error {
	// Under AppArmor you can only change your own attr, so use /proc/self/
	// instead of /proc/<tid>/ like libapparmor does
	path := fmt.Sprintf("/proc/self/attr/%s", attr)

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s", value)
	return err
}

// changeOnExec reimplements aa_change_onexec from libapparmor in Go
func changeOnExec(name string) error {
	value := "exec " + name
	if err := setprocattr("exec", value); err != nil {
		return fmt.Errorf("apparmor failed to apply profile: %s", err)
	}
	return nil
}

// ApplyProfile will apply the profile with the specified name to the process after
// the next exec.
func ApplyProfile(name string) error {
	if name == "" {
		return nil
	}

	return changeOnExec(name)
}
