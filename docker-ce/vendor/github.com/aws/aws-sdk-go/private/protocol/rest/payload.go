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
package rest

import "reflect"

// PayloadMember returns the payload field member of i if there is one, or nil.
func PayloadMember(i interface{}) interface{} {
	if i == nil {
		return nil
	}

	v := reflect.ValueOf(i).Elem()
	if !v.IsValid() {
		return nil
	}
	if field, ok := v.Type().FieldByName("_"); ok {
		if payloadName := field.Tag.Get("payload"); payloadName != "" {
			field, _ := v.Type().FieldByName(payloadName)
			if field.Tag.Get("type") != "structure" {
				return nil
			}

			payload := v.FieldByName(payloadName)
			if payload.IsValid() || (payload.Kind() == reflect.Ptr && !payload.IsNil()) {
				return payload.Interface()
			}
		}
	}
	return nil
}

// PayloadType returns the type of a payload field member of i if there is one, or "".
func PayloadType(i interface{}) string {
	v := reflect.Indirect(reflect.ValueOf(i))
	if !v.IsValid() {
		return ""
	}
	if field, ok := v.Type().FieldByName("_"); ok {
		if payloadName := field.Tag.Get("payload"); payloadName != "" {
			if member, ok := v.Type().FieldByName(payloadName); ok {
				return member.Tag.Get("type")
			}
		}
	}
	return ""
}
