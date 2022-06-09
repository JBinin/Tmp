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
package jsoniter

type trueAny struct {
	baseAny
}

func (any *trueAny) LastError() error {
	return nil
}

func (any *trueAny) ToBool() bool {
	return true
}

func (any *trueAny) ToInt() int {
	return 1
}

func (any *trueAny) ToInt32() int32 {
	return 1
}

func (any *trueAny) ToInt64() int64 {
	return 1
}

func (any *trueAny) ToUint() uint {
	return 1
}

func (any *trueAny) ToUint32() uint32 {
	return 1
}

func (any *trueAny) ToUint64() uint64 {
	return 1
}

func (any *trueAny) ToFloat32() float32 {
	return 1
}

func (any *trueAny) ToFloat64() float64 {
	return 1
}

func (any *trueAny) ToString() string {
	return "true"
}

func (any *trueAny) WriteTo(stream *Stream) {
	stream.WriteTrue()
}

func (any *trueAny) Parse() *Iterator {
	return nil
}

func (any *trueAny) GetInterface() interface{} {
	return true
}

func (any *trueAny) ValueType() ValueType {
	return BoolValue
}

func (any *trueAny) MustBeValid() Any {
	return any
}

type falseAny struct {
	baseAny
}

func (any *falseAny) LastError() error {
	return nil
}

func (any *falseAny) ToBool() bool {
	return false
}

func (any *falseAny) ToInt() int {
	return 0
}

func (any *falseAny) ToInt32() int32 {
	return 0
}

func (any *falseAny) ToInt64() int64 {
	return 0
}

func (any *falseAny) ToUint() uint {
	return 0
}

func (any *falseAny) ToUint32() uint32 {
	return 0
}

func (any *falseAny) ToUint64() uint64 {
	return 0
}

func (any *falseAny) ToFloat32() float32 {
	return 0
}

func (any *falseAny) ToFloat64() float64 {
	return 0
}

func (any *falseAny) ToString() string {
	return "false"
}

func (any *falseAny) WriteTo(stream *Stream) {
	stream.WriteFalse()
}

func (any *falseAny) Parse() *Iterator {
	return nil
}

func (any *falseAny) GetInterface() interface{} {
	return false
}

func (any *falseAny) ValueType() ValueType {
	return BoolValue
}

func (any *falseAny) MustBeValid() Any {
	return any
}