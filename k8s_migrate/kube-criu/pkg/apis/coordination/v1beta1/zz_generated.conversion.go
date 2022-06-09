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

package v1beta1

import (
	unsafe "unsafe"

	v1beta1 "k8s.io/api/coordination/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	coordination "k8s.io/kubernetes/pkg/apis/coordination"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.Lease)(nil), (*coordination.Lease)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Lease_To_coordination_Lease(a.(*v1beta1.Lease), b.(*coordination.Lease), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*coordination.Lease)(nil), (*v1beta1.Lease)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_coordination_Lease_To_v1beta1_Lease(a.(*coordination.Lease), b.(*v1beta1.Lease), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.LeaseList)(nil), (*coordination.LeaseList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_LeaseList_To_coordination_LeaseList(a.(*v1beta1.LeaseList), b.(*coordination.LeaseList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*coordination.LeaseList)(nil), (*v1beta1.LeaseList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_coordination_LeaseList_To_v1beta1_LeaseList(a.(*coordination.LeaseList), b.(*v1beta1.LeaseList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.LeaseSpec)(nil), (*coordination.LeaseSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_LeaseSpec_To_coordination_LeaseSpec(a.(*v1beta1.LeaseSpec), b.(*coordination.LeaseSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*coordination.LeaseSpec)(nil), (*v1beta1.LeaseSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_coordination_LeaseSpec_To_v1beta1_LeaseSpec(a.(*coordination.LeaseSpec), b.(*v1beta1.LeaseSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_Lease_To_coordination_Lease(in *v1beta1.Lease, out *coordination.Lease, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_LeaseSpec_To_coordination_LeaseSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_Lease_To_coordination_Lease is an autogenerated conversion function.
func Convert_v1beta1_Lease_To_coordination_Lease(in *v1beta1.Lease, out *coordination.Lease, s conversion.Scope) error {
	return autoConvert_v1beta1_Lease_To_coordination_Lease(in, out, s)
}

func autoConvert_coordination_Lease_To_v1beta1_Lease(in *coordination.Lease, out *v1beta1.Lease, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_coordination_LeaseSpec_To_v1beta1_LeaseSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_coordination_Lease_To_v1beta1_Lease is an autogenerated conversion function.
func Convert_coordination_Lease_To_v1beta1_Lease(in *coordination.Lease, out *v1beta1.Lease, s conversion.Scope) error {
	return autoConvert_coordination_Lease_To_v1beta1_Lease(in, out, s)
}

func autoConvert_v1beta1_LeaseList_To_coordination_LeaseList(in *v1beta1.LeaseList, out *coordination.LeaseList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]coordination.Lease)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1beta1_LeaseList_To_coordination_LeaseList is an autogenerated conversion function.
func Convert_v1beta1_LeaseList_To_coordination_LeaseList(in *v1beta1.LeaseList, out *coordination.LeaseList, s conversion.Scope) error {
	return autoConvert_v1beta1_LeaseList_To_coordination_LeaseList(in, out, s)
}

func autoConvert_coordination_LeaseList_To_v1beta1_LeaseList(in *coordination.LeaseList, out *v1beta1.LeaseList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1beta1.Lease)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_coordination_LeaseList_To_v1beta1_LeaseList is an autogenerated conversion function.
func Convert_coordination_LeaseList_To_v1beta1_LeaseList(in *coordination.LeaseList, out *v1beta1.LeaseList, s conversion.Scope) error {
	return autoConvert_coordination_LeaseList_To_v1beta1_LeaseList(in, out, s)
}

func autoConvert_v1beta1_LeaseSpec_To_coordination_LeaseSpec(in *v1beta1.LeaseSpec, out *coordination.LeaseSpec, s conversion.Scope) error {
	out.HolderIdentity = (*string)(unsafe.Pointer(in.HolderIdentity))
	out.LeaseDurationSeconds = (*int32)(unsafe.Pointer(in.LeaseDurationSeconds))
	out.AcquireTime = (*v1.MicroTime)(unsafe.Pointer(in.AcquireTime))
	out.RenewTime = (*v1.MicroTime)(unsafe.Pointer(in.RenewTime))
	out.LeaseTransitions = (*int32)(unsafe.Pointer(in.LeaseTransitions))
	return nil
}

// Convert_v1beta1_LeaseSpec_To_coordination_LeaseSpec is an autogenerated conversion function.
func Convert_v1beta1_LeaseSpec_To_coordination_LeaseSpec(in *v1beta1.LeaseSpec, out *coordination.LeaseSpec, s conversion.Scope) error {
	return autoConvert_v1beta1_LeaseSpec_To_coordination_LeaseSpec(in, out, s)
}

func autoConvert_coordination_LeaseSpec_To_v1beta1_LeaseSpec(in *coordination.LeaseSpec, out *v1beta1.LeaseSpec, s conversion.Scope) error {
	out.HolderIdentity = (*string)(unsafe.Pointer(in.HolderIdentity))
	out.LeaseDurationSeconds = (*int32)(unsafe.Pointer(in.LeaseDurationSeconds))
	out.AcquireTime = (*v1.MicroTime)(unsafe.Pointer(in.AcquireTime))
	out.RenewTime = (*v1.MicroTime)(unsafe.Pointer(in.RenewTime))
	out.LeaseTransitions = (*int32)(unsafe.Pointer(in.LeaseTransitions))
	return nil
}

// Convert_coordination_LeaseSpec_To_v1beta1_LeaseSpec is an autogenerated conversion function.
func Convert_coordination_LeaseSpec_To_v1beta1_LeaseSpec(in *coordination.LeaseSpec, out *v1beta1.LeaseSpec, s conversion.Scope) error {
	return autoConvert_coordination_LeaseSpec_To_v1beta1_LeaseSpec(in, out, s)
}
