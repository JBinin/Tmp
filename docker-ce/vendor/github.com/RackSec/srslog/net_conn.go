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
package srslog

import (
	"net"
)

// netConn has an internal net.Conn and adheres to the serverConn interface,
// allowing us to send syslog messages over the network.
type netConn struct {
	conn net.Conn
}

// writeString formats syslog messages using time.RFC3339 and includes the
// hostname, and sends the message to the connection.
func (n *netConn) writeString(framer Framer, formatter Formatter, p Priority, hostname, tag, msg string) error {
	if framer == nil {
		framer = DefaultFramer
	}
	if formatter == nil {
		formatter = DefaultFormatter
	}
	formattedMessage := framer(formatter(p, hostname, tag, msg))
	_, err := n.conn.Write([]byte(formattedMessage))
	return err
}

// close the network connection
func (n *netConn) close() error {
	return n.conn.Close()
}
