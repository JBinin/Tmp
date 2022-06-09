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

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "k8s.io/metrics/pkg/apis/metrics/v1alpha1"
	scheme "k8s.io/metrics/pkg/client/clientset/versioned/scheme"
)

// NodeMetricsesGetter has a method to return a NodeMetricsInterface.
// A group's client should implement this interface.
type NodeMetricsesGetter interface {
	NodeMetricses() NodeMetricsInterface
}

// NodeMetricsInterface has methods to work with NodeMetrics resources.
type NodeMetricsInterface interface {
	Get(name string, options v1.GetOptions) (*v1alpha1.NodeMetrics, error)
	List(opts v1.ListOptions) (*v1alpha1.NodeMetricsList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	NodeMetricsExpansion
}

// nodeMetricses implements NodeMetricsInterface
type nodeMetricses struct {
	client rest.Interface
}

// newNodeMetricses returns a NodeMetricses
func newNodeMetricses(c *MetricsV1alpha1Client) *nodeMetricses {
	return &nodeMetricses{
		client: c.RESTClient(),
	}
}

// Get takes name of the nodeMetrics, and returns the corresponding nodeMetrics object, and an error if there is any.
func (c *nodeMetricses) Get(name string, options v1.GetOptions) (result *v1alpha1.NodeMetrics, err error) {
	result = &v1alpha1.NodeMetrics{}
	err = c.client.Get().
		Resource("nodes").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NodeMetricses that match those selectors.
func (c *nodeMetricses) List(opts v1.ListOptions) (result *v1alpha1.NodeMetricsList, err error) {
	result = &v1alpha1.NodeMetricsList{}
	err = c.client.Get().
		Resource("nodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested nodeMetricses.
func (c *nodeMetricses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("nodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
