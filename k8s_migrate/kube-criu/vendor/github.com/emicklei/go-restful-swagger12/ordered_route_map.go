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
package swagger

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "github.com/emicklei/go-restful"

type orderedRouteMap struct {
	elements map[string][]restful.Route
	keys     []string
}

func newOrderedRouteMap() *orderedRouteMap {
	return &orderedRouteMap{
		elements: map[string][]restful.Route{},
		keys:     []string{},
	}
}

func (o *orderedRouteMap) Add(key string, route restful.Route) {
	routes, ok := o.elements[key]
	if ok {
		routes = append(routes, route)
		o.elements[key] = routes
		return
	}
	o.elements[key] = []restful.Route{route}
	o.keys = append(o.keys, key)
}

func (o *orderedRouteMap) Do(block func(key string, routes []restful.Route)) {
	for _, k := range o.keys {
		block(k, o.elements[k])
	}
}
