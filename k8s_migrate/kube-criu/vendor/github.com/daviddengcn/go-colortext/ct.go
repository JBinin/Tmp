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
/*
ct package provides functions to change the color of console text.

Under windows platform, the Console API is used. Under other systems, ANSI text mode is used.
*/
package ct

// Color is the type of color to be set.
type Color int

const (
	// No change of color
	None = Color(iota)
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// ResetColor resets the foreground and background to original colors
func ResetColor() {
	resetColor()
}

// ChangeColor sets the foreground and background colors. If the value of the color is None,
// the corresponding color keeps unchanged.
// If fgBright or bgBright is set true, corresponding color use bright color. bgBright may be
// ignored in some OS environment.
func ChangeColor(fg Color, fgBright bool, bg Color, bgBright bool) {
	changeColor(fg, fgBright, bg, bgBright)
}

// Foreground changes the foreground color.
func Foreground(cl Color, bright bool) {
	ChangeColor(cl, bright, None, false)
}

// Background changes the background color.
func Background(cl Color, bright bool) {
	ChangeColor(None, false, cl, bright)
}
