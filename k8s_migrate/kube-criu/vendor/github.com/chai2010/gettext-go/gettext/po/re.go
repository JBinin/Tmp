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

package po

import (
	"regexp"
)

var (
	reComment                = regexp.MustCompile(`^#`)              // #
	reExtractedComments      = regexp.MustCompile(`^#\.`)            // #.
	reReferenceComments      = regexp.MustCompile(`^#:`)             // #:
	reFlagsComments          = regexp.MustCompile(`^#,`)             // #, fuzzy,c-format
	rePrevMsgContextComments = regexp.MustCompile(`^#\|\s+msgctxt`)  // #| msgctxt
	rePrevMsgIdComments      = regexp.MustCompile(`^#\|\s+msgid`)    // #| msgid
	reStringLineComments     = regexp.MustCompile(`^#\|\s+".*"\s*$`) // #| "message"

	reMsgContext   = regexp.MustCompile(`^msgctxt\s+".*"\s*$`)            // msgctxt
	reMsgId        = regexp.MustCompile(`^msgid\s+".*"\s*$`)              // msgid
	reMsgIdPlural  = regexp.MustCompile(`^msgid_plural\s+".*"\s*$`)       // msgid_plural
	reMsgStr       = regexp.MustCompile(`^msgstr\s*".*"\s*$`)             // msgstr
	reMsgStrPlural = regexp.MustCompile(`^msgstr\s*(\[\d+\])\s*".*"\s*$`) // msgstr[0]
	reStringLine   = regexp.MustCompile(`^\s*".*"\s*$`)                   // "message"
	reBlankLine    = regexp.MustCompile(`^\s*$`)                          //
)

func (p *Message) isInvalidLine(s string) bool {
	if reComment.MatchString(s) {
		return false
	}
	if reBlankLine.MatchString(s) {
		return false
	}

	if reMsgContext.MatchString(s) {
		return false
	}
	if reMsgId.MatchString(s) {
		return false
	}
	if reMsgIdPlural.MatchString(s) {
		return false
	}
	if reMsgStr.MatchString(s) {
		return false
	}
	if reMsgStrPlural.MatchString(s) {
		return false
	}

	if reStringLine.MatchString(s) {
		return false
	}

	return true
}
