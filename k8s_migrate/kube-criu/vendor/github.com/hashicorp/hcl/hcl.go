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
// Package hcl decodes HCL into usable Go structures.
//
// hcl input can come in either pure HCL format or JSON format.
// It can be parsed into an AST, and then decoded into a structure,
// or it can be decoded directly from a string into a structure.
//
// If you choose to parse HCL into a raw AST, the benefit is that you
// can write custom visitor implementations to implement custom
// semantic checks. By default, HCL does not perform any semantic
// checks.
package hcl
