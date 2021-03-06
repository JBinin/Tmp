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
	"os"
	"time"
)

// Formatter is a type of function that takes the consituent parts of a
// syslog message and returns a formatted string. A different Formatter is
// defined for each different syslog protocol we support.
type Formatter func(p Priority, hostname, tag, content string) string

// DefaultFormatter is the original format supported by the Go syslog package,
// and is a non-compliant amalgamation of 3164 and 5424 that is intended to
// maximize compatibility.
func DefaultFormatter(p Priority, hostname, tag, content string) string {
	timestamp := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf("<%d> %s %s %s[%d]: %s",
		p, timestamp, hostname, tag, os.Getpid(), content)
	return msg
}

// UnixFormatter omits the hostname, because it is only used locally.
func UnixFormatter(p Priority, hostname, tag, content string) string {
	timestamp := time.Now().Format(time.Stamp)
	msg := fmt.Sprintf("<%d>%s %s[%d]: %s",
		p, timestamp, tag, os.Getpid(), content)
	return msg
}

// RFC3164Formatter provides an RFC 3164 compliant message.
func RFC3164Formatter(p Priority, hostname, tag, content string) string {
	timestamp := time.Now().Format(time.Stamp)
	msg := fmt.Sprintf("<%d>%s %s %s[%d]: %s",
		p, timestamp, hostname, tag, os.Getpid(), content)
	return msg
}

// RFC5424Formatter provides an RFC 5424 compliant message.
func RFC5424Formatter(p Priority, hostname, tag, content string) string {
	timestamp := time.Now().Format(time.RFC3339)
	pid := os.Getpid()
	appName := os.Args[0]
	msg := fmt.Sprintf("<%d>%d %s %s %s %d %s %s",
		p, 1, timestamp, hostname, appName, pid, tag, content)
	return msg
}
