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

package framework

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
)

// NewSingleContentTypeSerializer wraps a serializer in a NegotiatedSerializer that handles one content type
func NewSingleContentTypeSerializer(scheme *runtime.Scheme, info runtime.SerializerInfo) runtime.StorageSerializer {
	return &wrappedSerializer{
		scheme: scheme,
		info:   info,
	}
}

type wrappedSerializer struct {
	scheme *runtime.Scheme
	info   runtime.SerializerInfo
}

var _ runtime.StorageSerializer = &wrappedSerializer{}

func (s *wrappedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	return []runtime.SerializerInfo{s.info}
}

func (s *wrappedSerializer) UniversalDeserializer() runtime.Decoder {
	return s.info.Serializer
}

func (s *wrappedSerializer) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return versioning.NewCodec(encoder, nil, s.scheme, s.scheme, s.scheme, s.scheme, gv, nil, s.scheme.Name())
}

func (s *wrappedSerializer) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return versioning.NewCodec(nil, decoder, s.scheme, s.scheme, s.scheme, s.scheme, nil, gv, s.scheme.Name())
}
