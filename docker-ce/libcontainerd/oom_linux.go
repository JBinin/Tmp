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
package libcontainerd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/opencontainers/runc/libcontainer/system"
)

func setOOMScore(pid, score int) error {
	oomScoreAdjPath := fmt.Sprintf("/proc/%d/oom_score_adj", pid)
	f, err := os.OpenFile(oomScoreAdjPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	stringScore := strconv.Itoa(score)
	_, err = f.WriteString(stringScore)
	f.Close()
	if os.IsPermission(err) {
		// Setting oom_score_adj does not work in an
		// unprivileged container. Ignore the error, but log
		// it if we appear not to be in that situation.
		if !system.RunningInUserNS() {
			logrus.Debugf("Permission denied writing %q to %s", stringScore, oomScoreAdjPath)
		}
		return nil
	}
	return err
}
