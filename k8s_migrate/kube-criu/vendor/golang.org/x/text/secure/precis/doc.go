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

// Package precis contains types and functions for the preparation,
// enforcement, and comparison of internationalized strings ("PRECIS") as
// defined in RFC 7564. It also contains several pre-defined profiles for
// passwords, nicknames, and usernames as defined in RFC 7613 and RFC 7700.
//
// BE ADVISED: This package is under construction and the API may change in
// backwards incompatible ways and without notice.
package precis

//go:generate go run gen.go gen_trieval.go
