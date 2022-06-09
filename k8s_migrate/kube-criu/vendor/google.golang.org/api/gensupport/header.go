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
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gensupport

import (
	"fmt"
	"runtime"
	"strings"
)

// GoogleClientHeader returns the value to use for the x-goog-api-client
// header, which is used internally by Google.
func GoogleClientHeader(generatorVersion, clientElement string) string {
	elts := []string{"gl-go/" + strings.Replace(runtime.Version(), " ", "_", -1)}
	if clientElement != "" {
		elts = append(elts, clientElement)
	}
	elts = append(elts, fmt.Sprintf("gdcl/%s", generatorVersion))
	return strings.Join(elts, " ")
}
