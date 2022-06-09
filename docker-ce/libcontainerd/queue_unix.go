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
// +build linux solaris

package libcontainerd

import "sync"

type queue struct {
	sync.Mutex
	fns map[string]chan struct{}
}

func (q *queue) append(id string, f func()) {
	q.Lock()
	defer q.Unlock()

	if q.fns == nil {
		q.fns = make(map[string]chan struct{})
	}

	done := make(chan struct{})

	fn, ok := q.fns[id]
	q.fns[id] = done
	go func() {
		if ok {
			<-fn
		}
		f()
		close(done)
	}()
}
