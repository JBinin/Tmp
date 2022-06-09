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
package dbus

import (
	"encoding/binary"
	"errors"
	"io"
)

type genericTransport struct {
	io.ReadWriteCloser
}

func (t genericTransport) SendNullByte() error {
	_, err := t.Write([]byte{0})
	return err
}

func (t genericTransport) SupportsUnixFDs() bool {
	return false
}

func (t genericTransport) EnableUnixFDs() {}

func (t genericTransport) ReadMessage() (*Message, error) {
	return DecodeMessage(t)
}

func (t genericTransport) SendMessage(msg *Message) error {
	for _, v := range msg.Body {
		if _, ok := v.(UnixFD); ok {
			return errors.New("dbus: unix fd passing not enabled")
		}
	}
	return msg.EncodeTo(t, binary.LittleEndian)
}
