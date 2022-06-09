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
package serror

import (
	"encoding/json"
)

func NewTypedStorageOSError(kind StorageOSErrorKind, err error, msg string, help string) StorageOSError {
	return &typedStorageOSError{
		internal: &internal_TypedStorageOSError{
			ErrorKind:   &kind,
			Cause:       err,
			ErrMessage:  msg,
			HelpMessage: help,
		},
	}
}

func NewUntypedStorageOSError(err error, msg string, help string) StorageOSError {
	var kind StorageOSErrorKind = UnknownError

	return &typedStorageOSError{
		internal: &internal_TypedStorageOSError{
			ErrorKind:   &kind,
			Cause:       err,
			ErrMessage:  msg,
			HelpMessage: help,
		},
	}
}

type internal_TypedStorageOSError struct {
	ErrorKind   *StorageOSErrorKind `json:"error_kind"`
	Cause       error               `json:"caused_by"`
	ErrMessage  string              `json:"error_message"`
	HelpMessage string              `json:"help_message"`
	ExtraMap    map[string]string   `json:"extra"`
}

type typedStorageOSError struct {
	internal *internal_TypedStorageOSError
}

func (t *typedStorageOSError) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.internal)
}

func (t *typedStorageOSError) UnmarshalJSON(d []byte) error {
	internal := &internal_TypedStorageOSError{}

	err := json.Unmarshal(d, internal)
	if err != nil {
		return err
	}

	t.internal = internal
	return nil
}

func (t *typedStorageOSError) Error() string            { return t.String() }
func (t *typedStorageOSError) Err() error               { return t.internal.Cause }
func (t *typedStorageOSError) String() string           { return t.internal.ErrMessage }
func (t *typedStorageOSError) Help() string             { return t.internal.HelpMessage }
func (t *typedStorageOSError) Kind() StorageOSErrorKind { return *t.internal.ErrorKind }
func (t *typedStorageOSError) Extra() map[string]string { return t.internal.ExtraMap }
