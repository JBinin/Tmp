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

const (
	MIME_XML   = "application/xml"          // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_JSON  = "application/json"         // Accept or Content-Type used in Consumes() and/or Produces()
	MIME_OCTET = "application/octet-stream" // If Content-Type is not present in request, use the default

	HEADER_Allow                         = "Allow"
	HEADER_Accept                        = "Accept"
	HEADER_Origin                        = "Origin"
	HEADER_ContentType                   = "Content-Type"
	HEADER_LastModified                  = "Last-Modified"
	HEADER_AcceptEncoding                = "Accept-Encoding"
	HEADER_ContentEncoding               = "Content-Encoding"
	HEADER_AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HEADER_AccessControlRequestMethod    = "Access-Control-Request-Method"
	HEADER_AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HEADER_AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HEADER_AccessControlMaxAge           = "Access-Control-Max-Age"

	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"
)
