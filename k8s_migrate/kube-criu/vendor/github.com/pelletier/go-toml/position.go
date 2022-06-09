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
// Position support for go-toml

package toml

import (
	"fmt"
)

// Position of a document element within a TOML document.
//
// Line and Col are both 1-indexed positions for the element's line number and
// column number, respectively.  Values of zero or less will cause Invalid(),
// to return true.
type Position struct {
	Line int // line within the document
	Col  int // column within the line
}

// String representation of the position.
// Displays 1-indexed line and column numbers.
func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.Line, p.Col)
}

// Invalid returns whether or not the position is valid (i.e. with negative or
// null values)
func (p Position) Invalid() bool {
	return p.Line <= 0 || p.Col <= 0
}
