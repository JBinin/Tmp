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
// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"database/sql/driver"
	"io"
)

type mysqlField struct {
	tableName string
	name      string
	flags     fieldFlag
	fieldType byte
	decimals  byte
}

type mysqlRows struct {
	mc      *mysqlConn
	columns []mysqlField
}

type binaryRows struct {
	mysqlRows
}

type textRows struct {
	mysqlRows
}

type emptyRows struct{}

func (rows *mysqlRows) Columns() []string {
	columns := make([]string, len(rows.columns))
	if rows.mc != nil && rows.mc.cfg.ColumnsWithAlias {
		for i := range columns {
			if tableName := rows.columns[i].tableName; len(tableName) > 0 {
				columns[i] = tableName + "." + rows.columns[i].name
			} else {
				columns[i] = rows.columns[i].name
			}
		}
	} else {
		for i := range columns {
			columns[i] = rows.columns[i].name
		}
	}
	return columns
}

func (rows *mysqlRows) Close() error {
	mc := rows.mc
	if mc == nil {
		return nil
	}
	if mc.netConn == nil {
		return ErrInvalidConn
	}

	// Remove unread packets from stream
	err := mc.readUntilEOF()
	if err == nil {
		if err = mc.discardResults(); err != nil {
			return err
		}
	}

	rows.mc = nil
	return err
}

func (rows *binaryRows) Next(dest []driver.Value) error {
	if mc := rows.mc; mc != nil {
		if mc.netConn == nil {
			return ErrInvalidConn
		}

		// Fetch next row from stream
		return rows.readRow(dest)
	}
	return io.EOF
}

func (rows *textRows) Next(dest []driver.Value) error {
	if mc := rows.mc; mc != nil {
		if mc.netConn == nil {
			return ErrInvalidConn
		}

		// Fetch next row from stream
		return rows.readRow(dest)
	}
	return io.EOF
}

func (rows emptyRows) Columns() []string {
	return nil
}

func (rows emptyRows) Close() error {
	return nil
}

func (rows emptyRows) Next(dest []driver.Value) error {
	return io.EOF
}
