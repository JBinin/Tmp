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
// +build windows

package ioutils

import (
	"io/ioutil"

	"github.com/docker/docker/pkg/longpath"
)

// TempDir is the equivalent of ioutil.TempDir, except that the result is in Windows longpath format.
func TempDir(dir, prefix string) (string, error) {
	tempDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", err
	}
	return longpath.AddPrefix(tempDir), nil
}
