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
package watch

import "sort"

type Delta struct {
	ModifiedPackages []string

	NewSuites      []*Suite
	RemovedSuites  []*Suite
	modifiedSuites []*Suite
}

type DescendingByDelta []*Suite

func (a DescendingByDelta) Len() int           { return len(a) }
func (a DescendingByDelta) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DescendingByDelta) Less(i, j int) bool { return a[i].Delta() > a[j].Delta() }

func (d Delta) ModifiedSuites() []*Suite {
	sort.Sort(DescendingByDelta(d.modifiedSuites))
	return d.modifiedSuites
}
