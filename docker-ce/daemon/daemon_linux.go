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
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/mount"
)

// On Linux, plugins use a static path for storing execution state,
// instead of deriving path from daemon's exec-root. This is because
// plugin socket files are created here and they cannot exceed max
// path length of 108 bytes.
func getPluginExecRoot(root string) string {
	return "/run/docker/plugins"
}

func (daemon *Daemon) cleanupMountsByID(id string) error {
	logrus.Debugf("Cleaning up old mountid %s: start.", id)
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return err
	}
	defer f.Close()

	return daemon.cleanupMountsFromReaderByID(f, id, mount.Unmount)
}

func (daemon *Daemon) cleanupMountsFromReaderByID(reader io.Reader, id string, unmount func(target string) error) error {
	if daemon.root == "" {
		return nil
	}
	var errors []string

	regexps := getCleanPatterns(id)
	sc := bufio.NewScanner(reader)
	for sc.Scan() {
		if fields := strings.Fields(sc.Text()); len(fields) >= 4 {
			if mnt := fields[4]; strings.HasPrefix(mnt, daemon.root) {
				for _, p := range regexps {
					if p.MatchString(mnt) {
						if err := unmount(mnt); err != nil {
							logrus.Error(err)
							errors = append(errors, err.Error())
						}
					}
				}
			}
		}
	}

	if err := sc.Err(); err != nil {
		return err
	}

	if len(errors) > 0 {
		return fmt.Errorf("Error cleaning up mounts:\n%v", strings.Join(errors, "\n"))
	}

	logrus.Debugf("Cleaning up old mountid %v: done.", id)
	return nil
}

// cleanupMounts umounts shm/mqueue mounts for old containers
func (daemon *Daemon) cleanupMounts() error {
	return daemon.cleanupMountsByID("")
}

func getCleanPatterns(id string) (regexps []*regexp.Regexp) {
	var patterns []string
	if id == "" {
		id = "[0-9a-f]{64}"
		patterns = append(patterns, "containers/"+id+"/shm")
	}
	patterns = append(patterns, "aufs/mnt/"+id+"$", "overlay/"+id+"/merged$", "zfs/graph/"+id+"$")
	for _, p := range patterns {
		r, err := regexp.Compile(p)
		if err == nil {
			regexps = append(regexps, r)
		}
	}
	return
}
