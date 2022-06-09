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
package dhcp4client

import (
	"net"
	"time"
)

type inetSock struct {
	*net.UDPConn

	laddr net.UDPAddr
	raddr net.UDPAddr
}

func NewInetSock(options ...func(*inetSock) error) (*inetSock, error) {
	c := &inetSock{
		laddr: net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 68},
		raddr: net.UDPAddr{IP: net.IPv4bcast, Port: 67},
	}

	err := c.setOption(options...)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", &c.laddr)
	if err != nil {
		return nil, err
	}

	c.UDPConn = conn
	return c, err
}

func (c *inetSock) setOption(options ...func(*inetSock) error) error {
	for _, opt := range options {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

func SetLocalAddr(l net.UDPAddr) func(*inetSock) error {
	return func(c *inetSock) error {
		c.laddr = l
		return nil
	}
}

func SetRemoteAddr(r net.UDPAddr) func(*inetSock) error {
	return func(c *inetSock) error {
		c.raddr = r
		return nil
	}
}

func (c *inetSock) Write(packet []byte) error {
	_, err := c.WriteToUDP(packet, &c.raddr)
	return err
}

func (c *inetSock) ReadFrom() ([]byte, net.IP, error) {
	readBuffer := make([]byte, MaxDHCPLen)
	n, source, err := c.ReadFromUDP(readBuffer)
	if source != nil {
		return readBuffer[:n], source.IP, err
	} else {
		return readBuffer[:n], net.IP{}, err
	}
}

func (c *inetSock) SetReadTimeout(t time.Duration) error {
	return c.SetReadDeadline(time.Now().Add(t))
}
