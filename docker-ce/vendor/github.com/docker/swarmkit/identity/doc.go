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
// Package identity provides functionality for generating and manager
// identifiers within swarm. This includes entity identification, such as that
// of Service, Task and Network but also cryptographically-secure Node identity.
//
// Random Identifiers
//
// Identifiers provided by this package are cryptographically-strong, random
// 128 bit numbers encoded in Base36. This method is preferred over UUID4 since
// it requires less storage and leverages the full 128 bits of entropy.
//
// Generating an identifier is simple. Simply call the `NewID` function, check
// the error and proceed:
//
// 	id, err := NewID()
// 	if err != nil { /* ... handle it, please ... */ }
//
package identity
