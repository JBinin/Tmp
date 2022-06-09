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

package gettext

import (
	"os"
	"strings"
)

func getDefaultLocale() string {
	if v := os.Getenv("LC_MESSAGES"); v != "" {
		return simplifiedLocale(v)
	}
	if v := os.Getenv("LANG"); v != "" {
		return simplifiedLocale(v)
	}
	return "default"
}

func simplifiedLocale(lang string) string {
	// en_US/en_US.UTF-8/zh_CN/zh_TW/el_GR@euro/...
	if idx := strings.Index(lang, ":"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "@"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}
	return strings.TrimSpace(lang)
}
