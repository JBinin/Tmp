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

import (
	"reflect"
	"unsafe"
)

type safeType struct {
	reflect.Type
	cfg *frozenConfig
}

func (type2 *safeType) New() interface{} {
	return reflect.New(type2.Type).Interface()
}

func (type2 *safeType) UnsafeNew() unsafe.Pointer {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Elem() Type {
	return type2.cfg.Type2(type2.Type.Elem())
}

func (type2 *safeType) Type1() reflect.Type {
	return type2.Type
}

func (type2 *safeType) PackEFace(ptr unsafe.Pointer) interface{} {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Implements(thatType Type) bool {
	return type2.Type.Implements(thatType.Type1())
}

func (type2 *safeType) RType() uintptr {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Indirect(obj interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(obj)).Interface()
}

func (type2 *safeType) UnsafeIndirect(ptr unsafe.Pointer) interface{} {
	panic("does not support unsafe operation")
}

func (type2 *safeType) LikePtr() bool {
	panic("does not support unsafe operation")
}

func (type2 *safeType) IsNullable() bool {
	return IsNullable(type2.Kind())
}

func (type2 *safeType) IsNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	return reflect.ValueOf(obj).Elem().IsNil()
}

func (type2 *safeType) UnsafeIsNil(ptr unsafe.Pointer) bool {
	panic("does not support unsafe operation")
}

func (type2 *safeType) Set(obj interface{}, val interface{}) {
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(val).Elem())
}

func (type2 *safeType) UnsafeSet(ptr unsafe.Pointer, val unsafe.Pointer) {
	panic("does not support unsafe operation")
}

func (type2 *safeType) AssignableTo(anotherType Type) bool {
	return type2.Type1().AssignableTo(anotherType.Type1())
}
