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
Copyright 2014 The Kubernetes Authors.

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

package resource

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest/fake"

	// TODO we need to remove this linkage and create our own scheme
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func objBody(obj runtime.Object) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(runtime.EncodeOrDie(corev1Codec, obj))))
}

func header() http.Header {
	header := http.Header{}
	header.Set("Content-Type", runtime.ContentTypeJSON)
	return header
}

// splitPath returns the segments for a URL path.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

// V1DeepEqualSafePodSpec returns a PodSpec which is ready to be used with apiequality.Semantic.DeepEqual
func V1DeepEqualSafePodSpec() corev1.PodSpec {
	grace := int64(30)
	return corev1.PodSpec{
		RestartPolicy:                 corev1.RestartPolicyAlways,
		DNSPolicy:                     corev1.DNSClusterFirst,
		TerminationGracePeriodSeconds: &grace,
		SecurityContext:               &corev1.PodSecurityContext{},
	}
}

func TestHelperDelete(t *testing.T) {
	tests := []struct {
		name    string
		Err     bool
		Req     func(*http.Request) bool
		Resp    *http.Response
		HttpErr error
	}{
		{
			name:    "test1",
			HttpErr: errors.New("failure"),
			Err:     true,
		},
		{
			name: "test2",
			Resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusFailure}),
			},
			Err: true,
		},
		{
			name: "test3pkg/kubectl/genericclioptions/resource/helper_test.go",
			Resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusSuccess}),
			},
			Req: func(req *http.Request) bool {
				if req.Method != "DELETE" {
					t.Errorf("unexpected method: %#v", req)
					return false
				}
				parts := splitPath(req.URL.Path)
				if len(parts) < 3 {
					t.Errorf("expected URL path to have 3 parts: %s", req.URL.Path)
					return false
				}
				if parts[1] != "bar" {
					t.Errorf("url doesn't contain namespace: %#v", req)
					return false
				}
				if parts[2] != "foo" {
					t.Errorf("url doesn't contain name: %#v", req)
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fake.RESTClient{
				NegotiatedSerializer: scheme.Codecs,
				Resp:                 tt.Resp,
				Err:                  tt.HttpErr,
			}
			modifier := &Helper{
				RESTClient:      client,
				NamespaceScoped: true,
			}
			_, err := modifier.Delete("bar", "foo")
			if (err != nil) != tt.Err {
				t.Errorf("unexpected error: %t %v", tt.Err, err)
			}
			if err != nil {
				return
			}
			if tt.Req != nil && !tt.Req(client.Req) {
				t.Errorf("unexpected request: %#v", client.Req)
			}
		})
	}
}

func TestHelperCreate(t *testing.T) {
	expectPost := func(req *http.Request) bool {
		if req.Method != "POST" {
			t.Errorf("unexpected method: %#v", req)
			return false
		}
		parts := splitPath(req.URL.Path)
		if parts[1] != "bar" {
			t.Errorf("url doesn't contain namespace: %#v", req)
			return false
		}
		return true
	}

	tests := []struct {
		name    string
		Resp    *http.Response
		HttpErr error
		Modify  bool
		Object  runtime.Object

		ExpectObject runtime.Object
		Err          bool
		Req          func(*http.Request) bool
	}{
		{
			name:    "test1",
			HttpErr: errors.New("failure"),
			Err:     true,
		},
		{
			name: "test1",
			Resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusFailure}),
			},
			Err: true,
		},
		{
			name: "test1",
			Resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusSuccess}),
			},
			Object:       &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			ExpectObject: &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			Req:          expectPost,
		},
		{
			name:         "test1",
			Modify:       false,
			Object:       &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"}},
			ExpectObject: &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"}},
			Resp:         &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&metav1.Status{Status: metav1.StatusSuccess})},
			Req:          expectPost,
		},
		{
			name:   "test1",
			Modify: true,
			Object: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"},
				Spec:       V1DeepEqualSafePodSpec(),
			},
			ExpectObject: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec:       V1DeepEqualSafePodSpec(),
			},
			Resp: &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&metav1.Status{Status: metav1.StatusSuccess})},
			Req:  expectPost,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fake.RESTClient{
				GroupVersion:         corev1GV,
				NegotiatedSerializer: scheme.Codecs,
				Resp:                 tt.Resp,
				Err:                  tt.HttpErr,
			}
			modifier := &Helper{
				RESTClient:      client,
				NamespaceScoped: true,
			}
			_, err := modifier.Create("bar", tt.Modify, tt.Object)
			if (err != nil) != tt.Err {
				t.Errorf("%d: unexpected error: %t %v", i, tt.Err, err)
			}
			if err != nil {
				return
			}
			if tt.Req != nil && !tt.Req(client.Req) {
				t.Errorf("%d: unexpected request: %#v", i, client.Req)
			}
			body, err := ioutil.ReadAll(client.Req.Body)
			if err != nil {
				t.Fatalf("%d: unexpected error: %#v", i, err)
			}
			t.Logf("got body: %s", string(body))
			expect := []byte{}
			if tt.ExpectObject != nil {
				expect = []byte(runtime.EncodeOrDie(corev1Codec, tt.ExpectObject))
			}
			if !reflect.DeepEqual(expect, body) {
				t.Errorf("%d: unexpected body: %s (expected %s)", i, string(body), string(expect))
			}
		})
	}
}

