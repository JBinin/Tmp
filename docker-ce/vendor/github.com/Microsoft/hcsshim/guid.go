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
package hcsshim

import (
	"crypto/sha1"
	"fmt"
)

type GUID [16]byte

func NewGUID(source string) *GUID {
	h := sha1.Sum([]byte(source))
	var g GUID
	copy(g[0:], h[0:16])
	return &g
}

func (g *GUID) ToString() string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x-%02x", g[3], g[2], g[1], g[0], g[5], g[4], g[7], g[6], g[8:10], g[10:])
}
