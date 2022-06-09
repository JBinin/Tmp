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
package nl

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

type IfAddrmsg struct {
	unix.IfAddrmsg
}

func NewIfAddrmsg(family int) *IfAddrmsg {
	return &IfAddrmsg{
		IfAddrmsg: unix.IfAddrmsg{
			Family: uint8(family),
		},
	}
}

// struct ifaddrmsg {
//   __u8    ifa_family;
//   __u8    ifa_prefixlen;  /* The prefix length    */
//   __u8    ifa_flags;  /* Flags      */
//   __u8    ifa_scope;  /* Address scope    */
//   __u32   ifa_index;  /* Link index     */
// };

// type IfAddrmsg struct {
// 	Family    uint8
// 	Prefixlen uint8
// 	Flags     uint8
// 	Scope     uint8
// 	Index     uint32
// }
// SizeofIfAddrmsg     = 0x8

func DeserializeIfAddrmsg(b []byte) *IfAddrmsg {
	return (*IfAddrmsg)(unsafe.Pointer(&b[0:unix.SizeofIfAddrmsg][0]))
}

func (msg *IfAddrmsg) Serialize() []byte {
	return (*(*[unix.SizeofIfAddrmsg]byte)(unsafe.Pointer(msg)))[:]
}

func (msg *IfAddrmsg) Len() int {
	return unix.SizeofIfAddrmsg
}

// struct ifa_cacheinfo {
// 	__u32	ifa_prefered;
// 	__u32	ifa_valid;
// 	__u32	cstamp; /* created timestamp, hundredths of seconds */
// 	__u32	tstamp; /* updated timestamp, hundredths of seconds */
// };

const IFA_CACHEINFO = 6
const SizeofIfaCacheInfo = 0x10

type IfaCacheInfo struct {
	IfaPrefered uint32
	IfaValid    uint32
	Cstamp      uint32
	Tstamp      uint32
}

func (msg *IfaCacheInfo) Len() int {
	return SizeofIfaCacheInfo
}

func DeserializeIfaCacheInfo(b []byte) *IfaCacheInfo {
	return (*IfaCacheInfo)(unsafe.Pointer(&b[0:SizeofIfaCacheInfo][0]))
}

func (msg *IfaCacheInfo) Serialize() []byte {
	return (*(*[SizeofIfaCacheInfo]byte)(unsafe.Pointer(msg)))[:]
}
