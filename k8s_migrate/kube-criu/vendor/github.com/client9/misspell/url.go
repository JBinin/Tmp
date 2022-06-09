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
package misspell

import (
	"regexp"
)

// Regexp for URL https://mathiasbynens.be/demo/url-regex
//
// original @imme_emosol (54 chars) has trouble with dashes in hostname
// @(https?|ftp)://(-\.)?([^\s/?\.#-]+\.?)+(/[^\s]*)?$@iS
var reURL = regexp.MustCompile(`(?i)(https?|ftp)://(-\.)?([^\s/?\.#]+\.?)+(/[^\s]*)?`)

// StripURL attemps to replace URLs with blank spaces, e.g.
//  "xxx http://foo.com/ yyy -> "xxx          yyyy"
func StripURL(s string) string {
	return reURL.ReplaceAllStringFunc(s, replaceWithBlanks)
}
