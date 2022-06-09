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
package errcode

import (
	"encoding/json"
	"net/http"
)

// ServeJSON attempts to serve the errcode in a JSON envelope. It marshals err
// and sets the content-type header to 'application/json'. It will handle
// ErrorCoder and Errors, and if necessary will create an envelope.
func ServeJSON(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var sc int

	switch errs := err.(type) {
	case Errors:
		if len(errs) < 1 {
			break
		}

		if err, ok := errs[0].(ErrorCoder); ok {
			sc = err.ErrorCode().Descriptor().HTTPStatusCode
		}
	case ErrorCoder:
		sc = errs.ErrorCode().Descriptor().HTTPStatusCode
		err = Errors{err} // create an envelope.
	default:
		// We just have an unhandled error type, so just place in an envelope
		// and move along.
		err = Errors{err}
	}

	if sc == 0 {
		sc = http.StatusInternalServerError
	}

	w.WriteHeader(sc)

	if err := json.NewEncoder(w).Encode(err); err != nil {
		return err
	}

	return nil
}
