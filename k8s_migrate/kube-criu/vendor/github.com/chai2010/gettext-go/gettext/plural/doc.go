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

/*
Package plural provides standard plural formulas.

Examples:
	import (
		"code.google.com/p/gettext-go/gettext/plural"
	)

	func main() {
		enFormula := plural.Formula("en_US")
		xxFormula := plural.Formula("zh_CN")

		fmt.Printf("%s: %d\n", "en", enFormula(0))
		fmt.Printf("%s: %d\n", "en", enFormula(1))
		fmt.Printf("%s: %d\n", "en", enFormula(2))
		fmt.Printf("%s: %d\n", "??", xxFormula(0))
		fmt.Printf("%s: %d\n", "??", xxFormula(1))
		fmt.Printf("%s: %d\n", "??", xxFormula(2))
		fmt.Printf("%s: %d\n", "??", xxFormula(9))
		// Output:
		// en: 0
		// en: 0
		// en: 1
		// ??: 0
		// ??: 0
		// ??: 1
		// ??: 8
	}

See http://www.gnu.org/software/gettext/manual/html_node/Plural-forms.html
*/
package plural
