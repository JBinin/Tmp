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
Copyright 2018 The Kubernetes Authors.

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

package scheme

import (
	"testing"

	"k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
)

func TestCronJob(t *testing.T) {
	src := &v1beta1.CronJob{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}

	encoder := Codecs.LegacyCodec(v1.SchemeGroupVersion, v1beta1.SchemeGroupVersion)
	cronjobBytes, err := runtime.Encode(encoder, src)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cronjobBytes))
	t.Log(Scheme.PrioritizedVersionsAllGroups())

	decoder := Codecs.UniversalDecoder(Scheme.PrioritizedVersionsAllGroups()...)

	uncastDst, err := runtime.Decode(decoder, cronjobBytes)
	if err != nil {
		t.Fatal(err)
	}

	// clear typemeta
	uncastDst.(*v1beta1.CronJob).TypeMeta = metav1.TypeMeta{}

	if !equality.Semantic.DeepEqual(src, uncastDst) {
		t.Fatal(diff.ObjectDiff(src, uncastDst))
	}
}
