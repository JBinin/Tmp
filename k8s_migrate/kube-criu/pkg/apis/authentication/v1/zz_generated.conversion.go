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
// +build !ignore_autogenerated

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1

import (
	unsafe "unsafe"

	v1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	authentication "k8s.io/kubernetes/pkg/apis/authentication"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.BoundObjectReference)(nil), (*authentication.BoundObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BoundObjectReference_To_authentication_BoundObjectReference(a.(*v1.BoundObjectReference), b.(*authentication.BoundObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.BoundObjectReference)(nil), (*v1.BoundObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_BoundObjectReference_To_v1_BoundObjectReference(a.(*authentication.BoundObjectReference), b.(*v1.BoundObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenRequest)(nil), (*authentication.TokenRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenRequest_To_authentication_TokenRequest(a.(*v1.TokenRequest), b.(*authentication.TokenRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.TokenRequest)(nil), (*v1.TokenRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_TokenRequest_To_v1_TokenRequest(a.(*authentication.TokenRequest), b.(*v1.TokenRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenRequestSpec)(nil), (*authentication.TokenRequestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenRequestSpec_To_authentication_TokenRequestSpec(a.(*v1.TokenRequestSpec), b.(*authentication.TokenRequestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.TokenRequestSpec)(nil), (*v1.TokenRequestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_TokenRequestSpec_To_v1_TokenRequestSpec(a.(*authentication.TokenRequestSpec), b.(*v1.TokenRequestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenRequestStatus)(nil), (*authentication.TokenRequestStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenRequestStatus_To_authentication_TokenRequestStatus(a.(*v1.TokenRequestStatus), b.(*authentication.TokenRequestStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.TokenRequestStatus)(nil), (*v1.TokenRequestStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_TokenRequestStatus_To_v1_TokenRequestStatus(a.(*authentication.TokenRequestStatus), b.(*v1.TokenRequestStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenReview)(nil), (*authentication.TokenReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenReview_To_authentication_TokenReview(a.(*v1.TokenReview), b.(*authentication.TokenReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.TokenReview)(nil), (*v1.TokenReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_TokenReview_To_v1_TokenReview(a.(*authentication.TokenReview), b.(*v1.TokenReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenReviewSpec)(nil), (*authentication.TokenReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenReviewSpec_To_authentication_TokenReviewSpec(a.(*v1.TokenReviewSpec), b.(*authentication.TokenReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.TokenReviewSpec)(nil), (*v1.TokenReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_TokenReviewSpec_To_v1_TokenReviewSpec(a.(*authentication.TokenReviewSpec), b.(*v1.TokenReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenReviewStatus)(nil), (*authentication.TokenReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenReviewStatus_To_authentication_TokenReviewStatus(a.(*v1.TokenReviewStatus), b.(*authentication.TokenReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.TokenReviewStatus)(nil), (*v1.TokenReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_TokenReviewStatus_To_v1_TokenReviewStatus(a.(*authentication.TokenReviewStatus), b.(*v1.TokenReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.UserInfo)(nil), (*authentication.UserInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_UserInfo_To_authentication_UserInfo(a.(*v1.UserInfo), b.(*authentication.UserInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.UserInfo)(nil), (*v1.UserInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_UserInfo_To_v1_UserInfo(a.(*authentication.UserInfo), b.(*v1.UserInfo), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_BoundObjectReference_To_authentication_BoundObjectReference(in *v1.BoundObjectReference, out *authentication.BoundObjectReference, s conversion.Scope) error {
	out.Kind = in.Kind
	out.APIVersion = in.APIVersion
	out.Name = in.Name
	out.UID = types.UID(in.UID)
	return nil
}

// Convert_v1_BoundObjectReference_To_authentication_BoundObjectReference is an autogenerated conversion function.
func Convert_v1_BoundObjectReference_To_authentication_BoundObjectReference(in *v1.BoundObjectReference, out *authentication.BoundObjectReference, s conversion.Scope) error {
	return autoConvert_v1_BoundObjectReference_To_authentication_BoundObjectReference(in, out, s)
}

func autoConvert_authentication_BoundObjectReference_To_v1_BoundObjectReference(in *authentication.BoundObjectReference, out *v1.BoundObjectReference, s conversion.Scope) error {
	out.Kind = in.Kind
	out.APIVersion = in.APIVersion
	out.Name = in.Name
	out.UID = types.UID(in.UID)
	return nil
}

// Convert_authentication_BoundObjectReference_To_v1_BoundObjectReference is an autogenerated conversion function.
func Convert_authentication_BoundObjectReference_To_v1_BoundObjectReference(in *authentication.BoundObjectReference, out *v1.BoundObjectReference, s conversion.Scope) error {
	return autoConvert_authentication_BoundObjectReference_To_v1_BoundObjectReference(in, out, s)
}

func autoConvert_v1_TokenRequest_To_authentication_TokenRequest(in *v1.TokenRequest, out *authentication.TokenRequest, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_TokenRequestSpec_To_authentication_TokenRequestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_TokenRequestStatus_To_authentication_TokenRequestStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_TokenRequest_To_authentication_TokenRequest is an autogenerated conversion function.
func Convert_v1_TokenRequest_To_authentication_TokenRequest(in *v1.TokenRequest, out *authentication.TokenRequest, s conversion.Scope) error {
	return autoConvert_v1_TokenRequest_To_authentication_TokenRequest(in, out, s)
}

func autoConvert_authentication_TokenRequest_To_v1_TokenRequest(in *authentication.TokenRequest, out *v1.TokenRequest, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_authentication_TokenRequestSpec_To_v1_TokenRequestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_authentication_TokenRequestStatus_To_v1_TokenRequestStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_authentication_TokenRequest_To_v1_TokenRequest is an autogenerated conversion function.
func Convert_authentication_TokenRequest_To_v1_TokenRequest(in *authentication.TokenRequest, out *v1.TokenRequest, s conversion.Scope) error {
	return autoConvert_authentication_TokenRequest_To_v1_TokenRequest(in, out, s)
}

func autoConvert_v1_TokenRequestSpec_To_authentication_TokenRequestSpec(in *v1.TokenRequestSpec, out *authentication.TokenRequestSpec, s conversion.Scope) error {
	out.Audiences = *(*[]string)(unsafe.Pointer(&in.Audiences))
	if err := metav1.Convert_Pointer_int64_To_int64(&in.ExpirationSeconds, &out.ExpirationSeconds, s); err != nil {
		return err
	}
	out.BoundObjectRef = (*authentication.BoundObjectReference)(unsafe.Pointer(in.BoundObjectRef))
	return nil
}

// Convert_v1_TokenRequestSpec_To_authentication_TokenRequestSpec is an autogenerated conversion function.
func Convert_v1_TokenRequestSpec_To_authentication_TokenRequestSpec(in *v1.TokenRequestSpec, out *authentication.TokenRequestSpec, s conversion.Scope) error {
	return autoConvert_v1_TokenRequestSpec_To_authentication_TokenRequestSpec(in, out, s)
}

func autoConvert_authentication_TokenRequestSpec_To_v1_TokenRequestSpec(in *authentication.TokenRequestSpec, out *v1.TokenRequestSpec, s conversion.Scope) error {
	out.Audiences = *(*[]string)(unsafe.Pointer(&in.Audiences))
	if err := metav1.Convert_int64_To_Pointer_int64(&in.ExpirationSeconds, &out.ExpirationSeconds, s); err != nil {
		return err
	}
	out.BoundObjectRef = (*v1.BoundObjectReference)(unsafe.Pointer(in.BoundObjectRef))
	return nil
}

// Convert_authentication_TokenRequestSpec_To_v1_TokenRequestSpec is an autogenerated conversion function.
func Convert_authentication_TokenRequestSpec_To_v1_TokenRequestSpec(in *authentication.TokenRequestSpec, out *v1.TokenRequestSpec, s conversion.Scope) error {
	return autoConvert_authentication_TokenRequestSpec_To_v1_TokenRequestSpec(in, out, s)
}

func autoConvert_v1_TokenRequestStatus_To_authentication_TokenRequestStatus(in *v1.TokenRequestStatus, out *authentication.TokenRequestStatus, s conversion.Scope) error {
	out.Token = in.Token
	out.ExpirationTimestamp = in.ExpirationTimestamp
	return nil
}

// Convert_v1_TokenRequestStatus_To_authentication_TokenRequestStatus is an autogenerated conversion function.
func Convert_v1_TokenRequestStatus_To_authentication_TokenRequestStatus(in *v1.TokenRequestStatus, out *authentication.TokenRequestStatus, s conversion.Scope) error {
	return autoConvert_v1_TokenRequestStatus_To_authentication_TokenRequestStatus(in, out, s)
}

func autoConvert_authentication_TokenRequestStatus_To_v1_TokenRequestStatus(in *authentication.TokenRequestStatus, out *v1.TokenRequestStatus, s conversion.Scope) error {
	out.Token = in.Token
	out.ExpirationTimestamp = in.ExpirationTimestamp
	return nil
}

// Convert_authentication_TokenRequestStatus_To_v1_TokenRequestStatus is an autogenerated conversion function.
func Convert_authentication_TokenRequestStatus_To_v1_TokenRequestStatus(in *authentication.TokenRequestStatus, out *v1.TokenRequestStatus, s conversion.Scope) error {
	return autoConvert_authentication_TokenRequestStatus_To_v1_TokenRequestStatus(in, out, s)
}

func autoConvert_v1_TokenReview_To_authentication_TokenReview(in *v1.TokenReview, out *authentication.TokenReview, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_TokenReviewSpec_To_authentication_TokenReviewSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_TokenReviewStatus_To_authentication_TokenReviewStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_TokenReview_To_authentication_TokenReview is an autogenerated conversion function.
func Convert_v1_TokenReview_To_authentication_TokenReview(in *v1.TokenReview, out *authentication.TokenReview, s conversion.Scope) error {
	return autoConvert_v1_TokenReview_To_authentication_TokenReview(in, out, s)
}

func autoConvert_authentication_TokenReview_To_v1_TokenReview(in *authentication.TokenReview, out *v1.TokenReview, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_authentication_TokenReviewSpec_To_v1_TokenReviewSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_authentication_TokenReviewStatus_To_v1_TokenReviewStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_authentication_TokenReview_To_v1_TokenReview is an autogenerated conversion function.
func Convert_authentication_TokenReview_To_v1_TokenReview(in *authentication.TokenReview, out *v1.TokenReview, s conversion.Scope) error {
	return autoConvert_authentication_TokenReview_To_v1_TokenReview(in, out, s)
}

func autoConvert_v1_TokenReviewSpec_To_authentication_TokenReviewSpec(in *v1.TokenReviewSpec, out *authentication.TokenReviewSpec, s conversion.Scope) error {
	out.Token = in.Token
	return nil
}

// Convert_v1_TokenReviewSpec_To_authentication_TokenReviewSpec is an autogenerated conversion function.
func Convert_v1_TokenReviewSpec_To_authentication_TokenReviewSpec(in *v1.TokenReviewSpec, out *authentication.TokenReviewSpec, s conversion.Scope) error {
	return autoConvert_v1_TokenReviewSpec_To_authentication_TokenReviewSpec(in, out, s)
}

func autoConvert_authentication_TokenReviewSpec_To_v1_TokenReviewSpec(in *authentication.TokenReviewSpec, out *v1.TokenReviewSpec, s conversion.Scope) error {
	out.Token = in.Token
	return nil
}

// Convert_authentication_TokenReviewSpec_To_v1_TokenReviewSpec is an autogenerated conversion function.
func Convert_authentication_TokenReviewSpec_To_v1_TokenReviewSpec(in *authentication.TokenReviewSpec, out *v1.TokenReviewSpec, s conversion.Scope) error {
	return autoConvert_authentication_TokenReviewSpec_To_v1_TokenReviewSpec(in, out, s)
}

func autoConvert_v1_TokenReviewStatus_To_authentication_TokenReviewStatus(in *v1.TokenReviewStatus, out *authentication.TokenReviewStatus, s conversion.Scope) error {
	out.Authenticated = in.Authenticated
	if err := Convert_v1_UserInfo_To_authentication_UserInfo(&in.User, &out.User, s); err != nil {
		return err
	}
	out.Error = in.Error
	return nil
}

// Convert_v1_TokenReviewStatus_To_authentication_TokenReviewStatus is an autogenerated conversion function.
func Convert_v1_TokenReviewStatus_To_authentication_TokenReviewStatus(in *v1.TokenReviewStatus, out *authentication.TokenReviewStatus, s conversion.Scope) error {
	return autoConvert_v1_TokenReviewStatus_To_authentication_TokenReviewStatus(in, out, s)
}

func autoConvert_authentication_TokenReviewStatus_To_v1_TokenReviewStatus(in *authentication.TokenReviewStatus, out *v1.TokenReviewStatus, s conversion.Scope) error {
	out.Authenticated = in.Authenticated
	if err := Convert_authentication_UserInfo_To_v1_UserInfo(&in.User, &out.User, s); err != nil {
		return err
	}
	out.Error = in.Error
	return nil
}

// Convert_authentication_TokenReviewStatus_To_v1_TokenReviewStatus is an autogenerated conversion function.
func Convert_authentication_TokenReviewStatus_To_v1_TokenReviewStatus(in *authentication.TokenReviewStatus, out *v1.TokenReviewStatus, s conversion.Scope) error {
	return autoConvert_authentication_TokenReviewStatus_To_v1_TokenReviewStatus(in, out, s)
}

func autoConvert_v1_UserInfo_To_authentication_UserInfo(in *v1.UserInfo, out *authentication.UserInfo, s conversion.Scope) error {
	out.Username = in.Username
	out.UID = in.UID
	out.Groups = *(*[]string)(unsafe.Pointer(&in.Groups))
	out.Extra = *(*map[string]authentication.ExtraValue)(unsafe.Pointer(&in.Extra))
	return nil
}

// Convert_v1_UserInfo_To_authentication_UserInfo is an autogenerated conversion function.
func Convert_v1_UserInfo_To_authentication_UserInfo(in *v1.UserInfo, out *authentication.UserInfo, s conversion.Scope) error {
	return autoConvert_v1_UserInfo_To_authentication_UserInfo(in, out, s)
}

func autoConvert_authentication_UserInfo_To_v1_UserInfo(in *authentication.UserInfo, out *v1.UserInfo, s conversion.Scope) error {
	out.Username = in.Username
	out.UID = in.UID
	out.Groups = *(*[]string)(unsafe.Pointer(&in.Groups))
	out.Extra = *(*map[string]v1.ExtraValue)(unsafe.Pointer(&in.Extra))
	return nil
}

// Convert_authentication_UserInfo_To_v1_UserInfo is an autogenerated conversion function.
func Convert_authentication_UserInfo_To_v1_UserInfo(in *authentication.UserInfo, out *v1.UserInfo, s conversion.Scope) error {
	return autoConvert_authentication_UserInfo_To_v1_UserInfo(in, out, s)
}
