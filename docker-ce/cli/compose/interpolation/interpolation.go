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
	"fmt"

	"github.com/docker/docker/cli/compose/template"
	"github.com/docker/docker/cli/compose/types"
)

// Interpolate replaces variables in a string with the values from a mapping
func Interpolate(config types.Dict, section string, mapping template.Mapping) (types.Dict, error) {
	out := types.Dict{}

	for name, item := range config {
		if item == nil {
			out[name] = nil
			continue
		}
		interpolatedItem, err := interpolateSectionItem(name, item.(types.Dict), section, mapping)
		if err != nil {
			return nil, err
		}
		out[name] = interpolatedItem
	}

	return out, nil
}

func interpolateSectionItem(
	name string,
	item types.Dict,
	section string,
	mapping template.Mapping,
) (types.Dict, error) {

	out := types.Dict{}

	for key, value := range item {
		interpolatedValue, err := recursiveInterpolate(value, mapping)
		if err != nil {
			return nil, fmt.Errorf(
				"Invalid interpolation format for %#v option in %s %#v: %#v",
				key, section, name, err.Template,
			)
		}
		out[key] = interpolatedValue
	}

	return out, nil

}

func recursiveInterpolate(
	value interface{},
	mapping template.Mapping,
) (interface{}, *template.InvalidTemplateError) {

	switch value := value.(type) {

	case string:
		return template.Substitute(value, mapping)

	case types.Dict:
		out := types.Dict{}
		for key, elem := range value {
			interpolatedElem, err := recursiveInterpolate(elem, mapping)
			if err != nil {
				return nil, err
			}
			out[key] = interpolatedElem
		}
		return out, nil

	case []interface{}:
		out := make([]interface{}, len(value))
		for i, elem := range value {
			interpolatedElem, err := recursiveInterpolate(elem, mapping)
			if err != nil {
				return nil, err
			}
			out[i] = interpolatedElem
		}
		return out, nil

	default:
		return value, nil

	}

}
