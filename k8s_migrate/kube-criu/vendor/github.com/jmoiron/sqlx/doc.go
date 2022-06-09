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
// Package sqlx provides general purpose extensions to database/sql.
//
// It is intended to seamlessly wrap database/sql and provide convenience
// methods which are useful in the development of database driven applications.
// None of the underlying database/sql methods are changed.  Instead all extended
// behavior is implemented through new methods defined on wrapper types.
//
// Additions include scanning into structs, named query support, rebinding
// queries for different drivers, convenient shorthands for common error handling
// and more.
//
package sqlx
