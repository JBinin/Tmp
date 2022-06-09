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
	"encoding/binary"

	"github.com/vishvananda/netlink/nl"
)

var (
	native       = nl.NativeEndian()
	networkOrder = binary.BigEndian
)

func htonl(val uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, val)
	return bytes
}

func htons(val uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, val)
	return bytes
}

func ntohl(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func ntohs(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}
