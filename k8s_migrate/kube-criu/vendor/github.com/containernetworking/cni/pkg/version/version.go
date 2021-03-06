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
// Copyright 2016 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/020"
	"github.com/containernetworking/cni/pkg/types/current"
)

// Current reports the version of the CNI spec implemented by this library
func Current() string {
	return "0.3.1"
}

// Legacy PluginInfo describes a plugin that is backwards compatible with the
// CNI spec version 0.1.0.  In particular, a runtime compiled against the 0.1.0
// library ought to work correctly with a plugin that reports support for
// Legacy versions.
//
// Any future CNI spec versions which meet this definition should be added to
// this list.
var Legacy = PluginSupports("0.1.0", "0.2.0")
var All = PluginSupports("0.1.0", "0.2.0", "0.3.0", "0.3.1")

var resultFactories = []struct {
	supportedVersions []string
	newResult         types.ResultFactoryFunc
}{
	{current.SupportedVersions, current.NewResult},
	{types020.SupportedVersions, types020.NewResult},
}

// Finds a Result object matching the requested version (if any) and asks
// that object to parse the plugin result, returning an error if parsing failed.
func NewResult(version string, resultBytes []byte) (types.Result, error) {
	reconciler := &Reconciler{}
	for _, resultFactory := range resultFactories {
		err := reconciler.CheckRaw(version, resultFactory.supportedVersions)
		if err == nil {
			// Result supports this version
			return resultFactory.newResult(resultBytes)
		}
	}

	return nil, fmt.Errorf("unsupported CNI result version %q", version)
}
