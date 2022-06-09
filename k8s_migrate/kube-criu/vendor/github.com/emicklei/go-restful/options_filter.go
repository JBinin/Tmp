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

import "strings"

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

// OPTIONSFilter is a filter function that inspects the Http Request for the OPTIONS method
// and provides the response with a set of allowed methods for the request URL Path.
// As for any filter, you can also install it for a particular WebService within a Container.
// Note: this filter is not needed when using CrossOriginResourceSharing (for CORS).
func (c *Container) OPTIONSFilter(req *Request, resp *Response, chain *FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}
	resp.AddHeader(HEADER_Allow, strings.Join(c.computeAllowedMethods(req), ","))
}

// OPTIONSFilter is a filter function that inspects the Http Request for the OPTIONS method
// and provides the response with a set of allowed methods for the request URL Path.
// Note: this filter is not needed when using CrossOriginResourceSharing (for CORS).
func OPTIONSFilter() FilterFunction {
	return DefaultContainer.OPTIONSFilter
}
