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
package bundlefile

import (
	"bytes"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestLoadFileV01Success(t *testing.T) {
	reader := strings.NewReader(`{
		"Version": "0.1",
		"Services": {
			"redis": {
				"Image": "redis@sha256:4b24131101fa0117bcaa18ac37055fffd9176aa1a240392bb8ea85e0be50f2ce",
				"Networks": ["default"]
			},
			"web": {
				"Image": "dockercloud/hello-world@sha256:fe79a2cfbd17eefc344fb8419420808df95a1e22d93b7f621a7399fd1e9dca1d",
				"Networks": ["default"],
				"User": "web"
			}
		}
	}`)

	bundle, err := LoadFile(reader)
	assert.NilError(t, err)
	assert.Equal(t, bundle.Version, "0.1")
	assert.Equal(t, len(bundle.Services), 2)
}

func TestLoadFileSyntaxError(t *testing.T) {
	reader := strings.NewReader(`{
		"Version": "0.1",
		"Services": unquoted string
	}`)

	_, err := LoadFile(reader)
	assert.Error(t, err, "syntax error at byte 37: invalid character 'u'")
}

func TestLoadFileTypeError(t *testing.T) {
	reader := strings.NewReader(`{
		"Version": "0.1",
		"Services": {
			"web": {
				"Image": "redis",
				"Networks": "none"
			}
		}
	}`)

	_, err := LoadFile(reader)
	assert.Error(t, err, "Unexpected type at byte 94. Expected []string but received string")
}

func TestPrint(t *testing.T) {
	var buffer bytes.Buffer
	bundle := &Bundlefile{
		Version: "0.1",
		Services: map[string]Service{
			"web": {
				Image:   "image",
				Command: []string{"echo", "something"},
			},
		},
	}
	assert.NilError(t, Print(&buffer, bundle))
	output := buffer.String()
	assert.Contains(t, output, "\"Image\": \"image\"")
	assert.Contains(t, output,
		`"Command": [
                "echo",
                "something"
            ]`)
}
