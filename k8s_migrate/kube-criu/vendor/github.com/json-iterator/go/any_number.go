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

import (
	"io"
	"unsafe"
)

type numberLazyAny struct {
	baseAny
	cfg *frozenConfig
	buf []byte
	err error
}

func (any *numberLazyAny) ValueType() ValueType {
	return NumberValue
}

func (any *numberLazyAny) MustBeValid() Any {
	return any
}

func (any *numberLazyAny) LastError() error {
	return any.err
}

func (any *numberLazyAny) ToBool() bool {
	return any.ToFloat64() != 0
}

func (any *numberLazyAny) ToInt() int {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadInt()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToInt32() int32 {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadInt32()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToInt64() int64 {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadInt64()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToUint() uint {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadUint()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToUint32() uint32 {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadUint32()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToUint64() uint64 {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadUint64()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToFloat32() float32 {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadFloat32()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToFloat64() float64 {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	val := iter.ReadFloat64()
	if iter.Error != nil && iter.Error != io.EOF {
		any.err = iter.Error
	}
	return val
}

func (any *numberLazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *numberLazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *numberLazyAny) GetInterface() interface{} {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	return iter.Read()
}
