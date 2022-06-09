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

package options

import (
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/server"
)

type FeatureOptions struct {
	EnableProfiling           bool
	EnableContentionProfiling bool
	EnableSwaggerUI           bool
}

func NewFeatureOptions() *FeatureOptions {
	defaults := server.NewConfig(serializer.CodecFactory{})

	return &FeatureOptions{
		EnableProfiling:           defaults.EnableProfiling,
		EnableContentionProfiling: defaults.EnableContentionProfiling,
		EnableSwaggerUI:           defaults.EnableSwaggerUI,
	}
}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "profiling", o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")
	fs.BoolVar(&o.EnableContentionProfiling, "contention-profiling", o.EnableContentionProfiling,
		"Enable lock contention profiling, if profiling is enabled")
	fs.BoolVar(&o.EnableSwaggerUI, "enable-swagger-ui", o.EnableSwaggerUI,
		"Enables swagger ui on the apiserver at /swagger-ui")
}

func (o *FeatureOptions) ApplyTo(c *server.Config) error {
	if o == nil {
		return nil
	}

	c.EnableProfiling = o.EnableProfiling
	c.EnableContentionProfiling = o.EnableContentionProfiling
	c.EnableSwaggerUI = o.EnableSwaggerUI

	return nil
}

func (o *FeatureOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