func TestHelperGet(t *testing.T) {
	tests := []struct {
		name    string
		Err     bool
		Req     func(*http.Request) bool
		Resp    *http.Response
		HttpErr error
	}{
		{
			name:    "test1",
			HttpErr: errors.New("failure"),
			Err:     true,
		},
		{
			name: "test1",
			Resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusFailure}),
			},
			Err: true,
		},
		{
			name: "test1",
			Resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     header(),
				Body:       objBody(&corev1.Pod{TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"}, ObjectMeta: metav1.ObjectMeta{Name: "foo"}}),
			},
			Req: func(req *http.Request) bool {
				if req.Method != "GET" {
					t.Errorf("unexpected method: %#v", req)
					return false
				}
				parts := splitPath(req.URL.Path)
				if parts[1] != "bar" {
					t.Errorf("url doesn't contain namespace: %#v", req)
					return false
				}
				if parts[2] != "foo" {
					t.Errorf("url doesn't contain name: %#v", req)
					return false
				}
				return true
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fake.RESTClient{
				GroupVersion:         corev1GV,
				NegotiatedSerializer: serializer.DirectCodecFactory{CodecFactory: scheme.Codecs},
				Resp:                 tt.Resp,
				Err:                  tt.HttpErr,
			}
			modifier := &Helper{
				RESTClient:      client,
				NamespaceScoped: true,
			}
			obj, err := modifier.Get("bar", "foo", false)

			if (err != nil) != tt.Err {
				t.Errorf("unexpected error: %d %t %v", i, tt.Err, err)
			}
			if err != nil {
				return
			}
			if obj.(*corev1.Pod).Name != "foo" {
				t.Errorf("unexpected object: %#v", obj)
			}
			if tt.Req != nil && !tt.Req(client.Req) {
				t.Errorf("unexpected request: %#v", client.Req)
			}
		})
	}
}

func TestHelperList(t *testing.T) {
	tests := []struct {
		name    string
		Err     bool
		Req     func(*http.Request) bool
		Resp    *http.Response
		HttpErr error
	}{
		{
			name:    "test1",
			HttpErr: errors.New("failure"),
			Err:     true,
		},
		{
			name: "test2",
			Resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusFailure}),
			},
			Err: true,
		},
		{
			name: "test3",
			Resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     header(),
				Body: objBody(&corev1.PodList{
					Items: []corev1.Pod{{
						ObjectMeta: metav1.ObjectMeta{Name: "foo"},
					},
					},
				}),
			},
			Req: func(req *http.Request) bool {
				if req.Method != "GET" {
					t.Errorf("unexpected method: %#v", req)
					return false
				}
				if req.URL.Path != "/namespaces/bar" {
					t.Errorf("url doesn't contain name: %#v", req.URL)
					return false
				}
				if req.URL.Query().Get(metav1.LabelSelectorQueryParam(corev1GV.String())) != labels.SelectorFromSet(labels.Set{"foo": "baz"}).String() {
					t.Errorf("url doesn't contain query parameters: %#v", req.URL)
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fake.RESTClient{
				GroupVersion:         corev1GV,
				NegotiatedSerializer: serializer.DirectCodecFactory{CodecFactory: scheme.Codecs},
				Resp:                 tt.Resp,
				Err:                  tt.HttpErr,
			}
			modifier := &Helper{
				RESTClient:      client,
				NamespaceScoped: true,
			}
			obj, err := modifier.List("bar", corev1GV.String(), false, &metav1.ListOptions{LabelSelector: "foo=baz"})
			if (err != nil) != tt.Err {
				t.Errorf("unexpected error: %t %v", tt.Err, err)
			}
			if err != nil {
				return
			}
			if obj.(*corev1.PodList).Items[0].Name != "foo" {
				t.Errorf("unexpected object: %#v", obj)
			}
			if tt.Req != nil && !tt.Req(client.Req) {
				t.Errorf("unexpected request: %#v", client.Req)
			}
		})
	}
}

