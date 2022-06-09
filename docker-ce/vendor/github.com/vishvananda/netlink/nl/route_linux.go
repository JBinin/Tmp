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
	"syscall"
	"unsafe"
)

type RtMsg struct {
	syscall.RtMsg
}

func NewRtMsg() *RtMsg {
	return &RtMsg{
		RtMsg: syscall.RtMsg{
			Table:    syscall.RT_TABLE_MAIN,
			Scope:    syscall.RT_SCOPE_UNIVERSE,
			Protocol: syscall.RTPROT_BOOT,
			Type:     syscall.RTN_UNICAST,
		},
	}
}

func NewRtDelMsg() *RtMsg {
	return &RtMsg{
		RtMsg: syscall.RtMsg{
			Table: syscall.RT_TABLE_MAIN,
			Scope: syscall.RT_SCOPE_NOWHERE,
		},
	}
}

func (msg *RtMsg) Len() int {
	return syscall.SizeofRtMsg
}

func DeserializeRtMsg(b []byte) *RtMsg {
	return (*RtMsg)(unsafe.Pointer(&b[0:syscall.SizeofRtMsg][0]))
}

func (msg *RtMsg) Serialize() []byte {
	return (*(*[syscall.SizeofRtMsg]byte)(unsafe.Pointer(msg)))[:]
}

type RtNexthop struct {
	syscall.RtNexthop
}

func DeserializeRtNexthop(b []byte) *RtNexthop {
	return (*RtNexthop)(unsafe.Pointer(&b[0:syscall.SizeofRtNexthop][0]))
}

func (msg *RtNexthop) Serialize() []byte {
	return (*(*[syscall.SizeofRtNexthop]byte)(unsafe.Pointer(msg)))[:]
}
