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
package query

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
)

type xmlErrorResponse struct {
	XMLName   xml.Name `xml:"ErrorResponse"`
	Code      string   `xml:"Error>Code"`
	Message   string   `xml:"Error>Message"`
	RequestID string   `xml:"RequestId"`
}

type xmlServiceUnavailableResponse struct {
	XMLName xml.Name `xml:"ServiceUnavailableException"`
}

// UnmarshalErrorHandler is a name request handler to unmarshal request errors
var UnmarshalErrorHandler = request.NamedHandler{Name: "awssdk.query.UnmarshalError", Fn: UnmarshalError}

// UnmarshalError unmarshals an error response for an AWS Query service.
func UnmarshalError(r *request.Request) {
	defer r.HTTPResponse.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.HTTPResponse.Body)
	if err != nil {
		r.Error = awserr.New("SerializationError", "failed to read from query HTTP response body", err)
		return
	}

	// First check for specific error
	resp := xmlErrorResponse{}
	decodeErr := xml.Unmarshal(bodyBytes, &resp)
	if decodeErr == nil {
		reqID := resp.RequestID
		if reqID == "" {
			reqID = r.RequestID
		}
		r.Error = awserr.NewRequestFailure(
			awserr.New(resp.Code, resp.Message, nil),
			r.HTTPResponse.StatusCode,
			reqID,
		)
		return
	}

	// Check for unhandled error
	servUnavailResp := xmlServiceUnavailableResponse{}
	unavailErr := xml.Unmarshal(bodyBytes, &servUnavailResp)
	if unavailErr == nil {
		r.Error = awserr.NewRequestFailure(
			awserr.New("ServiceUnavailableException", "service is unavailable", nil),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}

	// Failed to retrieve any error message from the response body
	r.Error = awserr.New("SerializationError",
		"failed to decode query XML error response", decodeErr)
}
