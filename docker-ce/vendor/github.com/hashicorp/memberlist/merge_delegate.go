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

// MergeDelegate is used to involve a client in
// a potential cluster merge operation. Namely, when
// a node does a TCP push/pull (as part of a join),
// the delegate is involved and allowed to cancel the join
// based on custom logic. The merge delegate is NOT invoked
// as part of the push-pull anti-entropy.
type MergeDelegate interface {
	// NotifyMerge is invoked when a merge could take place.
	// Provides a list of the nodes known by the peer. If
	// the return value is non-nil, the merge is canceled.
	NotifyMerge(peers []*Node) error
}
