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
Package po provides support for reading and writing GNU PO file.

Examples:
	import (
		"github.com/chai2010/gettext-go/gettext/po"
	)

	func main() {
		poFile, err := po.Load("test.po")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", poFile)
	}

The GNU PO file specification is at
http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html.
*/
package po
