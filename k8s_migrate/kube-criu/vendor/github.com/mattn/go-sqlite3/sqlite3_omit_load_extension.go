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
// Copyright (C) 2014 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
// +build sqlite_omit_load_extension

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_OMIT_LOAD_EXTENSION
*/
import "C"
import (
	"errors"
)

func (c *SQLiteConn) loadExtensions(extensions []string) error {
	return errors.New("Extensions have been disabled for static builds")
}

func (c *SQLiteConn) LoadExtension(lib string, entry string) error {
	return errors.New("Extensions have been disabled for static builds")
}
