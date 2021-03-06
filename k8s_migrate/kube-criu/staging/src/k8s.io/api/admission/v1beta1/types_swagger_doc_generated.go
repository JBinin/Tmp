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
Copyright The Kubernetes Authors.

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

package v1beta1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
var map_AdmissionRequest = map[string]string{
	"":            "AdmissionRequest describes the admission.Attributes for the admission request.",
	"uid":         "UID is an identifier for the individual request/response. It allows us to distinguish instances of requests which are otherwise identical (parallel requests, requests when earlier requests did not modify etc) The UID is meant to track the round trip (request/response) between the KAS and the WebHook, not the user request. It is suitable for correlating log entries between the webhook and apiserver, for either auditing or debugging.",
	"kind":        "Kind is the type of object being manipulated.  For example: Pod",
	"resource":    "Resource is the name of the resource being requested.  This is not the kind.  For example: pods",
	"subResource": "SubResource is the name of the subresource being requested.  This is a different resource, scoped to the parent resource, but it may have a different kind. For instance, /pods has the resource \"pods\" and the kind \"Pod\", while /pods/foo/status has the resource \"pods\", the sub resource \"status\", and the kind \"Pod\" (because status operates on pods). The binding resource for a pod though may be /pods/foo/binding, which has resource \"pods\", subresource \"binding\", and kind \"Binding\".",
	"name":        "Name is the name of the object as presented in the request.  On a CREATE operation, the client may omit name and rely on the server to generate the name.  If that is the case, this method will return the empty string.",
	"namespace":   "Namespace is the namespace associated with the request (if any).",
	"operation":   "Operation is the operation being performed",
	"userInfo":    "UserInfo is information about the requesting user",
	"object":      "Object is the object from the incoming request prior to default values being applied",
	"oldObject":   "OldObject is the existing object. Only populated for UPDATE requests.",
}

func (AdmissionRequest) SwaggerDoc() map[string]string {
	return map_AdmissionRequest
}

var map_AdmissionResponse = map[string]string{
	"":          "AdmissionResponse describes an admission response.",
	"uid":       "UID is an identifier for the individual request/response. This should be copied over from the corresponding AdmissionRequest.",
	"allowed":   "Allowed indicates whether or not the admission request was permitted.",
	"status":    "Result contains extra details into why an admission request was denied. This field IS NOT consulted in any way if \"Allowed\" is \"true\".",
	"patch":     "The patch body. Currently we only support \"JSONPatch\" which implements RFC 6902.",
	"patchType": "The type of Patch. Currently we only allow \"JSONPatch\".",
}

func (AdmissionResponse) SwaggerDoc() map[string]string {
	return map_AdmissionResponse
}

var map_AdmissionReview = map[string]string{
	"":         "AdmissionReview describes an admission review request/response.",
	"request":  "Request describes the attributes for the admission request.",
	"response": "Response describes the attributes for the admission response.",
}

func (AdmissionReview) SwaggerDoc() map[string]string {
	return map_AdmissionReview
}

// AUTO-GENERATED FUNCTIONS END HERE
