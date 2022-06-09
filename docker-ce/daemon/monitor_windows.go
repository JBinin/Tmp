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
package daemon

import (
	"fmt"

	"github.com/docker/docker/container"
	"github.com/docker/docker/libcontainerd"
)

// platformConstructExitStatus returns a platform specific exit status structure
func platformConstructExitStatus(e libcontainerd.StateInfo) *container.ExitStatus {
	return &container.ExitStatus{
		ExitCode: int(e.ExitCode),
	}
}

// postRunProcessing perfoms any processing needed on the container after it has stopped.
func (daemon *Daemon) postRunProcessing(container *container.Container, e libcontainerd.StateInfo) error {
	if e.ExitCode == 0 && e.UpdatePending {
		spec, err := daemon.createSpec(container)
		if err != nil {
			return err
		}

		newOpts := []libcontainerd.CreateOption{&libcontainerd.ServicingOption{
			IsServicing: true,
		}}

		copts, err := daemon.getLibcontainerdCreateOptions(container)
		if err != nil {
			return err
		}

		if copts != nil {
			newOpts = append(newOpts, copts...)
		}

		// Create a new servicing container, which will start, complete the update, and merge back the
		// results if it succeeded, all as part of the below function call.
		if err := daemon.containerd.Create((container.ID + "_servicing"), "", "", *spec, container.InitializeStdio, newOpts...); err != nil {
			container.SetExitCode(-1)
			return fmt.Errorf("Post-run update servicing failed: %s", err)
		}
	}
	return nil
}
