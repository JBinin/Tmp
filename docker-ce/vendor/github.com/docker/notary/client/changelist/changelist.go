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
package changelist

// memChangeList implements a simple in memory change list.
type memChangelist struct {
	changes []Change
}

// NewMemChangelist instantiates a new in-memory changelist
func NewMemChangelist() Changelist {
	return &memChangelist{}
}

// List returns a list of Changes
func (cl memChangelist) List() []Change {
	return cl.changes
}

// Add adds a change to the in-memory change list
func (cl *memChangelist) Add(c Change) error {
	cl.changes = append(cl.changes, c)
	return nil
}

// Remove deletes the changes found at the given indices
func (cl *memChangelist) Remove(idxs []int) error {
	remove := make(map[int]struct{})
	for _, i := range idxs {
		remove[i] = struct{}{}
	}
	var keep []Change

	for i, c := range cl.changes {
		if _, ok := remove[i]; ok {
			continue
		}
		keep = append(keep, c)
	}
	cl.changes = keep
	return nil
}

// Clear empties the changelist file.
func (cl *memChangelist) Clear(archive string) error {
	// appending to a nil list initializes it.
	cl.changes = nil
	return nil
}

// Close is a no-op in this in-memory change-list
func (cl *memChangelist) Close() error {
	return nil
}

func (cl *memChangelist) NewIterator() (ChangeIterator, error) {
	return &MemChangeListIterator{index: 0, collection: cl.changes}, nil
}

// MemChangeListIterator is a concrete instance of ChangeIterator
type MemChangeListIterator struct {
	index      int
	collection []Change // Same type as memChangeList.changes
}

// Next returns the next Change
func (m *MemChangeListIterator) Next() (item Change, err error) {
	if m.index >= len(m.collection) {
		return nil, IteratorBoundsError(m.index)
	}
	item = m.collection[m.index]
	m.index++
	return item, err
}

// HasNext indicates whether the iterator is exhausted
func (m *MemChangeListIterator) HasNext() bool {
	return m.index < len(m.collection)
}
