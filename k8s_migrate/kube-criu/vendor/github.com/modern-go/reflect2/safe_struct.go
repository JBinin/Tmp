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
package reflect2

type safeStructType struct {
	safeType
}

func (type2 *safeStructType) FieldByName(name string) StructField {
	field, found := type2.Type.FieldByName(name)
	if !found {
		panic("field " + name + " not found")
	}
	return &safeField{StructField: field}
}

func (type2 *safeStructType) Field(i int) StructField {
	return &safeField{StructField: type2.Type.Field(i)}
}

func (type2 *safeStructType) FieldByIndex(index []int) StructField {
	return &safeField{StructField: type2.Type.FieldByIndex(index)}
}

func (type2 *safeStructType) FieldByNameFunc(match func(string) bool) StructField {
	field, found := type2.Type.FieldByNameFunc(match)
	if !found {
		panic("field match condition not found in " + type2.Type.String())
	}
	return &safeField{StructField: field}
}
