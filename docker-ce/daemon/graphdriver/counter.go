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
package graphdriver

import "sync"

type minfo struct {
	check bool
	count int
}

// RefCounter is a generic counter for use by graphdriver Get/Put calls
type RefCounter struct {
	counts  map[string]*minfo
	mu      sync.Mutex
	checker Checker
}

// NewRefCounter returns a new RefCounter
func NewRefCounter(c Checker) *RefCounter {
	return &RefCounter{
		checker: c,
		counts:  make(map[string]*minfo),
	}
}

// Increment increaes the ref count for the given id and returns the current count
func (c *RefCounter) Increment(path string) int {
	c.mu.Lock()
	m := c.counts[path]
	if m == nil {
		m = &minfo{}
		c.counts[path] = m
	}
	// if we are checking this path for the first time check to make sure
	// if it was already mounted on the system and make sure we have a correct ref
	// count if it is mounted as it is in use.
	if !m.check {
		m.check = true
		if c.checker.IsMounted(path) {
			m.count++
		}
	}
	m.count++
	count := m.count
	c.mu.Unlock()
	return count
}

// Decrement decreases the ref count for the given id and returns the current count
func (c *RefCounter) Decrement(path string) int {
	c.mu.Lock()
	m := c.counts[path]
	if m == nil {
		m = &minfo{}
		c.counts[path] = m
	}
	// if we are checking this path for the first time check to make sure
	// if it was already mounted on the system and make sure we have a correct ref
	// count if it is mounted as it is in use.
	if !m.check {
		m.check = true
		if c.checker.IsMounted(path) {
			m.count++
		}
	}
	m.count--
	count := m.count
	c.mu.Unlock()
	return count
}
