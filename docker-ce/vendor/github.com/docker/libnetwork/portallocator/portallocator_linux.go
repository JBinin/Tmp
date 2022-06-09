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
package portallocator

import (
	"bufio"
	"fmt"
	"os"
)

func getDynamicPortRange() (start int, end int, err error) {
	const portRangeKernelParam = "/proc/sys/net/ipv4/ip_local_port_range"
	portRangeFallback := fmt.Sprintf("using fallback port range %d-%d", DefaultPortRangeStart, DefaultPortRangeEnd)
	file, err := os.Open(portRangeKernelParam)
	if err != nil {
		return 0, 0, fmt.Errorf("port allocator - %s due to error: %v", portRangeFallback, err)
	}

	defer file.Close()

	n, err := fmt.Fscanf(bufio.NewReader(file), "%d\t%d", &start, &end)
	if n != 2 || err != nil {
		if err == nil {
			err = fmt.Errorf("unexpected count of parsed numbers (%d)", n)
		}
		return 0, 0, fmt.Errorf("port allocator - failed to parse system ephemeral port range from %s - %s: %v", portRangeKernelParam, portRangeFallback, err)
	}
	return start, end, nil
}
