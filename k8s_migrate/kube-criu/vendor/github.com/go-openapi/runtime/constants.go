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
// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

const (
	// HeaderContentType represents a http content-type header, it's value is supposed to be a mime type
	HeaderContentType = "Content-Type"

	// HeaderTransferEncoding represents a http transfer-encoding header.
	HeaderTransferEncoding = "Transfer-Encoding"

	// HeaderAccept the Accept header
	HeaderAccept = "Accept"

	charsetKey = "charset"

	// DefaultMime the default fallback mime type
	DefaultMime = "application/octet-stream"
	// JSONMime the json mime type
	JSONMime = "application/json"
	// YAMLMime the yaml mime type
	YAMLMime = "application/x-yaml"
	// XMLMime the xml mime type
	XMLMime = "application/xml"
	// TextMime the text mime type
	TextMime = "text/plain"
	// MultipartFormMime the multipart form mime type
	MultipartFormMime = "multipart/form-data"
)
