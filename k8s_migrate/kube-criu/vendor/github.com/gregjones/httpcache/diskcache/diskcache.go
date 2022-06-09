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
// Package diskcache provides an implementation of httpcache.Cache that uses the diskv package
// to supplement an in-memory map with persistent storage
//
package diskcache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/peterbourgon/diskv"
	"io"
)

// Cache is an implementation of httpcache.Cache that supplements the in-memory map with persistent storage
type Cache struct {
	d *diskv.Diskv
}

// Get returns the response corresponding to key if present
func (c *Cache) Get(key string) (resp []byte, ok bool) {
	key = keyToFilename(key)
	resp, err := c.d.Read(key)
	if err != nil {
		return []byte{}, false
	}
	return resp, true
}

// Set saves a response to the cache as key
func (c *Cache) Set(key string, resp []byte) {
	key = keyToFilename(key)
	c.d.WriteStream(key, bytes.NewReader(resp), true)
}

// Delete removes the response with key from the cache
func (c *Cache) Delete(key string) {
	key = keyToFilename(key)
	c.d.Erase(key)
}

func keyToFilename(key string) string {
	h := md5.New()
	io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}

// New returns a new Cache that will store files in basePath
func New(basePath string) *Cache {
	return &Cache{
		d: diskv.New(diskv.Options{
			BasePath:     basePath,
			CacheSizeMax: 100 * 1024 * 1024, // 100MB
		}),
	}
}

// NewWithDiskv returns a new Cache using the provided Diskv as underlying
// storage.
func NewWithDiskv(d *diskv.Diskv) *Cache {
	return &Cache{d}
}
