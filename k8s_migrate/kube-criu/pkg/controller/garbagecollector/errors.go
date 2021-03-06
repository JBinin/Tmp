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
/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package garbagecollector

import (
	"fmt"
)

type restMappingError struct {
	kind    string
	version string
}

func (r *restMappingError) Error() string {
	versionKind := fmt.Sprintf("%s/%s", r.version, r.kind)
	return fmt.Sprintf("unable to get REST mapping for %s.", versionKind)
}

// Message prints more details
func (r *restMappingError) Message() string {
	versionKind := fmt.Sprintf("%s/%s", r.version, r.kind)
	errMsg := fmt.Sprintf("unable to get REST mapping for %s. ", versionKind)
	errMsg += fmt.Sprintf(" If %s is an invalid resource, then you should manually remove ownerReferences that refer %s objects.", versionKind, versionKind)
	return errMsg
}

func newRESTMappingError(kind, version string) *restMappingError {
	return &restMappingError{kind: kind, version: version}
}
