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

//go:generate go run gen.go

// Package internal contains non-exported functionality that are used by
// packages in the text repository.
package internal

import (
	"sort"

	"golang.org/x/text/language"
)

// SortTags sorts tags in place.
func SortTags(tags []language.Tag) {
	sort.Sort(sorter(tags))
}

type sorter []language.Tag

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sorter) Less(i, j int) bool {
	return s[i].String() < s[j].String()
}

// UniqueTags sorts and filters duplicate tags in place and returns a slice with
// only unique tags.
func UniqueTags(tags []language.Tag) []language.Tag {
	if len(tags) <= 1 {
		return tags
	}
	SortTags(tags)
	k := 0
	for i := 1; i < len(tags); i++ {
		if tags[k].String() < tags[i].String() {
			k++
			tags[k] = tags[i]
		}
	}
	return tags[:k+1]
}
