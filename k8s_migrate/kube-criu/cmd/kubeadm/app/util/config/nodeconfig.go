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

package config

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1alpha3 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/validation"
	"k8s.io/kubernetes/pkg/util/node"
)

// SetJoinDynamicDefaults checks and sets configuration values for the JoinConfiguration object
func SetJoinDynamicDefaults(cfg *kubeadmapi.JoinConfiguration) error {
	nodeName, err := node.GetHostname(cfg.NodeRegistration.Name)
	if err != nil {
		return err
	}
	cfg.NodeRegistration.Name = nodeName

	if cfg.AdvertiseAddress == "" {
		ip, err := netutil.ChooseBindAddress(nil)
		if err != nil {
			return err
		}
		cfg.AdvertiseAddress = ip.String()
	}

	return nil
}

// NodeConfigFileAndDefaultsToInternalConfig
func NodeConfigFileAndDefaultsToInternalConfig(cfgPath string, defaultversionedcfg *kubeadmapiv1alpha3.JoinConfiguration) (*kubeadmapi.JoinConfiguration, error) {
	internalcfg := &kubeadmapi.JoinConfiguration{}

	if cfgPath != "" {
		// Loads configuration from config file, if provided
		// Nb. --config overrides command line flags, TODO: fix this
		glog.V(1).Infoln("loading configuration from the given file")

		b, err := ioutil.ReadFile(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read config from %q [%v]", cfgPath, err)
		}

		if err := DetectUnsupportedVersion(b); err != nil {
			return nil, err
		}

		if err := runtime.DecodeInto(kubeadmscheme.Codecs.UniversalDecoder(), b, internalcfg); err != nil {
			return nil, err
		}
	} else {
		// Takes passed flags into account; the defaulting is executed once again enforcing assignement of
		// static default values to cfg only for values not provided with flags
		kubeadmscheme.Scheme.Default(defaultversionedcfg)
		kubeadmscheme.Scheme.Convert(defaultversionedcfg, internalcfg, nil)
	}

	// Applies dynamic defaults to settings not provided with flags
	if err := SetJoinDynamicDefaults(internalcfg); err != nil {
		return nil, err
	}
	// Validates cfg (flags/configs + defaults)
	if err := validation.ValidateJoinConfiguration(internalcfg).ToAggregate(); err != nil {
		return nil, err
	}

	return internalcfg, nil
}
