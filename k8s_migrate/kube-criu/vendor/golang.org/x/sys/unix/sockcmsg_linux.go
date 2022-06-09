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
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Socket control messages

package unix

import "unsafe"

// UnixCredentials encodes credentials into a socket control message
// for sending to another process. This can be used for
// authentication.
func UnixCredentials(ucred *Ucred) []byte {
	b := make([]byte, CmsgSpace(SizeofUcred))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = SOL_SOCKET
	h.Type = SCM_CREDENTIALS
	h.SetLen(CmsgLen(SizeofUcred))
	*((*Ucred)(cmsgData(h))) = *ucred
	return b
}

// ParseUnixCredentials decodes a socket control message that contains
// credentials in a Ucred structure. To receive such a message, the
// SO_PASSCRED option must be enabled on the socket.
func ParseUnixCredentials(m *SocketControlMessage) (*Ucred, error) {
	if m.Header.Level != SOL_SOCKET {
		return nil, EINVAL
	}
	if m.Header.Type != SCM_CREDENTIALS {
		return nil, EINVAL
	}
	ucred := *(*Ucred)(unsafe.Pointer(&m.Data[0]))
	return &ucred, nil
}
