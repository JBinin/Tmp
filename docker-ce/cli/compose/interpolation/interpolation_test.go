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
package interpolation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/cli/compose/types"
)

var defaults = map[string]string{
	"USER": "jenny",
	"FOO":  "bar",
}

func defaultMapping(name string) (string, bool) {
	val, ok := defaults[name]
	return val, ok
}

func TestInterpolate(t *testing.T) {
	services := types.Dict{
		"servicea": types.Dict{
			"image":   "example:${USER}",
			"volumes": []interface{}{"$FOO:/target"},
			"logging": types.Dict{
				"driver": "${FOO}",
				"options": types.Dict{
					"user": "$USER",
				},
			},
		},
	}
	expected := types.Dict{
		"servicea": types.Dict{
			"image":   "example:jenny",
			"volumes": []interface{}{"bar:/target"},
			"logging": types.Dict{
				"driver": "bar",
				"options": types.Dict{
					"user": "jenny",
				},
			},
		},
	}
	result, err := Interpolate(services, "service", defaultMapping)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestInvalidInterpolation(t *testing.T) {
	services := types.Dict{
		"servicea": types.Dict{
			"image": "${",
		},
	}
	_, err := Interpolate(services, "service", defaultMapping)
	assert.EqualError(t, err, `Invalid interpolation format for "image" option in service "servicea": "${"`)
}
