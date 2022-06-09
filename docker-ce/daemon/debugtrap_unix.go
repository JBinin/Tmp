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
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	stackdump "github.com/docker/docker/pkg/signal"
)

func (d *Daemon) setupDumpStackTrap(root string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for range c {
			path, err := stackdump.DumpStacks(root)
			if err != nil {
				logrus.WithError(err).Error("failed to write goroutines dump")
			} else {
				logrus.Infof("goroutine stacks written to %s", path)
			}
			path, err = d.dumpDaemon(root)
			if err != nil {
				logrus.WithError(err).Error("failed to write daemon datastructure dump")
			} else {
				logrus.Infof("daemon datastructure dump written to %s", path)
			}
		}
	}()
}
