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
package storage

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// PutAppendBlob initializes an empty append blob with specified name. An
// append blob must be created using this method before appending blocks.
//
// See CreateBlockBlobFromReader for more info on creating blobs.
//
// See https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/Put-Blob
func (b *Blob) PutAppendBlob(options *PutBlobOptions) error {
	params := url.Values{}
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeAppend)
	headers = mergeHeaders(headers, headersFromStruct(b.Properties))
	headers = b.Container.bsc.client.addMetadataToHeaders(headers, b.Metadata)

	if options != nil {
		params = addTimeout(params, options.Timeout)
		headers = mergeHeaders(headers, headersFromStruct(*options))
	}
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), params)

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return err
	}
	return b.respondCreation(resp, BlobTypeAppend)
}

// AppendBlockOptions includes the options for an append block operation
type AppendBlockOptions struct {
	Timeout           uint
	LeaseID           string     `header:"x-ms-lease-id"`
	MaxSize           *uint      `header:"x-ms-blob-condition-maxsize"`
	AppendPosition    *uint      `header:"x-ms-blob-condition-appendpos"`
	IfModifiedSince   *time.Time `header:"If-Modified-Since"`
	IfUnmodifiedSince *time.Time `header:"If-Unmodified-Since"`
	IfMatch           string     `header:"If-Match"`
	IfNoneMatch       string     `header:"If-None-Match"`
	RequestID         string     `header:"x-ms-client-request-id"`
	ContentMD5        bool
}

// AppendBlock appends a block to an append blob.
//
// See https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/Append-Block
func (b *Blob) AppendBlock(chunk []byte, options *AppendBlockOptions) error {
	params := url.Values{"comp": {"appendblock"}}
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeAppend)
	headers["Content-Length"] = fmt.Sprintf("%v", len(chunk))

	if options != nil {
		params = addTimeout(params, options.Timeout)
		headers = mergeHeaders(headers, headersFromStruct(*options))
		if options.ContentMD5 {
			md5sum := md5.Sum(chunk)
			headers[headerContentMD5] = base64.StdEncoding.EncodeToString(md5sum[:])
		}
	}
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), params)

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, bytes.NewReader(chunk), b.Container.bsc.auth)
	if err != nil {
		return err
	}
	return b.respondCreation(resp, BlobTypeAppend)
}
