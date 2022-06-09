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
package types

// BoolValues defines the name and value mappings for ParseBool.
var BoolValues = map[string]interface{}{
	"true": true, "yes": true, "on": true, "1": true,
	"false": false, "no": false, "off": false, "0": false,
}

var boolParser = func() *EnumParser {
	ep := &EnumParser{}
	ep.AddVals(BoolValues)
	return ep
}()

// ParseBool parses bool values according to the definitions in BoolValues.
// Parsing is case-insensitive.
func ParseBool(s string) (bool, error) {
	v, err := boolParser.Parse(s)
	if err != nil {
		return false, err
	}
	return v.(bool), nil
}
