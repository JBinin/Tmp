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

// +build ignore

package main

// Class is the Unicode BiDi class. Each rune has a single class.
type Class uint

const (
	L       Class = iota // LeftToRight
	R                    // RightToLeft
	EN                   // EuropeanNumber
	ES                   // EuropeanSeparator
	ET                   // EuropeanTerminator
	AN                   // ArabicNumber
	CS                   // CommonSeparator
	B                    // ParagraphSeparator
	S                    // SegmentSeparator
	WS                   // WhiteSpace
	ON                   // OtherNeutral
	BN                   // BoundaryNeutral
	NSM                  // NonspacingMark
	AL                   // ArabicLetter
	Control              // Control LRO - PDI

	numClass

	LRO // LeftToRightOverride
	RLO // RightToLeftOverride
	LRE // LeftToRightEmbedding
	RLE // RightToLeftEmbedding
	PDF // PopDirectionalFormat
	LRI // LeftToRightIsolate
	RLI // RightToLeftIsolate
	FSI // FirstStrongIsolate
	PDI // PopDirectionalIsolate

	unknownClass = ^Class(0)
)

var controlToClass = map[rune]Class{
	0x202D: LRO, // LeftToRightOverride,
	0x202E: RLO, // RightToLeftOverride,
	0x202A: LRE, // LeftToRightEmbedding,
	0x202B: RLE, // RightToLeftEmbedding,
	0x202C: PDF, // PopDirectionalFormat,
	0x2066: LRI, // LeftToRightIsolate,
	0x2067: RLI, // RightToLeftIsolate,
	0x2068: FSI, // FirstStrongIsolate,
	0x2069: PDI, // PopDirectionalIsolate,
}

// A trie entry has the following bits:
// 7..5  XOR mask for brackets
// 4     1: Bracket open, 0: Bracket close
// 3..0  Class type

const (
	openMask     = 0x10
	xorMaskShift = 5
)
