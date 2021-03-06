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

// Package mutating delegates admission checks to dynamically configured
// mutating webhooks.
package mutating

import (
	"context"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/api/admissionregistration/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	admissionmetrics "k8s.io/apiserver/pkg/admission/metrics"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/config"
	webhookerrors "k8s.io/apiserver/pkg/admission/plugin/webhook/errors"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/generic"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/request"
)

type mutatingDispatcher struct {
	cm     *config.ClientManager
	plugin *Plugin
}

func newMutatingDispatcher(p *Plugin) func(cm *config.ClientManager) generic.Dispatcher {
	return func(cm *config.ClientManager) generic.Dispatcher {
		return &mutatingDispatcher{cm, p}
	}
}

var _ generic.Dispatcher = &mutatingDispatcher{}

func (a *mutatingDispatcher) Dispatch(ctx context.Context, attr *generic.VersionedAttributes, relevantHooks []*v1beta1.Webhook) error {
	for _, hook := range relevantHooks {
		t := time.Now()
		err := a.callAttrMutatingHook(ctx, hook, attr)
		admissionmetrics.Metrics.ObserveWebhook(time.Since(t), err != nil, attr.Attributes, "admit", hook.Name)
		if err == nil {
			continue
		}

		ignoreClientCallFailures := hook.FailurePolicy != nil && *hook.FailurePolicy == v1beta1.Ignore
		if callErr, ok := err.(*webhookerrors.ErrCallingWebhook); ok {
			if ignoreClientCallFailures {
				glog.Warningf("Failed calling webhook, failing open %v: %v", hook.Name, callErr)
				utilruntime.HandleError(callErr)
				continue
			}
			glog.Warningf("Failed calling webhook, failing closed %v: %v", hook.Name, err)
		}
		return apierrors.NewInternalError(err)
	}

	// convert attr.VersionedObject to the internal version in the underlying admission.Attributes
	if attr.VersionedObject != nil {
		return a.plugin.scheme.Convert(attr.VersionedObject, attr.Attributes.GetObject(), nil)
	}
	return nil
}

// note that callAttrMutatingHook updates attr
func (a *mutatingDispatcher) callAttrMutatingHook(ctx context.Context, h *v1beta1.Webhook, attr *generic.VersionedAttributes) error {
	if attr.IsDryRun() {
		// TODO: support this
		return webhookerrors.NewDryRunUnsupportedErr(h.Name)
	}

	// Make the webhook request
	request := request.CreateAdmissionReview(attr)
	client, err := a.cm.HookClient(h)
	if err != nil {
		return &webhookerrors.ErrCallingWebhook{WebhookName: h.Name, Reason: err}
	}
	response := &admissionv1beta1.AdmissionReview{}
	if err := client.Post().Context(ctx).Body(&request).Do().Into(response); err != nil {
		return &webhookerrors.ErrCallingWebhook{WebhookName: h.Name, Reason: err}
	}

	if response.Response == nil {
		return &webhookerrors.ErrCallingWebhook{WebhookName: h.Name, Reason: fmt.Errorf("Webhook response was absent")}
	}

	if !response.Response.Allowed {
		return webhookerrors.ToStatusErr(h.Name, response.Response.Result)
	}

	patchJS := response.Response.Patch
	if len(patchJS) == 0 {
		return nil
	}
	patchObj, err := jsonpatch.DecodePatch(patchJS)
	if err != nil {
		return apierrors.NewInternalError(err)
	}
	if len(patchObj) == 0 {
		return nil
	}

	// if a non-empty patch was provided, and we have no object we can apply it to (e.g. a DELETE admission operation), error
	if attr.VersionedObject == nil {
		return apierrors.NewInternalError(fmt.Errorf("admission webhook %q attempted to modify the object, which is not supported for this operation", h.Name))
	}

	objJS, err := runtime.Encode(a.plugin.jsonSerializer, attr.VersionedObject)
	if err != nil {
		return apierrors.NewInternalError(err)
	}
	patchedJS, err := patchObj.Apply(objJS)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	var newVersionedObject runtime.Object
	if _, ok := attr.VersionedObject.(*unstructured.Unstructured); ok {
		// Custom Resources don't have corresponding Go struct's.
		// They are represented as Unstructured.
		newVersionedObject = &unstructured.Unstructured{}
	} else {
		newVersionedObject, err = a.plugin.scheme.New(attr.GetKind())
		if err != nil {
			return apierrors.NewInternalError(err)
		}
	}
	// TODO: if we have multiple mutating webhooks, we can remember the json
	// instead of encoding and decoding for each one.
	if _, _, err := a.plugin.jsonSerializer.Decode(patchedJS, nil, newVersionedObject); err != nil {
		return apierrors.NewInternalError(err)
	}
	attr.VersionedObject = newVersionedObject
	a.plugin.scheme.Default(attr.VersionedObject)
	return nil
}
