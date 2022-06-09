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
// +build windows

package stats

// platformNewStatsCollector performs platform specific initialisation of the
// Collector structure. This is a no-op on Windows.
func platformNewStatsCollector(s *Collector) {
}

// getSystemCPUUsage returns the host system's cpu usage in
// nanoseconds. An error is returned if the format of the underlying
// file does not match. This is a no-op on Windows.
func (s *Collector) getSystemCPUUsage() (uint64, error) {
	return 0, nil
}
