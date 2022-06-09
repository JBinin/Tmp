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
package memberlist

// ConflictDelegate is a used to inform a client that
// a node has attempted to join which would result in a
// name conflict. This happens if two clients are configured
// with the same name but different addresses.
type ConflictDelegate interface {
	// NotifyConflict is invoked when a name conflict is detected
	NotifyConflict(existing, other *Node)
}
