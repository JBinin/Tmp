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
package runtime

import (
	"errors"
	"net/http"
)

// MIMEWildcard is the fallback MIME type used for requests which do not match
// a registered MIME type.
const MIMEWildcard = "*"

var (
	acceptHeader      = http.CanonicalHeaderKey("Accept")
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")

	defaultMarshaler = &JSONPb{OrigName: true}
)

// MarshalerForRequest returns the inbound/outbound marshalers for this request.
// It checks the registry on the ServeMux for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
func MarshalerForRequest(mux *ServeMux, r *http.Request) (inbound Marshaler, outbound Marshaler) {
	for _, acceptVal := range r.Header[acceptHeader] {
		if m, ok := mux.marshalers.mimeMap[acceptVal]; ok {
			outbound = m
			break
		}
	}

	for _, contentTypeVal := range r.Header[contentTypeHeader] {
		if m, ok := mux.marshalers.mimeMap[contentTypeVal]; ok {
			inbound = m
			break
		}
	}

	if inbound == nil {
		inbound = mux.marshalers.mimeMap[MIMEWildcard]
	}
	if outbound == nil {
		outbound = inbound
	}

	return inbound, outbound
}

// marshalerRegistry is a mapping from MIME types to Marshalers.
type marshalerRegistry struct {
	mimeMap map[string]Marshaler
}

// add adds a marshaler for a case-sensitive MIME type string ("*" to match any
// MIME type).
func (m marshalerRegistry) add(mime string, marshaler Marshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}

	m.mimeMap[mime] = marshaler

	return nil
}

// makeMarshalerMIMERegistry returns a new registry of marshalers.
// It allows for a mapping of case-sensitive Content-Type MIME type string to runtime.Marshaler interfaces.
//
// For example, you could allow the client to specify the use of the runtime.JSONPb marshaler
// with a "applicaton/jsonpb" Content-Type and the use of the runtime.JSONBuiltin marshaler
// with a "application/json" Content-Type.
// "*" can be used to match any Content-Type.
// This can be attached to a ServerMux with the marshaler option.
func makeMarshalerMIMERegistry() marshalerRegistry {
	return marshalerRegistry{
		mimeMap: map[string]Marshaler{
			MIMEWildcard: defaultMarshaler,
		},
	}
}

// WithMarshalerOption returns a ServeMuxOption which associates inbound and outbound
// Marshalers to a MIME type in mux.
func WithMarshalerOption(mime string, marshaler Marshaler) ServeMuxOption {
	return func(mux *ServeMux) {
		if err := mux.marshalers.add(mime, marshaler); err != nil {
			panic(err)
		}
	}
}
