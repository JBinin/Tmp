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
// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mo

import (
	"bytes"
	"fmt"
)

// A MO file is made up of many entries,
// each entry holding the relation between an original untranslated string
// and its corresponding translation.
//
// See http://www.gnu.org/software/gettext/manual/html_node/MO-Files.html
type Message struct {
	MsgContext   string   // msgctxt context
	MsgId        string   // msgid untranslated-string
	MsgIdPlural  string   // msgid_plural untranslated-string-plural
	MsgStr       string   // msgstr translated-string
	MsgStrPlural []string // msgstr[0] translated-string-case-0
}

// String returns the po format entry string.
func (p Message) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "msgid %s", encodePoString(p.MsgId))
	if p.MsgIdPlural != "" {
		fmt.Fprintf(&buf, "msgid_plural %s", encodePoString(p.MsgIdPlural))
	}
	if p.MsgStr != "" {
		fmt.Fprintf(&buf, "msgstr %s", encodePoString(p.MsgStr))
	}
	for i := 0; i < len(p.MsgStrPlural); i++ {
		fmt.Fprintf(&buf, "msgstr[%d] %s", i, encodePoString(p.MsgStrPlural[i]))
	}
	return buf.String()
}
