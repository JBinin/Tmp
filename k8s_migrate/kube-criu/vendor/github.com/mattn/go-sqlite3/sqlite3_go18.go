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

// +build go1.8

package sqlite3

import (
	"database/sql/driver"
	"errors"

	"context"
)

// Ping implement Pinger.
func (c *SQLiteConn) Ping(ctx context.Context) error {
	if c.db == nil {
		return errors.New("Connection was closed")
	}
	return nil
}

// QueryContext implement QueryerContext.
func (c *SQLiteConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	list := make([]namedValue, len(args))
	for i, nv := range args {
		list[i] = namedValue(nv)
	}
	return c.query(ctx, query, list)
}

// ExecContext implement ExecerContext.
func (c *SQLiteConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	list := make([]namedValue, len(args))
	for i, nv := range args {
		list[i] = namedValue(nv)
	}
	return c.exec(ctx, query, list)
}

// PrepareContext implement ConnPrepareContext.
func (c *SQLiteConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.prepare(ctx, query)
}

// BeginTx implement ConnBeginTx.
func (c *SQLiteConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.begin(ctx)
}

// QueryContext implement QueryerContext.
func (s *SQLiteStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	list := make([]namedValue, len(args))
	for i, nv := range args {
		list[i] = namedValue(nv)
	}
	return s.query(ctx, list)
}

// ExecContext implement ExecerContext.
func (s *SQLiteStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	list := make([]namedValue, len(args))
	for i, nv := range args {
		list[i] = namedValue(nv)
	}
	return s.exec(ctx, list)
}
