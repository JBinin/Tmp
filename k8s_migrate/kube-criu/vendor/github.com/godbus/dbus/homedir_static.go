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
// +build static_build

package dbus

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func lookupHomeDir() string {
	myUid := os.Getuid()

	f, err := os.Open("/etc/passwd")
	if err != nil {
		return "/"
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		if err := s.Err(); err != nil {
			break
		}

		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")

		if len(parts) >= 6 {
			uid, err := strconv.Atoi(parts[2])
			if err == nil && uid == myUid {
				return parts[5]
			}
		}
	}

	// Default to / if we can't get a better value
	return "/"
}
