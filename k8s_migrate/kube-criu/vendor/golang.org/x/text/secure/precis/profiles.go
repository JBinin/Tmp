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
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	// Implements the Nickname profile specified in RFC 7700.
	// The nickname profile is not idempotent and may need to be applied multiple
	// times before being used for comparisons.
	Nickname *Profile = nickname

	// Implements the UsernameCaseMapped profile specified in RFC 7613.
	UsernameCaseMapped *Profile = usernameCaseMap

	// Implements the UsernameCasePreserved profile specified in RFC 7613.
	UsernameCasePreserved *Profile = usernameNoCaseMap

	// Implements the OpaqueString profile defined in RFC 7613 for passwords and other secure labels.
	OpaqueString *Profile = opaquestring
)

var (
	nickname = &Profile{
		options: getOpts(
			AdditionalMapping(func() transform.Transformer {
				return &nickAdditionalMapping{}
			}),
			IgnoreCase,
			Norm(norm.NFKC),
			DisallowEmpty,
		),
		class: freeform,
	}
	usernameCaseMap = &Profile{
		options: getOpts(
			FoldWidth,
			LowerCase(),
			Norm(norm.NFC),
			BidiRule,
		),
		class: identifier,
	}
	usernameNoCaseMap = &Profile{
		options: getOpts(
			FoldWidth,
			Norm(norm.NFC),
			BidiRule,
		),
		class: identifier,
	}
	opaquestring = &Profile{
		options: getOpts(
			AdditionalMapping(func() transform.Transformer {
				return mapSpaces
			}),
			Norm(norm.NFC),
			DisallowEmpty,
		),
		class: freeform,
	}
)

// mapSpaces is a shared value of a runes.Map transformer.
var mapSpaces transform.Transformer = runes.Map(func(r rune) rune {
	if unicode.Is(unicode.Zs, r) {
		return ' '
	}
	return r
})
