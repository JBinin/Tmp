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
package container

// StoreFilter defines a function to filter
// container in the store.
type StoreFilter func(*Container) bool

// StoreReducer defines a function to
// manipulate containers in the store
type StoreReducer func(*Container)

// Store defines an interface that
// any container store must implement.
type Store interface {
	// Add appends a new container to the store.
	Add(string, *Container)
	// Get returns a container from the store by the identifier it was stored with.
	Get(string) *Container
	// Delete removes a container from the store by the identifier it was stored with.
	Delete(string)
	// List returns a list of containers from the store.
	List() []*Container
	// Size returns the number of containers in the store.
	Size() int
	// First returns the first container found in the store by a given filter.
	First(StoreFilter) *Container
	// ApplyAll calls the reducer function with every container in the store.
	ApplyAll(StoreReducer)
}
