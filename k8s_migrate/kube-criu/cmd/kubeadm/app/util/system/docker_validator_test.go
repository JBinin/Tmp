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
Copyright 2016 The Kubernetes Authors.

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

package system

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateDockerInfo(t *testing.T) {
	v := &DockerValidator{
		Reporter: DefaultReporter,
	}
	spec := &DockerSpec{
		Version:     []string{`1\.1[1-3]\..*`, `17\.03\..*`}, // Requires [1.11, 17.03].
		GraphDriver: []string{"driver_1", "driver_2"},
	}
	for _, test := range []struct {
		info types.Info
		err  bool
		warn bool
	}{
		{
			info: types.Info{Driver: "driver_1", ServerVersion: "1.10.1"},
			err:  true,
			warn: false,
		},
		{
			info: types.Info{Driver: "bad_driver", ServerVersion: "1.11.1"},
			err:  true,
			warn: false,
		},
		{
			info: types.Info{Driver: "driver_1", ServerVersion: "1.11.1"},
			err:  false,
			warn: false,
		},
		{
			info: types.Info{Driver: "driver_2", ServerVersion: "1.12.1"},
			err:  false,
			warn: false,
		},
		{
			info: types.Info{Driver: "driver_2", ServerVersion: "1.13.1"},
			err:  false,
			warn: false,
		},
		{
			info: types.Info{Driver: "driver_2", ServerVersion: "17.03.0-ce"},
			err:  false,
			warn: false,
		},
		{
			info: types.Info{Driver: "driver_2", ServerVersion: "17.06.0-ce"},
			err:  false,
			warn: true,
		},
	} {
		warn, err := v.validateDockerInfo(spec, test.info)
		if !test.err {
			assert.Nil(t, err, "Expect error not to occur with docker info %+v", test.info)
		} else {
			assert.NotNil(t, err, "Expect error to occur with docker info %+v", test.info)
		}
		if !test.warn {
			assert.Nil(t, warn, "Expect error not to occur with docker info %+v", test.info)
		} else {
			assert.NotNil(t, warn, "Expect error to occur with docker info %+v", test.info)
		}

	}
}
