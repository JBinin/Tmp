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
package memberlist

import "time"

// PingDelegate is used to notify an observer how long it took for a ping message to
// complete a round trip.  It can also be used for writing arbitrary byte slices
// into ack messages. Note that in order to be meaningful for RTT estimates, this
// delegate does not apply to indirect pings, nor fallback pings sent over TCP.
type PingDelegate interface {
	// AckPayload is invoked when an ack is being sent; the returned bytes will be appended to the ack
	AckPayload() []byte
	// NotifyPing is invoked when an ack for a ping is received
	NotifyPingComplete(other *Node, rtt time.Duration, payload []byte)
}
