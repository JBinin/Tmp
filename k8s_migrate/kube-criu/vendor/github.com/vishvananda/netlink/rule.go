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
	"fmt"
	"net"
)

// Rule represents a netlink rule.
type Rule struct {
	Priority          int
	Family            int
	Table             int
	Mark              int
	Mask              int
	TunID             uint
	Goto              int
	Src               *net.IPNet
	Dst               *net.IPNet
	Flow              int
	IifName           string
	OifName           string
	SuppressIfgroup   int
	SuppressPrefixlen int
	Invert            bool
}

func (r Rule) String() string {
	return fmt.Sprintf("ip rule %d: from %s table %d", r.Priority, r.Src, r.Table)
}

// NewRule return empty rules.
func NewRule() *Rule {
	return &Rule{
		SuppressIfgroup:   -1,
		SuppressPrefixlen: -1,
		Priority:          -1,
		Mark:              -1,
		Mask:              -1,
		Goto:              -1,
		Flow:              -1,
	}
}
