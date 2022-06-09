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

package autoscalingv1

import (
	v1 "k8s.io/api/autoscaling/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	scheme "k8s.io/client-go/scale/scheme"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.Scale)(nil), (*scheme.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Scale_To_scheme_Scale(a.(*v1.Scale), b.(*scheme.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*scheme.Scale)(nil), (*v1.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_scheme_Scale_To_v1_Scale(a.(*scheme.Scale), b.(*v1.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScaleSpec)(nil), (*scheme.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleSpec_To_scheme_ScaleSpec(a.(*v1.ScaleSpec), b.(*scheme.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*scheme.ScaleSpec)(nil), (*v1.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_scheme_ScaleSpec_To_v1_ScaleSpec(a.(*scheme.ScaleSpec), b.(*v1.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScaleStatus)(nil), (*scheme.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleStatus_To_scheme_ScaleStatus(a.(*v1.ScaleStatus), b.(*scheme.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*scheme.ScaleStatus)(nil), (*v1.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_scheme_ScaleStatus_To_v1_ScaleStatus(a.(*scheme.ScaleStatus), b.(*v1.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*scheme.ScaleStatus)(nil), (*v1.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_scheme_ScaleStatus_To_v1_ScaleStatus(a.(*scheme.ScaleStatus), b.(*v1.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ScaleStatus)(nil), (*scheme.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleStatus_To_scheme_ScaleStatus(a.(*v1.ScaleStatus), b.(*scheme.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_Scale_To_scheme_Scale(in *v1.Scale, out *scheme.Scale, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_ScaleSpec_To_scheme_ScaleSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_ScaleStatus_To_scheme_ScaleStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_Scale_To_scheme_Scale is an autogenerated conversion function.
func Convert_v1_Scale_To_scheme_Scale(in *v1.Scale, out *scheme.Scale, s conversion.Scope) error {
	return autoConvert_v1_Scale_To_scheme_Scale(in, out, s)
}

func autoConvert_scheme_Scale_To_v1_Scale(in *scheme.Scale, out *v1.Scale, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_scheme_ScaleSpec_To_v1_ScaleSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_scheme_ScaleStatus_To_v1_ScaleStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_scheme_Scale_To_v1_Scale is an autogenerated conversion function.
func Convert_scheme_Scale_To_v1_Scale(in *scheme.Scale, out *v1.Scale, s conversion.Scope) error {
	return autoConvert_scheme_Scale_To_v1_Scale(in, out, s)
}

func autoConvert_v1_ScaleSpec_To_scheme_ScaleSpec(in *v1.ScaleSpec, out *scheme.ScaleSpec, s conversion.Scope) error {
	out.Replicas = in.Replicas
	return nil
}

// Convert_v1_ScaleSpec_To_scheme_ScaleSpec is an autogenerated conversion function.
func Convert_v1_ScaleSpec_To_scheme_ScaleSpec(in *v1.ScaleSpec, out *scheme.ScaleSpec, s conversion.Scope) error {
	return autoConvert_v1_ScaleSpec_To_scheme_ScaleSpec(in, out, s)
}

func autoConvert_scheme_ScaleSpec_To_v1_ScaleSpec(in *scheme.ScaleSpec, out *v1.ScaleSpec, s conversion.Scope) error {
	out.Replicas = in.Replicas
	return nil
}

// Convert_scheme_ScaleSpec_To_v1_ScaleSpec is an autogenerated conversion function.
func Convert_scheme_ScaleSpec_To_v1_ScaleSpec(in *scheme.ScaleSpec, out *v1.ScaleSpec, s conversion.Scope) error {
	return autoConvert_scheme_ScaleSpec_To_v1_ScaleSpec(in, out, s)
}

func autoConvert_v1_ScaleStatus_To_scheme_ScaleStatus(in *v1.ScaleStatus, out *scheme.ScaleStatus, s conversion.Scope) error {
	out.Replicas = in.Replicas
	// WARNING: in.Selector requires manual conversion: inconvertible types (string vs *k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector)
	return nil
}

func autoConvert_scheme_ScaleStatus_To_v1_ScaleStatus(in *scheme.ScaleStatus, out *v1.ScaleStatus, s conversion.Scope) error {
	out.Replicas = in.Replicas
	// WARNING: in.Selector requires manual conversion: inconvertible types (*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector vs string)
	return nil
}
