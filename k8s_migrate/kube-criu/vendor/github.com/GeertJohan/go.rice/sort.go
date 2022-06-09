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
package rice

import "os"

// SortByName allows an array of os.FileInfo objects
// to be easily sorted by filename using sort.Sort(SortByName(array))
type SortByName []os.FileInfo

func (f SortByName) Len() int           { return len(f) }
func (f SortByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f SortByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// SortByModified allows an array of os.FileInfo objects
// to be easily sorted by modified date using sort.Sort(SortByModified(array))
type SortByModified []os.FileInfo

func (f SortByModified) Len() int           { return len(f) }
func (f SortByModified) Less(i, j int) bool { return f[i].ModTime().Unix() > f[j].ModTime().Unix() }
func (f SortByModified) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
