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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	rest "k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type BatchInterface interface {
	RESTClient() rest.Interface
	CronJobsGetter
	JobsGetter
}

// BatchClient is used to interact with features provided by the batch group.
type BatchClient struct {
	restClient rest.Interface
}

func (c *BatchClient) CronJobs(namespace string) CronJobInterface {
	return newCronJobs(c, namespace)
}

func (c *BatchClient) Jobs(namespace string) JobInterface {
	return newJobs(c, namespace)
}

// NewForConfig creates a new BatchClient for the given config.
func NewForConfig(c *rest.Config) (*BatchClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &BatchClient{client}, nil
}

// NewForConfigOrDie creates a new BatchClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *BatchClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new BatchClient for the given RESTClient.
func New(c rest.Interface) *BatchClient {
	return &BatchClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != scheme.Scheme.PrioritizedVersionsForGroup("batch")[0].Group {
		gv := scheme.Scheme.PrioritizedVersionsForGroup("batch")[0]
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *BatchClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
