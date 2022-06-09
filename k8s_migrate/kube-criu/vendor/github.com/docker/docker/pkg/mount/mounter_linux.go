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
package mount

import (
	"golang.org/x/sys/unix"
)

const (
	// ptypes is the set propagation types.
	ptypes = unix.MS_SHARED | unix.MS_PRIVATE | unix.MS_SLAVE | unix.MS_UNBINDABLE

	// pflags is the full set valid flags for a change propagation call.
	pflags = ptypes | unix.MS_REC | unix.MS_SILENT

	// broflags is the combination of bind and read only
	broflags = unix.MS_BIND | unix.MS_RDONLY
)

// isremount returns true if either device name or flags identify a remount request, false otherwise.
func isremount(device string, flags uintptr) bool {
	switch {
	// We treat device "" and "none" as a remount request to provide compatibility with
	// requests that don't explicitly set MS_REMOUNT such as those manipulating bind mounts.
	case flags&unix.MS_REMOUNT != 0, device == "", device == "none":
		return true
	default:
		return false
	}
}

func mount(device, target, mType string, flags uintptr, data string) error {
	oflags := flags &^ ptypes
	if !isremount(device, flags) || data != "" {
		// Initial call applying all non-propagation flags for mount
		// or remount with changed data
		if err := unix.Mount(device, target, mType, oflags, data); err != nil {
			return err
		}
	}

	if flags&ptypes != 0 {
		// Change the propagation type.
		if err := unix.Mount("", target, "", flags&pflags, ""); err != nil {
			return err
		}
	}

	if oflags&broflags == broflags {
		// Remount the bind to apply read only.
		return unix.Mount("", target, "", oflags|unix.MS_REMOUNT, "")
	}

	return nil
}

func unmount(target string, flag int) error {
	return unix.Unmount(target, flag)
}
