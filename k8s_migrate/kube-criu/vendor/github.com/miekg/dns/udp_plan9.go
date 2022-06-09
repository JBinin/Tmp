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
package dns

import (
	"net"
)

func setUDPSocketOptions(conn *net.UDPConn) error { return nil }

// SessionUDP holds the remote address and the associated
// out-of-band data.
type SessionUDP struct {
	raddr   *net.UDPAddr
	context []byte
}

// RemoteAddr returns the remote network address.
func (s *SessionUDP) RemoteAddr() net.Addr { return s.raddr }

// ReadFromSessionUDP acts just like net.UDPConn.ReadFrom(), but returns a session object instead of a
// net.UDPAddr.
func ReadFromSessionUDP(conn *net.UDPConn, b []byte) (int, *SessionUDP, error) {
	oob := make([]byte, 40)
	n, oobn, _, raddr, err := conn.ReadMsgUDP(b, oob)
	if err != nil {
		return n, nil, err
	}
	return n, &SessionUDP{raddr, oob[:oobn]}, err
}

// WriteToSessionUDP acts just like net.UDPConn.WritetTo(), but uses a *SessionUDP instead of a net.Addr.
func WriteToSessionUDP(conn *net.UDPConn, b []byte, session *SessionUDP) (int, error) {
	n, _, err := conn.WriteMsgUDP(b, session.context, session.raddr)
	return n, err
}
