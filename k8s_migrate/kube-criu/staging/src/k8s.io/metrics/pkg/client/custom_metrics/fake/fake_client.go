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
Copyright 2017 The Kubernetes Authors.

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

package fake

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/testing"
	"k8s.io/metrics/pkg/apis/custom_metrics/v1beta1"
	cmclient "k8s.io/metrics/pkg/client/custom_metrics"
)

type GetForActionImpl struct {
	testing.GetAction
	MetricName    string
	LabelSelector labels.Selector
}

type GetForAction interface {
	testing.GetAction
	GetMetricName() string
	GetLabelSelector() labels.Selector
}

func (i GetForActionImpl) GetMetricName() string {
	return i.MetricName
}

func (i GetForActionImpl) GetLabelSelector() labels.Selector {
	return i.LabelSelector
}

func (i GetForActionImpl) GetSubresource() string {
	return i.MetricName
}

func (i GetForActionImpl) DeepCopy() testing.Action {
	var labelSelector labels.Selector
	if i.LabelSelector != nil {
		labelSelector = i.LabelSelector.DeepCopySelector()
	}
	return GetForActionImpl{
		GetAction:     i.GetAction.DeepCopy().(testing.GetAction),
		MetricName:    i.MetricName,
		LabelSelector: labelSelector,
	}
}

func NewGetForAction(groupKind schema.GroupKind, namespace, name string, metricName string, labelSelector labels.Selector) GetForActionImpl {
	// the version doesn't matter
	gvk := groupKind.WithVersion("")
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	groupResourceForKind := schema.GroupResource{
		Group:    gvr.Group,
		Resource: gvr.Resource,
	}
	resource := schema.GroupResource{
		Group:    v1beta1.SchemeGroupVersion.Group,
		Resource: groupResourceForKind.String(),
	}
	return GetForActionImpl{
		GetAction:     testing.NewGetAction(resource.WithVersion(""), namespace, name),
		MetricName:    metricName,
		LabelSelector: labelSelector,
	}
}

func NewRootGetForAction(groupKind schema.GroupKind, name string, metricName string, labelSelector labels.Selector) GetForActionImpl {
	// the version doesn't matter
	gvk := groupKind.WithVersion("")
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	groupResourceForKind := schema.GroupResource{
		Group:    gvr.Group,
		Resource: gvr.Resource,
	}
	resource := schema.GroupResource{
		Group:    v1beta1.SchemeGroupVersion.Group,
		Resource: groupResourceForKind.String(),
	}
	return GetForActionImpl{
		GetAction:     testing.NewRootGetAction(resource.WithVersion(""), name),
		MetricName:    metricName,
		LabelSelector: labelSelector,
	}
}

type FakeCustomMetricsClient struct {
	testing.Fake
}

func (c *FakeCustomMetricsClient) RootScopedMetrics() cmclient.MetricsInterface {
	return &fakeRootScopedMetrics{
		Fake: c,
	}
}

func (c *FakeCustomMetricsClient) NamespacedMetrics(namespace string) cmclient.MetricsInterface {
	return &fakeNamespacedMetrics{
		Fake: c,
		ns:   namespace,
	}
}

type fakeNamespacedMetrics struct {
	Fake *FakeCustomMetricsClient
	ns   string
}

func (m *fakeNamespacedMetrics) GetForObject(groupKind schema.GroupKind, name string, metricName string) (*v1beta1.MetricValue, error) {
	obj, err := m.Fake.
		Invokes(NewGetForAction(groupKind, m.ns, name, metricName, nil), &v1beta1.MetricValueList{})

	if obj == nil {
		return nil, err
	}

	objList := obj.(*v1beta1.MetricValueList)
	if len(objList.Items) != 1 {
		return nil, fmt.Errorf("the custom metrics API server returned %v results when we asked for exactly one", len(objList.Items))
	}

	return &objList.Items[0], err
}

func (m *fakeNamespacedMetrics) GetForObjects(groupKind schema.GroupKind, selector labels.Selector, metricName string) (*v1beta1.MetricValueList, error) {
	obj, err := m.Fake.
		Invokes(NewGetForAction(groupKind, m.ns, "*", metricName, selector), &v1beta1.MetricValueList{})

	if obj == nil {
		return nil, err
	}

	return obj.(*v1beta1.MetricValueList), err
}

type fakeRootScopedMetrics struct {
	Fake *FakeCustomMetricsClient
}

func (m *fakeRootScopedMetrics) GetForObject(groupKind schema.GroupKind, name string, metricName string) (*v1beta1.MetricValue, error) {
	obj, err := m.Fake.
		Invokes(NewRootGetForAction(groupKind, name, metricName, nil), &v1beta1.MetricValueList{})

	if obj == nil {
		return nil, err
	}

	objList := obj.(*v1beta1.MetricValueList)
	if len(objList.Items) != 1 {
		return nil, fmt.Errorf("the custom metrics API server returned %v results when we asked for exactly one", len(objList.Items))
	}

	return &objList.Items[0], err
}

func (m *fakeRootScopedMetrics) GetForObjects(groupKind schema.GroupKind, selector labels.Selector, metricName string) (*v1beta1.MetricValueList, error) {
	obj, err := m.Fake.
		Invokes(NewRootGetForAction(groupKind, "*", metricName, selector), &v1beta1.MetricValueList{})

	if obj == nil {
		return nil, err
	}

	return obj.(*v1beta1.MetricValueList), err
}
