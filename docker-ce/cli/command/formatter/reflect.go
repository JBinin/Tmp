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
package formatter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

func marshalJSON(x interface{}) ([]byte, error) {
	m, err := marshalMap(x)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// marshalMap marshals x to map[string]interface{}
func marshalMap(x interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(x)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected a pointer to a struct, got %v", val.Kind())
	}
	if val.IsNil() {
		return nil, fmt.Errorf("expected a pointer to a struct, got nil pointer")
	}
	valElem := val.Elem()
	if valElem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a pointer to a struct, got a pointer to %v", valElem.Kind())
	}
	typ := val.Type()
	m := make(map[string]interface{})
	for i := 0; i < val.NumMethod(); i++ {
		k, v, err := marshalForMethod(typ.Method(i), val.Method(i))
		if err != nil {
			return nil, err
		}
		if k != "" {
			m[k] = v
		}
	}
	return m, nil
}

var unmarshallableNames = map[string]struct{}{"FullHeader": {}}

// marshalForMethod returns the map key and the map value for marshalling the method.
// It returns ("", nil, nil) for valid but non-marshallable parameter. (e.g. "unexportedFunc()")
func marshalForMethod(typ reflect.Method, val reflect.Value) (string, interface{}, error) {
	if val.Kind() != reflect.Func {
		return "", nil, fmt.Errorf("expected func, got %v", val.Kind())
	}
	name, numIn, numOut := typ.Name, val.Type().NumIn(), val.Type().NumOut()
	_, blackListed := unmarshallableNames[name]
	// FIXME: In text/template, (numOut == 2) is marshallable,
	//        if the type of the second param is error.
	marshallable := unicode.IsUpper(rune(name[0])) && !blackListed &&
		numIn == 0 && numOut == 1
	if !marshallable {
		return "", nil, nil
	}
	result := val.Call(make([]reflect.Value, numIn))
	intf := result[0].Interface()
	return name, intf, nil
}
