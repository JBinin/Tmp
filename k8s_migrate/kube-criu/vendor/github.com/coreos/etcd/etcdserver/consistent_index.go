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
// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdserver

import (
	"sync/atomic"
)

// consistentIndex represents the offset of an entry in a consistent replica log.
// It implements the mvcc.ConsistentIndexGetter interface.
// It is always set to the offset of current entry before executing the entry,
// so ConsistentWatchableKV could get the consistent index from it.
type consistentIndex uint64

func (i *consistentIndex) setConsistentIndex(v uint64) {
	atomic.StoreUint64((*uint64)(i), v)
}

func (i *consistentIndex) ConsistentIndex() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}
