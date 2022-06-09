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
package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "net/http"

// A RouteSelector finds the best matching Route given the input HTTP Request
type RouteSelector interface {

	// SelectRoute finds a Route given the input HTTP Request and a list of WebServices.
	// It returns a selected Route and its containing WebService or an error indicating
	// a problem.
	SelectRoute(
		webServices []*WebService,
		httpRequest *http.Request) (selectedService *WebService, selected *Route, err error)
}
