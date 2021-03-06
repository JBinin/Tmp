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

package rollout

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/rest/fake"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
	cmdtesting "k8s.io/kubernetes/pkg/kubectl/cmd/testing"
)

var rolloutPauseGroupVersionEncoder = schema.GroupVersion{Group: "extensions", Version: "v1beta1"}
var rolloutPauseGroupVersionDecoder = schema.GroupVersion{Group: "extensions", Version: runtime.APIVersionInternal}

func TestRolloutPause(t *testing.T) {
	deploymentName := "deployment/nginx-deployment"
	ns := legacyscheme.Codecs
	tf := cmdtesting.NewTestFactory().WithNamespace("test")

	info, _ := runtime.SerializerInfoForMediaType(ns.SupportedMediaTypes(), runtime.ContentTypeJSON)
	encoder := ns.EncoderForVersion(info.Serializer, rolloutPauseGroupVersionEncoder)
	tf.Client = &RolloutPauseRESTClient{
		RESTClient: &fake.RESTClient{
			NegotiatedSerializer: ns,
			Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
				switch p, m := req.URL.Path, req.Method; {
				case p == "/namespaces/test/deployments/nginx-deployment" && (m == "GET" || m == "PATCH"):
					responseDeployment := &extensions.Deployment{}
					responseDeployment.Name = deploymentName
					body := ioutil.NopCloser(bytes.NewReader([]byte(runtime.EncodeOrDie(encoder, responseDeployment))))
					return &http.Response{StatusCode: http.StatusOK, Header: defaultHeader(), Body: body}, nil
				default:
					t.Fatalf("unexpected request: %#v\n%#v", req.URL, req)
					return nil, nil
				}
			}),
		},
	}

	streams, _, buf, _ := genericclioptions.NewTestIOStreams()
	cmd := NewCmdRolloutPause(tf, streams)

	cmd.Run(cmd, []string{deploymentName})
	expectedOutput := "deployment.extensions/" + deploymentName + " paused\n"
	if buf.String() != expectedOutput {
		t.Errorf("expected output: %s, but got: %s", expectedOutput, buf.String())
	}
}

type RolloutPauseRESTClient struct {
	*fake.RESTClient
}

func (c *RolloutPauseRESTClient) Get() *restclient.Request {
	config := restclient.ContentConfig{
		ContentType:          runtime.ContentTypeJSON,
		GroupVersion:         &rolloutPauseGroupVersionEncoder,
		NegotiatedSerializer: c.NegotiatedSerializer,
	}

	info, _ := runtime.SerializerInfoForMediaType(c.NegotiatedSerializer.SupportedMediaTypes(), runtime.ContentTypeJSON)
	serializers := restclient.Serializers{
		Encoder: c.NegotiatedSerializer.EncoderForVersion(info.Serializer, rolloutPauseGroupVersionEncoder),
		Decoder: c.NegotiatedSerializer.DecoderToVersion(info.Serializer, rolloutPauseGroupVersionDecoder),
	}
	if info.StreamSerializer != nil {
		serializers.StreamingSerializer = info.StreamSerializer.Serializer
		serializers.Framer = info.StreamSerializer.Framer
	}
	return restclient.NewRequest(c, "GET", &url.URL{Host: "localhost"}, c.VersionedAPIPath, config, serializers, nil, nil, 0)
}

func (c *RolloutPauseRESTClient) Patch(pt types.PatchType) *restclient.Request {
	config := restclient.ContentConfig{
		ContentType:          runtime.ContentTypeJSON,
		GroupVersion:         &rolloutPauseGroupVersionEncoder,
		NegotiatedSerializer: c.NegotiatedSerializer,
	}

	info, _ := runtime.SerializerInfoForMediaType(c.NegotiatedSerializer.SupportedMediaTypes(), runtime.ContentTypeJSON)
	serializers := restclient.Serializers{
		Encoder: c.NegotiatedSerializer.EncoderForVersion(info.Serializer, rolloutPauseGroupVersionEncoder),
		Decoder: c.NegotiatedSerializer.DecoderToVersion(info.Serializer, rolloutPauseGroupVersionDecoder),
	}
	if info.StreamSerializer != nil {
		serializers.StreamingSerializer = info.StreamSerializer.Serializer
		serializers.Framer = info.StreamSerializer.Framer
	}
	return restclient.NewRequest(c, "PATCH", &url.URL{Host: "localhost"}, c.VersionedAPIPath, config, serializers, nil, nil, 0)
}

func defaultHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", runtime.ContentTypeJSON)
	return header
}
