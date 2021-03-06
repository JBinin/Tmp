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
package sqlite3

/*
#ifndef USE_LIBSQLITE3
#include <sqlite3-binding.h>
#else
#include <sqlite3.h>
#endif
*/
import "C"
import (
	"reflect"
	"time"
)

// ColumnTypeDatabaseTypeName implement RowsColumnTypeDatabaseTypeName.
func (rc *SQLiteRows) ColumnTypeDatabaseTypeName(i int) string {
	return C.GoString(C.sqlite3_column_decltype(rc.s.s, C.int(i)))
}

/*
func (rc *SQLiteRows) ColumnTypeLength(index int) (length int64, ok bool) {
	return 0, false
}

func (rc *SQLiteRows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	return 0, 0, false
}
*/

// ColumnTypeNullable implement RowsColumnTypeNullable.
func (rc *SQLiteRows) ColumnTypeNullable(i int) (nullable, ok bool) {
	return true, true
}

// ColumnTypeScanType implement RowsColumnTypeScanType.
func (rc *SQLiteRows) ColumnTypeScanType(i int) reflect.Type {
	switch C.sqlite3_column_type(rc.s.s, C.int(i)) {
	case C.SQLITE_INTEGER:
		switch C.GoString(C.sqlite3_column_decltype(rc.s.s, C.int(i))) {
		case "timestamp", "datetime", "date":
			return reflect.TypeOf(time.Time{})
		case "boolean":
			return reflect.TypeOf(false)
		}
		return reflect.TypeOf(int64(0))
	case C.SQLITE_FLOAT:
		return reflect.TypeOf(float64(0))
	case C.SQLITE_BLOB:
		return reflect.SliceOf(reflect.TypeOf(byte(0)))
	case C.SQLITE_NULL:
		return reflect.TypeOf(nil)
	case C.SQLITE_TEXT:
		return reflect.TypeOf("")
	}
	return reflect.SliceOf(reflect.TypeOf(byte(0)))
}
