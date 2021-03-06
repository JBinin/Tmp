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
// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wal

import "io"

type multiReadCloser struct {
	closers []io.Closer
	reader  io.Reader
}

func (mc *multiReadCloser) Close() error {
	var err error
	for i := range mc.closers {
		err = mc.closers[i].Close()
	}
	return err
}

func (mc *multiReadCloser) Read(p []byte) (int, error) {
	return mc.reader.Read(p)
}

func MultiReadCloser(readClosers ...io.ReadCloser) io.ReadCloser {
	cs := make([]io.Closer, len(readClosers))
	rs := make([]io.Reader, len(readClosers))
	for i := range readClosers {
		cs[i] = readClosers[i]
		rs[i] = readClosers[i]
	}
	r := io.MultiReader(rs...)
	return &multiReadCloser{cs, r}
}
