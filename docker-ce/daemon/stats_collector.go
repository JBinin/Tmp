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
	"runtime"
	"time"

	"github.com/docker/docker/daemon/stats"
	"github.com/docker/docker/pkg/system"
)

// newStatsCollector returns a new statsCollector that collections
// stats for a registered container at the specified interval.
// The collector allows non-running containers to be added
// and will start processing stats when they are started.
func (daemon *Daemon) newStatsCollector(interval time.Duration) *stats.Collector {
	// FIXME(vdemeester) move this elsewhere
	if runtime.GOOS == "linux" {
		meminfo, err := system.ReadMemInfo()
		if err == nil && meminfo.MemTotal > 0 {
			daemon.machineMemory = uint64(meminfo.MemTotal)
		}
	}
	s := stats.NewCollector(daemon, interval)
	go s.Run()
	return s
}
