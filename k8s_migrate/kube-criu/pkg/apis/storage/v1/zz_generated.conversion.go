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

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/kubernetes/pkg/apis/core"
	storage "k8s.io/kubernetes/pkg/apis/storage"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.StorageClass)(nil), (*storage.StorageClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StorageClass_To_storage_StorageClass(a.(*v1.StorageClass), b.(*storage.StorageClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.StorageClass)(nil), (*v1.StorageClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_StorageClass_To_v1_StorageClass(a.(*storage.StorageClass), b.(*v1.StorageClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StorageClassList)(nil), (*storage.StorageClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StorageClassList_To_storage_StorageClassList(a.(*v1.StorageClassList), b.(*storage.StorageClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.StorageClassList)(nil), (*v1.StorageClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_StorageClassList_To_v1_StorageClassList(a.(*storage.StorageClassList), b.(*v1.StorageClassList), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_StorageClass_To_storage_StorageClass(in *v1.StorageClass, out *storage.StorageClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Provisioner = in.Provisioner
	out.Parameters = *(*map[string]string)(unsafe.Pointer(&in.Parameters))
	out.ReclaimPolicy = (*core.PersistentVolumeReclaimPolicy)(unsafe.Pointer(in.ReclaimPolicy))
	out.MountOptions = *(*[]string)(unsafe.Pointer(&in.MountOptions))
	out.AllowVolumeExpansion = (*bool)(unsafe.Pointer(in.AllowVolumeExpansion))
	out.VolumeBindingMode = (*storage.VolumeBindingMode)(unsafe.Pointer(in.VolumeBindingMode))
	out.AllowedTopologies = *(*[]core.TopologySelectorTerm)(unsafe.Pointer(&in.AllowedTopologies))
	return nil
}

// Convert_v1_StorageClass_To_storage_StorageClass is an autogenerated conversion function.
func Convert_v1_StorageClass_To_storage_StorageClass(in *v1.StorageClass, out *storage.StorageClass, s conversion.Scope) error {
	return autoConvert_v1_StorageClass_To_storage_StorageClass(in, out, s)
}

func autoConvert_storage_StorageClass_To_v1_StorageClass(in *storage.StorageClass, out *v1.StorageClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Provisioner = in.Provisioner
	out.Parameters = *(*map[string]string)(unsafe.Pointer(&in.Parameters))
	out.ReclaimPolicy = (*corev1.PersistentVolumeReclaimPolicy)(unsafe.Pointer(in.ReclaimPolicy))
	out.MountOptions = *(*[]string)(unsafe.Pointer(&in.MountOptions))
	out.AllowVolumeExpansion = (*bool)(unsafe.Pointer(in.AllowVolumeExpansion))
	out.VolumeBindingMode = (*v1.VolumeBindingMode)(unsafe.Pointer(in.VolumeBindingMode))
	out.AllowedTopologies = *(*[]corev1.TopologySelectorTerm)(unsafe.Pointer(&in.AllowedTopologies))
	return nil
}

// Convert_storage_StorageClass_To_v1_StorageClass is an autogenerated conversion function.
func Convert_storage_StorageClass_To_v1_StorageClass(in *storage.StorageClass, out *v1.StorageClass, s conversion.Scope) error {
	return autoConvert_storage_StorageClass_To_v1_StorageClass(in, out, s)
}

func autoConvert_v1_StorageClassList_To_storage_StorageClassList(in *v1.StorageClassList, out *storage.StorageClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]storage.StorageClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1_StorageClassList_To_storage_StorageClassList is an autogenerated conversion function.
func Convert_v1_StorageClassList_To_storage_StorageClassList(in *v1.StorageClassList, out *storage.StorageClassList, s conversion.Scope) error {
	return autoConvert_v1_StorageClassList_To_storage_StorageClassList(in, out, s)
}

func autoConvert_storage_StorageClassList_To_v1_StorageClassList(in *storage.StorageClassList, out *v1.StorageClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1.StorageClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_storage_StorageClassList_To_v1_StorageClassList is an autogenerated conversion function.
func Convert_storage_StorageClassList_To_v1_StorageClassList(in *storage.StorageClassList, out *v1.StorageClassList, s conversion.Scope) error {
	return autoConvert_storage_StorageClassList_To_v1_StorageClassList(in, out, s)
}
