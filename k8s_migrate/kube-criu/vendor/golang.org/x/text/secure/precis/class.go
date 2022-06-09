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
	"unicode/utf8"
)

// TODO: Add contextual character rules from Appendix A of RFC5892.

// A class is a set of characters that match certain derived properties. The
// PRECIS framework defines two classes: The Freeform class and the Identifier
// class. The freeform class should be used for profiles where expressiveness is
// prioritized over safety such as nicknames or passwords. The identifier class
// should be used for profiles where safety is the first priority such as
// addressable network labels and usernames.
type class struct {
	validFrom property
}

// Contains satisfies the runes.Set interface and returns whether the given rune
// is a member of the class.
func (c class) Contains(r rune) bool {
	b := make([]byte, 4)
	n := utf8.EncodeRune(b, r)

	trieval, _ := dpTrie.lookup(b[:n])
	return c.validFrom <= property(trieval)
}

var (
	identifier = &class{validFrom: pValid}
	freeform   = &class{validFrom: idDisOrFreePVal}
)
