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
Copyright 2016 The Kubernetes Authors.

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

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TokenReview attempts to authenticate a token to a known user.
// Note: TokenReview requests may be cached by the webhook token authenticator
// plugin in the kube-apiserver.
type TokenReview struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec holds information about the request being evaluated
	Spec TokenReviewSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Status is filled in by the server and indicates whether the request can be authenticated.
	// +optional
	Status TokenReviewStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// TokenReviewSpec is a description of the token authentication request.
type TokenReviewSpec struct {
	// Token is the opaque bearer token.
	// +optional
	Token string `json:"token,omitempty" protobuf:"bytes,1,opt,name=token"`
}

// TokenReviewStatus is the result of the token authentication request.
type TokenReviewStatus struct {
	// Authenticated indicates that the token was associated with a known user.
	// +optional
	Authenticated bool `json:"authenticated,omitempty" protobuf:"varint,1,opt,name=authenticated"`
	// User is the UserInfo associated with the provided token.
	// +optional
	User UserInfo `json:"user,omitempty" protobuf:"bytes,2,opt,name=user"`
	// Error indicates that the token couldn't be checked
	// +optional
	Error string `json:"error,omitempty" protobuf:"bytes,3,opt,name=error"`
}

// UserInfo holds the information about the user needed to implement the
// user.Info interface.
type UserInfo struct {
	// The name that uniquely identifies this user among all active users.
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,1,opt,name=username"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	// +optional
	UID string `json:"uid,omitempty" protobuf:"bytes,2,opt,name=uid"`
	// The names of groups this user is a part of.
	// +optional
	Groups []string `json:"groups,omitempty" protobuf:"bytes,3,rep,name=groups"`
	// Any additional information provided by the authenticator.
	// +optional
	Extra map[string]ExtraValue `json:"extra,omitempty" protobuf:"bytes,4,rep,name=extra"`
}

// ExtraValue masks the value so protobuf can generate
// +protobuf.nullable=true
// +protobuf.options.(gogoproto.goproto_stringer)=false
type ExtraValue []string

func (t ExtraValue) String() string {
	return fmt.Sprintf("%v", []string(t))
}
