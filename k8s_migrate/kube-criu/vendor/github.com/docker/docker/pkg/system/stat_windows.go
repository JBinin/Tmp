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
package system

import (
	"os"
	"time"
)

// StatT type contains status of a file. It contains metadata
// like permission, size, etc about a file.
type StatT struct {
	mode os.FileMode
	size int64
	mtim time.Time
}

// Size returns file's size.
func (s StatT) Size() int64 {
	return s.size
}

// Mode returns file's permission mode.
func (s StatT) Mode() os.FileMode {
	return os.FileMode(s.mode)
}

// Mtim returns file's last modification time.
func (s StatT) Mtim() time.Time {
	return time.Time(s.mtim)
}

// Stat takes a path to a file and returns
// a system.StatT type pertaining to that file.
//
// Throws an error if the file does not exist
func Stat(path string) (*StatT, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return fromStatT(&fi)
}

// fromStatT converts a os.FileInfo type to a system.StatT type
func fromStatT(fi *os.FileInfo) (*StatT, error) {
	return &StatT{
		size: (*fi).Size(),
		mode: (*fi).Mode(),
		mtim: (*fi).ModTime()}, nil
}
