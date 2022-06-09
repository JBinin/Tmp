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

package v1alpha1

import (
	unsafe "unsafe"

	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/kubernetes/pkg/apis/core"
	resourcequota "k8s.io/kubernetes/plugin/pkg/admission/resourcequota/apis/resourcequota"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Configuration)(nil), (*resourcequota.Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Configuration_To_resourcequota_Configuration(a.(*Configuration), b.(*resourcequota.Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*resourcequota.Configuration)(nil), (*Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_resourcequota_Configuration_To_v1alpha1_Configuration(a.(*resourcequota.Configuration), b.(*Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LimitedResource)(nil), (*resourcequota.LimitedResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LimitedResource_To_resourcequota_LimitedResource(a.(*LimitedResource), b.(*resourcequota.LimitedResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*resourcequota.LimitedResource)(nil), (*LimitedResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_resourcequota_LimitedResource_To_v1alpha1_LimitedResource(a.(*resourcequota.LimitedResource), b.(*LimitedResource), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Configuration_To_resourcequota_Configuration(in *Configuration, out *resourcequota.Configuration, s conversion.Scope) error {
	out.LimitedResources = *(*[]resourcequota.LimitedResource)(unsafe.Pointer(&in.LimitedResources))
	return nil
}

// Convert_v1alpha1_Configuration_To_resourcequota_Configuration is an autogenerated conversion function.
func Convert_v1alpha1_Configuration_To_resourcequota_Configuration(in *Configuration, out *resourcequota.Configuration, s conversion.Scope) error {
	return autoConvert_v1alpha1_Configuration_To_resourcequota_Configuration(in, out, s)
}

func autoConvert_resourcequota_Configuration_To_v1alpha1_Configuration(in *resourcequota.Configuration, out *Configuration, s conversion.Scope) error {
	out.LimitedResources = *(*[]LimitedResource)(unsafe.Pointer(&in.LimitedResources))
	return nil
}

// Convert_resourcequota_Configuration_To_v1alpha1_Configuration is an autogenerated conversion function.
func Convert_resourcequota_Configuration_To_v1alpha1_Configuration(in *resourcequota.Configuration, out *Configuration, s conversion.Scope) error {
	return autoConvert_resourcequota_Configuration_To_v1alpha1_Configuration(in, out, s)
}

func autoConvert_v1alpha1_LimitedResource_To_resourcequota_LimitedResource(in *LimitedResource, out *resourcequota.LimitedResource, s conversion.Scope) error {
	out.APIGroup = in.APIGroup
	out.Resource = in.Resource
	out.MatchContains = *(*[]string)(unsafe.Pointer(&in.MatchContains))
	out.MatchScopes = *(*[]core.ScopedResourceSelectorRequirement)(unsafe.Pointer(&in.MatchScopes))
	return nil
}

// Convert_v1alpha1_LimitedResource_To_resourcequota_LimitedResource is an autogenerated conversion function.
func Convert_v1alpha1_LimitedResource_To_resourcequota_LimitedResource(in *LimitedResource, out *resourcequota.LimitedResource, s conversion.Scope) error {
	return autoConvert_v1alpha1_LimitedResource_To_resourcequota_LimitedResource(in, out, s)
}

func autoConvert_resourcequota_LimitedResource_To_v1alpha1_LimitedResource(in *resourcequota.LimitedResource, out *LimitedResource, s conversion.Scope) error {
	out.APIGroup = in.APIGroup
	out.Resource = in.Resource
	out.MatchContains = *(*[]string)(unsafe.Pointer(&in.MatchContains))
	out.MatchScopes = *(*[]v1.ScopedResourceSelectorRequirement)(unsafe.Pointer(&in.MatchScopes))
	return nil
}

// Convert_resourcequota_LimitedResource_To_v1alpha1_LimitedResource is an autogenerated conversion function.
func Convert_resourcequota_LimitedResource_To_v1alpha1_LimitedResource(in *resourcequota.LimitedResource, out *LimitedResource, s conversion.Scope) error {
	return autoConvert_resourcequota_LimitedResource_To_v1alpha1_LimitedResource(in, out, s)
}
