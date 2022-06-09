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
	"fmt"
)

// Framer is a type of function that takes an input string (typically an
// already-formatted syslog message) and applies "message framing" to it. We
// have different framers because different versions of the syslog protocol
// and its transport requirements define different framing behavior.
type Framer func(in string) string

// DefaultFramer does nothing, since there is no framing to apply. This is
// the original behavior of the Go syslog package, and is also typically used
// for UDP syslog.
func DefaultFramer(in string) string {
	return in
}

// RFC5425MessageLengthFramer prepends the message length to the front of the
// provided message, as defined in RFC 5425.
func RFC5425MessageLengthFramer(in string) string {
	return fmt.Sprintf("%d %s", len(in), in)
}
