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
// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package language

// This file contains code common to the maketables.go and the package code.

// langAliasType is the type of an alias in langAliasMap.
type langAliasType int8

const (
	langDeprecated langAliasType = iota
	langMacro
	langLegacy

	langAliasTypeUnknown langAliasType = -1
)