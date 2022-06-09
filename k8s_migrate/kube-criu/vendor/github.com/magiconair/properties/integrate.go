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
// Copyright 2016 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import "flag"

// MustFlag sets flags that are skipped by dst.Parse when p contains
// the respective key for flag.Flag.Name.
//
// It's use is recommended with command line arguments as in:
// 	flag.Parse()
// 	p.MustFlag(flag.CommandLine)
func (p *Properties) MustFlag(dst *flag.FlagSet) {
	m := make(map[string]*flag.Flag)
	dst.VisitAll(func(f *flag.Flag) {
		m[f.Name] = f
	})
	dst.Visit(func(f *flag.Flag) {
		delete(m, f.Name) // overridden
	})

	for name, f := range m {
		v, ok := p.Get(name)
		if !ok {
			continue
		}

		if err := f.Value.Set(v); err != nil {
			ErrorHandler(err)
		}
	}
}
