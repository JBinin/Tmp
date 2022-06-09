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
//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package utils

type StringStack struct {
	list []string
}

func NewStringStack() *StringStack {
	a := &StringStack{}
	a.list = make([]string, 0)
	return a
}

func (a *StringStack) IsEmpty() bool {
	return len(a.list) == 0
}

func (a *StringStack) Pop() (x string) {
	x, a.list = a.list[0], a.list[1:len(a.list)]
	return
}

func (a *StringStack) Push(x string) {
	a.list = append(a.list, x)
}
