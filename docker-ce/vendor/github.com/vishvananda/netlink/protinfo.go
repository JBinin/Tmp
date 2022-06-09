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
package netlink

import (
	"strings"
)

// Protinfo represents bridge flags from netlink.
type Protinfo struct {
	Hairpin   bool
	Guard     bool
	FastLeave bool
	RootBlock bool
	Learning  bool
	Flood     bool
}

// String returns a list of enabled flags
func (prot *Protinfo) String() string {
	var boolStrings []string
	if prot.Hairpin {
		boolStrings = append(boolStrings, "Hairpin")
	}
	if prot.Guard {
		boolStrings = append(boolStrings, "Guard")
	}
	if prot.FastLeave {
		boolStrings = append(boolStrings, "FastLeave")
	}
	if prot.RootBlock {
		boolStrings = append(boolStrings, "RootBlock")
	}
	if prot.Learning {
		boolStrings = append(boolStrings, "Learning")
	}
	if prot.Flood {
		boolStrings = append(boolStrings, "Flood")
	}
	return strings.Join(boolStrings, " ")
}

func boolToByte(x bool) []byte {
	if x {
		return []byte{1}
	}
	return []byte{0}
}

func byteToBool(x byte) bool {
	return uint8(x) != 0
}
