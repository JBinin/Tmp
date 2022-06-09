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

package daemon

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-check/check"
)

func cleanupExecRoot(c *check.C, execRoot string) {
	// Cleanup network namespaces in the exec root of this
	// daemon because this exec root is specific to this
	// daemon instance and has no chance of getting
	// cleaned up when a new daemon is instantiated with a
	// new exec root.
	netnsPath := filepath.Join(execRoot, "netns")
	filepath.Walk(netnsPath, func(path string, info os.FileInfo, err error) error {
		if err := syscall.Unmount(path, syscall.MNT_FORCE); err != nil {
			c.Logf("unmount of %s failed: %v", path, err)
		}
		os.Remove(path)
		return nil
	})
}

// SignalDaemonDump sends a signal to the daemon to write a dump file
func SignalDaemonDump(pid int) {
	syscall.Kill(pid, syscall.SIGQUIT)
}

func signalDaemonReload(pid int) error {
	return syscall.Kill(pid, syscall.SIGHUP)
}
