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
	"fmt"

	aaprofile "github.com/docker/docker/profiles/apparmor"
	"github.com/opencontainers/runc/libcontainer/apparmor"
)

// Define constants for native driver
const (
	defaultApparmorProfile = "docker-default"
)

func ensureDefaultAppArmorProfile() error {
	if apparmor.IsEnabled() {
		loaded, err := aaprofile.IsLoaded(defaultApparmorProfile)
		if err != nil {
			return fmt.Errorf("Could not check if %s AppArmor profile was loaded: %s", defaultApparmorProfile, err)
		}

		// Nothing to do.
		if loaded {
			return nil
		}

		// Load the profile.
		if err := aaprofile.InstallDefault(defaultApparmorProfile); err != nil {
			return fmt.Errorf("AppArmor enabled on system but the %s profile could not be loaded.", defaultApparmorProfile)
		}
	}

	return nil
}
