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
// +build linux freebsd

package system

import (
	"strings"
	"testing"

	"github.com/docker/go-units"
)

// TestMemInfo tests parseMemInfo with a static meminfo string
func TestMemInfo(t *testing.T) {
	const input = `
	MemTotal:      1 kB
	MemFree:       2 kB
	SwapTotal:     3 kB
	SwapFree:      4 kB
	Malformed1:
	Malformed2:    1
	Malformed3:    2 MB
	Malformed4:    X kB
	`
	meminfo, err := parseMemInfo(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if meminfo.MemTotal != 1*units.KiB {
		t.Fatalf("Unexpected MemTotal: %d", meminfo.MemTotal)
	}
	if meminfo.MemFree != 2*units.KiB {
		t.Fatalf("Unexpected MemFree: %d", meminfo.MemFree)
	}
	if meminfo.SwapTotal != 3*units.KiB {
		t.Fatalf("Unexpected SwapTotal: %d", meminfo.SwapTotal)
	}
	if meminfo.SwapFree != 4*units.KiB {
		t.Fatalf("Unexpected SwapFree: %d", meminfo.SwapFree)
	}
}
