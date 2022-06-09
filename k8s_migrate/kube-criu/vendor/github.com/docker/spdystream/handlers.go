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
package spdystream

import (
	"io"
	"net/http"
)

// MirrorStreamHandler mirrors all streams.
func MirrorStreamHandler(stream *Stream) {
	replyErr := stream.SendReply(http.Header{}, false)
	if replyErr != nil {
		return
	}

	go func() {
		io.Copy(stream, stream)
		stream.Close()
	}()
	go func() {
		for {
			header, receiveErr := stream.ReceiveHeader()
			if receiveErr != nil {
				return
			}
			sendErr := stream.SendHeader(header, false)
			if sendErr != nil {
				return
			}
		}
	}()
}

// NoopStreamHandler does nothing when stream connects, most
// likely used with RejectAuthHandler which will not allow any
// streams to make it to the stream handler.
func NoOpStreamHandler(stream *Stream) {
	stream.SendReply(http.Header{}, false)
}
