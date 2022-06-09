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
package store

import (
	"strings"
)

// CreateEndpoints creates a list of endpoints given the right scheme
func CreateEndpoints(addrs []string, scheme string) (entries []string) {
	for _, addr := range addrs {
		entries = append(entries, scheme+"://"+addr)
	}
	return entries
}

// Normalize the key for each store to the form:
//
//     /path/to/key
//
func Normalize(key string) string {
	return "/" + join(SplitKey(key))
}

// GetDirectory gets the full directory part of
// the key to the form:
//
//     /path/to/
//
func GetDirectory(key string) string {
	parts := SplitKey(key)
	parts = parts[:len(parts)-1]
	return "/" + join(parts)
}

// SplitKey splits the key to extract path informations
func SplitKey(key string) (path []string) {
	if strings.Contains(key, "/") {
		path = strings.Split(key, "/")
	} else {
		path = []string{key}
	}
	return path
}

// join the path parts with '/'
func join(parts []string) string {
	return strings.Join(parts, "/")
}
