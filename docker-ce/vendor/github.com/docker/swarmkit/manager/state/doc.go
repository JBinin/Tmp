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
// Package state provides interfaces to work with swarm cluster state.
//
// The primary interface is Store, which abstracts storage of this cluster
// state. Store exposes a transactional interface for both reads and writes.
// To perform a read transaction, View accepts a callback function that it
// will invoke with a ReadTx object that gives it a consistent view of the
// state. Similarly, Update accepts a callback function that it will invoke with
// a Tx object that allows reads and writes to happen without interference from
// other transactions.
//
// This is an example of making an update to a Store:
//
//	err := store.Update(func(tx state.Tx) {
//		if err := tx.Nodes().Update(newNode); err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return fmt.Errorf("transaction failed: %v", err)
//	}
//
// WatchableStore is a version of Store that exposes watch functionality.
// These expose a publish/subscribe queue where code can subscribe to
// changes of interest. This can be combined with the ViewAndWatch function to
// "fork" a store, by making a snapshot and then applying future changes
// to keep the copy in sync. This approach lets consumers of the data
// use their own data structures and implement their own concurrency
// strategies. It can lead to more efficient code because data consumers
// don't necessarily have to lock the main data store if they are
// maintaining their own copies of the state.
package state
