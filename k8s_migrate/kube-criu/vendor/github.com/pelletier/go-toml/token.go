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
package toml

import (
	"fmt"
	"strconv"
	"unicode"
)

// Define tokens
type tokenType int

const (
	eof = -(iota + 1)
)

const (
	tokenError tokenType = iota
	tokenEOF
	tokenComment
	tokenKey
	tokenString
	tokenInteger
	tokenTrue
	tokenFalse
	tokenFloat
	tokenInf
	tokenNan
	tokenEqual
	tokenLeftBracket
	tokenRightBracket
	tokenLeftCurlyBrace
	tokenRightCurlyBrace
	tokenLeftParen
	tokenRightParen
	tokenDoubleLeftBracket
	tokenDoubleRightBracket
	tokenDate
	tokenKeyGroup
	tokenKeyGroupArray
	tokenComma
	tokenColon
	tokenDollar
	tokenStar
	tokenQuestion
	tokenDot
	tokenDotDot
	tokenEOL
)

var tokenTypeNames = []string{
	"Error",
	"EOF",
	"Comment",
	"Key",
	"String",
	"Integer",
	"True",
	"False",
	"Float",
	"Inf",
	"NaN",
	"=",
	"[",
	"]",
	"{",
	"}",
	"(",
	")",
	"]]",
	"[[",
	"Date",
	"KeyGroup",
	"KeyGroupArray",
	",",
	":",
	"$",
	"*",
	"?",
	".",
	"..",
	"EOL",
}

type token struct {
	Position
	typ tokenType
	val string
}

func (tt tokenType) String() string {
	idx := int(tt)
	if idx < len(tokenTypeNames) {
		return tokenTypeNames[idx]
	}
	return "Unknown"
}

func (t token) Int() int {
	if result, err := strconv.Atoi(t.val); err != nil {
		panic(err)
	} else {
		return result
	}
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}

	return fmt.Sprintf("%q", t.val)
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isKeyChar(r rune) bool {
	// Keys start with the first character that isn't whitespace or [ and end
	// with the last non-whitespace character before the equals sign. Keys
	// cannot contain a # character."
	return !(r == '\r' || r == '\n' || r == eof || r == '=')
}

func isKeyStartChar(r rune) bool {
	return !(isSpace(r) || r == '\r' || r == '\n' || r == eof || r == '[')
}

func isDigit(r rune) bool {
	return unicode.IsNumber(r)
}

func isHexDigit(r rune) bool {
	return isDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
