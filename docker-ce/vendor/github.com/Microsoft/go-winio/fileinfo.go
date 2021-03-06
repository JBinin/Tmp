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

package winio

import (
	"os"
	"syscall"
	"unsafe"
)

//sys getFileInformationByHandleEx(h syscall.Handle, class uint32, buffer *byte, size uint32) (err error) = GetFileInformationByHandleEx
//sys setFileInformationByHandle(h syscall.Handle, class uint32, buffer *byte, size uint32) (err error) = SetFileInformationByHandle

const (
	fileBasicInfo = 0
	fileIDInfo    = 0x12
)

// FileBasicInfo contains file access time and file attributes information.
type FileBasicInfo struct {
	CreationTime, LastAccessTime, LastWriteTime, ChangeTime syscall.Filetime
	FileAttributes                                          uintptr // includes padding
}

// GetFileBasicInfo retrieves times and attributes for a file.
func GetFileBasicInfo(f *os.File) (*FileBasicInfo, error) {
	bi := &FileBasicInfo{}
	if err := getFileInformationByHandleEx(syscall.Handle(f.Fd()), fileBasicInfo, (*byte)(unsafe.Pointer(bi)), uint32(unsafe.Sizeof(*bi))); err != nil {
		return nil, &os.PathError{Op: "GetFileInformationByHandleEx", Path: f.Name(), Err: err}
	}
	return bi, nil
}

// SetFileBasicInfo sets times and attributes for a file.
func SetFileBasicInfo(f *os.File, bi *FileBasicInfo) error {
	if err := setFileInformationByHandle(syscall.Handle(f.Fd()), fileBasicInfo, (*byte)(unsafe.Pointer(bi)), uint32(unsafe.Sizeof(*bi))); err != nil {
		return &os.PathError{Op: "SetFileInformationByHandle", Path: f.Name(), Err: err}
	}
	return nil
}

// FileIDInfo contains the volume serial number and file ID for a file. This pair should be
// unique on a system.
type FileIDInfo struct {
	VolumeSerialNumber uint64
	FileID             [16]byte
}

// GetFileID retrieves the unique (volume, file ID) pair for a file.
func GetFileID(f *os.File) (*FileIDInfo, error) {
	fileID := &FileIDInfo{}
	if err := getFileInformationByHandleEx(syscall.Handle(f.Fd()), fileIDInfo, (*byte)(unsafe.Pointer(fileID)), uint32(unsafe.Sizeof(*fileID))); err != nil {
		return nil, &os.PathError{Op: "GetFileInformationByHandleEx", Path: f.Name(), Err: err}
	}
	return fileID, nil
}