func TestHelperListSelectorCombination(t *testing.T) {
	tests := []struct {
		Name          string
		Err           bool
		ErrMsg        string
		FieldSelector string
		LabelSelector string
	}{
		{
			Name: "No selector",
			Err:  false,
		},
		{
			Name:          "Only Label Selector",
			Err:           false,
			LabelSelector: "foo=baz",
		},
		{
			Name:          "Only Field Selector",
			Err:           false,
			FieldSelector: "xyz=zyx",
		},
		{
			Name:          "Both Label and Field Selector",
			Err:           false,
			LabelSelector: "foo=baz",
			FieldSelector: "xyz=zyx",
		},
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     header(),
		Body: objBody(&corev1.PodList{
			Items: []corev1.Pod{{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			},
			},
		}),
	}
	client := &fake.RESTClient{
		NegotiatedSerializer: scheme.Codecs,
		Resp:                 resp,
		Err:                  nil,
	}
	modifier := &Helper{
		RESTClient:      client,
		NamespaceScoped: true,
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := modifier.List("bar",
				corev1GV.String(),
				false,
				&metav1.ListOptions{LabelSelector: tt.LabelSelector, FieldSelector: tt.FieldSelector})
			if tt.Err {
				if err == nil {
					t.Errorf("%q expected error: %q", tt.Name, tt.ErrMsg)
				}
				if err != nil && err.Error() != tt.ErrMsg {
					t.Errorf("%q expected error: %q", tt.Name, tt.ErrMsg)
				}
			}
		})
	}
}

func TestHelperReplace(t *testing.T) {
	expectPut := func(path string, req *http.Request) bool {
		if req.Method != "PUT" {
			t.Errorf("unexpected method: %#v", req)
			return false
		}
		if req.URL.Path != path {
			t.Errorf("unexpected url: %v", req.URL)
			return false
		}
		return true
	}

	tests := []struct {
		Name            string
		Resp            *http.Response
		HTTPClient      *http.Client
		HttpErr         error
		Overwrite       bool
		Object          runtime.Object
		Namespace       string
		NamespaceScoped bool

		ExpectPath   string
		ExpectObject runtime.Object
		Err          bool
		Req          func(string, *http.Request) bool
	}{
		{
			Name:            "test1",
			Namespace:       "bar",
			NamespaceScoped: true,
			HttpErr:         errors.New("failure"),
			Err:             true,
		},
		{
			Name:            "test2",
			Namespace:       "bar",
			NamespaceScoped: true,
			Object:          &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			Resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusFailure}),
			},
			Err: true,
		},
		{
			Name:            "test3",
			Namespace:       "bar",
			NamespaceScoped: true,
			Object:          &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			ExpectPath:      "/namespaces/bar/foo",
			ExpectObject:    &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			Resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     header(),
				Body:       objBody(&metav1.Status{Status: metav1.StatusSuccess}),
			},
			Req: expectPut,
		},
		// namespace scoped resource
		{
			Name:            "test4",
			Namespace:       "bar",
			NamespaceScoped: true,
			Object: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec:       V1DeepEqualSafePodSpec(),
			},
			ExpectPath: "/namespaces/bar/foo",
			ExpectObject: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"},
				Spec:       V1DeepEqualSafePodSpec(),
			},
			Overwrite: true,
			HTTPClient: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
				if req.Method == "PUT" {
					return &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&metav1.Status{Status: metav1.StatusSuccess})}, nil
				}
				return &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"}})}, nil
			}),
			Req: expectPut,
		},
		// cluster scoped resource
		{
			Name: "test5",
			Object: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			},
			ExpectObject: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"},
			},
			Overwrite:  true,
			ExpectPath: "/foo",
			HTTPClient: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
				if req.Method == "PUT" {
					return &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&metav1.Status{Status: metav1.StatusSuccess})}, nil
				}
				return &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"}})}, nil
			}),
			Req: expectPut,
		},
		{
			Name:            "test6",
			Namespace:       "bar",
			NamespaceScoped: true,
			Object:          &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"}},
			ExpectPath:      "/namespaces/bar/foo",
			ExpectObject:    &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"}},
			Resp:            &http.Response{StatusCode: http.StatusOK, Header: header(), Body: objBody(&metav1.Status{Status: metav1.StatusSuccess})},
			Req:             expectPut,
		},
	}
	for i, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			client := &fake.RESTClient{
				GroupVersion:         corev1GV,
				NegotiatedSerializer: serializer.DirectCodecFactory{CodecFactory: scheme.Codecs},
				Client:               tt.HTTPClient,
				Resp:                 tt.Resp,
				Err:                  tt.HttpErr,
			}
			modifier := &Helper{
				RESTClient:      client,
				NamespaceScoped: tt.NamespaceScoped,
			}
			_, err := modifier.Replace(tt.Namespace, "foo", tt.Overwrite, tt.Object)
			if (err != nil) != tt.Err {
				t.Errorf("%d: unexpected error: %t %v", i, tt.Err, err)
			}
			if err != nil {
				return
			}
			if tt.Req != nil && !tt.Req(tt.ExpectPath, client.Req) {
				t.Errorf("%d: unexpected request: %#v", i, client.Req)
			}
			body, err := ioutil.ReadAll(client.Req.Body)
			if err != nil {
				t.Fatalf("%d: unexpected error: %#v", i, err)
			}
			expect := []byte{}
			if tt.ExpectObject != nil {
				expect = []byte(runtime.EncodeOrDie(corev1Codec, tt.ExpectObject))
			}
			if !reflect.DeepEqual(expect, body) {
				t.Errorf("%d: unexpected body: %s", i, string(body))
			}
		})
	}
}
