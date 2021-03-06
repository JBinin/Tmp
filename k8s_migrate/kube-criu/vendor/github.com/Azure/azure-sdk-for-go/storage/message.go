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
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Message represents an Azure message.
type Message struct {
	Queue        *Queue
	Text         string      `xml:"MessageText"`
	ID           string      `xml:"MessageId"`
	Insertion    TimeRFC1123 `xml:"InsertionTime"`
	Expiration   TimeRFC1123 `xml:"ExpirationTime"`
	PopReceipt   string      `xml:"PopReceipt"`
	NextVisible  TimeRFC1123 `xml:"TimeNextVisible"`
	DequeueCount int         `xml:"DequeueCount"`
}

func (m *Message) buildPath() string {
	return fmt.Sprintf("%s/%s", m.Queue.buildPathMessages(), m.ID)
}

// PutMessageOptions is the set of options can be specified for Put Messsage
// operation. A zero struct does not use any preferences for the request.
type PutMessageOptions struct {
	Timeout           uint
	VisibilityTimeout int
	MessageTTL        int
	RequestID         string `header:"x-ms-client-request-id"`
}

// Put operation adds a new message to the back of the message queue.
//
// See https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/Put-Message
func (m *Message) Put(options *PutMessageOptions) error {
	query := url.Values{}
	headers := m.Queue.qsc.client.getStandardHeaders()

	req := putMessageRequest{MessageText: m.Text}
	body, nn, err := xmlMarshal(req)
	if err != nil {
		return err
	}
	headers["Content-Length"] = strconv.Itoa(nn)

	if options != nil {
		if options.VisibilityTimeout != 0 {
			query.Set("visibilitytimeout", strconv.Itoa(options.VisibilityTimeout))
		}
		if options.MessageTTL != 0 {
			query.Set("messagettl", strconv.Itoa(options.MessageTTL))
		}
		query = addTimeout(query, options.Timeout)
		headers = mergeHeaders(headers, headersFromStruct(*options))
	}

	uri := m.Queue.qsc.client.getEndpoint(queueServiceName, m.Queue.buildPathMessages(), query)
	resp, err := m.Queue.qsc.client.exec(http.MethodPost, uri, headers, body, m.Queue.qsc.auth)
	if err != nil {
		return err
	}
	defer drainRespBody(resp)
	err = checkRespCode(resp, []int{http.StatusCreated})
	if err != nil {
		return err
	}
	err = xmlUnmarshal(resp.Body, m)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMessageOptions is the set of options can be specified for Update Messsage
// operation. A zero struct does not use any preferences for the request.
type UpdateMessageOptions struct {
	Timeout           uint
	VisibilityTimeout int
	RequestID         string `header:"x-ms-client-request-id"`
}

// Update operation updates the specified message.
//
// See https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/Update-Message
func (m *Message) Update(options *UpdateMessageOptions) error {
	query := url.Values{}
	if m.PopReceipt != "" {
		query.Set("popreceipt", m.PopReceipt)
	}

	headers := m.Queue.qsc.client.getStandardHeaders()
	req := putMessageRequest{MessageText: m.Text}
	body, nn, err := xmlMarshal(req)
	if err != nil {
		return err
	}
	headers["Content-Length"] = strconv.Itoa(nn)
	// visibilitytimeout is required for Update (zero or greater) so set the default here
	query.Set("visibilitytimeout", "0")
	if options != nil {
		if options.VisibilityTimeout != 0 {
			query.Set("visibilitytimeout", strconv.Itoa(options.VisibilityTimeout))
		}
		query = addTimeout(query, options.Timeout)
		headers = mergeHeaders(headers, headersFromStruct(*options))
	}
	uri := m.Queue.qsc.client.getEndpoint(queueServiceName, m.buildPath(), query)

	resp, err := m.Queue.qsc.client.exec(http.MethodPut, uri, headers, body, m.Queue.qsc.auth)
	if err != nil {
		return err
	}
	defer drainRespBody(resp)

	m.PopReceipt = resp.Header.Get("x-ms-popreceipt")
	nextTimeStr := resp.Header.Get("x-ms-time-next-visible")
	if nextTimeStr != "" {
		nextTime, err := time.Parse(time.RFC1123, nextTimeStr)
		if err != nil {
			return err
		}
		m.NextVisible = TimeRFC1123(nextTime)
	}

	return checkRespCode(resp, []int{http.StatusNoContent})
}

// Delete operation deletes the specified message.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179347.aspx
func (m *Message) Delete(options *QueueServiceOptions) error {
	params := url.Values{"popreceipt": {m.PopReceipt}}
	headers := m.Queue.qsc.client.getStandardHeaders()

	if options != nil {
		params = addTimeout(params, options.Timeout)
		headers = mergeHeaders(headers, headersFromStruct(*options))
	}
	uri := m.Queue.qsc.client.getEndpoint(queueServiceName, m.buildPath(), params)

	resp, err := m.Queue.qsc.client.exec(http.MethodDelete, uri, headers, nil, m.Queue.qsc.auth)
	if err != nil {
		return err
	}
	defer drainRespBody(resp)
	return checkRespCode(resp, []int{http.StatusNoContent})
}

type putMessageRequest struct {
	XMLName     xml.Name `xml:"QueueMessage"`
	MessageText string   `xml:"MessageText"`
}
