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

package validating

import (
	"io"

	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/configuration"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/generic"
)

const (
	// Name of admission plug-in
	PluginName = "ValidatingAdmissionWebhook"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(configFile io.Reader) (admission.Interface, error) {
		plugin, err := NewValidatingAdmissionWebhook(configFile)
		if err != nil {
			return nil, err
		}

		return plugin, nil
	})
}

// Plugin is an implementation of admission.Interface.
type Plugin struct {
	*generic.Webhook
}

var _ admission.ValidationInterface = &Plugin{}

// NewValidatingAdmissionWebhook returns a generic admission webhook plugin.
func NewValidatingAdmissionWebhook(configFile io.Reader) (*Plugin, error) {
	handler := admission.NewHandler(admission.Connect, admission.Create, admission.Delete, admission.Update)
	webhook, err := generic.NewWebhook(handler, configFile, configuration.NewValidatingWebhookConfigurationManager, newValidatingDispatcher)
	if err != nil {
		return nil, err
	}
	return &Plugin{webhook}, nil
}

// Validate makes an admission decision based on the request attributes.
func (a *Plugin) Validate(attr admission.Attributes) error {
	return a.Webhook.Dispatch(attr)
}
