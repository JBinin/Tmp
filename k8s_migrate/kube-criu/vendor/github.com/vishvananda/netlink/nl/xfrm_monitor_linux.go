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
)

const (
	SizeofXfrmUserExpire = 0xe8
)

// struct xfrm_user_expire {
// 	struct xfrm_usersa_info		state;
// 	__u8				hard;
// };

type XfrmUserExpire struct {
	XfrmUsersaInfo XfrmUsersaInfo
	Hard           uint8
	Pad            [7]byte
}

func (msg *XfrmUserExpire) Len() int {
	return SizeofXfrmUserExpire
}

func DeserializeXfrmUserExpire(b []byte) *XfrmUserExpire {
	return (*XfrmUserExpire)(unsafe.Pointer(&b[0:SizeofXfrmUserExpire][0]))
}

func (msg *XfrmUserExpire) Serialize() []byte {
	return (*(*[SizeofXfrmUserExpire]byte)(unsafe.Pointer(msg)))[:]
}
