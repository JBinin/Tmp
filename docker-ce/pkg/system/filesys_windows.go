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

package system

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	winio "github.com/Microsoft/go-winio"
)

// MkdirAllWithACL is a wrapper for MkdirAll that creates a directory
// ACL'd for Builtin Administrators and Local System.
func MkdirAllWithACL(path string, perm os.FileMode) error {
	return mkdirall(path, true)
}

// MkdirAll implementation that is volume path aware for Windows.
func MkdirAll(path string, _ os.FileMode) error {
	return mkdirall(path, false)
}

// mkdirall is a custom version of os.MkdirAll modified for use on Windows
// so that it is both volume path aware, and can create a directory with
// a DACL.
func mkdirall(path string, adminAndLocalSystem bool) error {
	if re := regexp.MustCompile(`^\\\\\?\\Volume{[a-z0-9-]+}$`); re.MatchString(path) {
		return nil
	}

	// The rest of this method is largely copied from os.MkdirAll and should be kept
	// as-is to ensure compatibility.

	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{
			Op:   "mkdir",
			Path: path,
			Err:  syscall.ENOTDIR,
		}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = mkdirall(path[0:j-1], false)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke os.Mkdir or mkdirWithACL and use its result.
	if adminAndLocalSystem {
		err = mkdirWithACL(path)
	} else {
		err = os.Mkdir(path, 0)
	}

	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

// mkdirWithACL creates a new directory. If there is an error, it will be of
// type *PathError. .
//
// This is a modified and combined version of os.Mkdir and syscall.Mkdir
// in golang to cater for creating a directory am ACL permitting full
// access, with inheritance, to any subfolder/file for Built-in Administrators
// and Local System.
func mkdirWithACL(name string) error {
	sa := syscall.SecurityAttributes{Length: 0}
	sddl := "D:P(A;OICI;GA;;;BA)(A;OICI;GA;;;SY)"
	sd, err := winio.SddlToSecurityDescriptor(sddl)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	sa.SecurityDescriptor = uintptr(unsafe.Pointer(&sd[0]))

	namep, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}

	e := syscall.CreateDirectory(namep, &sa)
	if e != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: e}
	}
	return nil
}

// IsAbs is a platform-specific wrapper for filepath.IsAbs. On Windows,
// golang filepath.IsAbs does not consider a path \windows\system32 as absolute
// as it doesn't start with a drive-letter/colon combination. However, in
// docker we need to verify things such as WORKDIR /windows/system32 in
// a Dockerfile (which gets translated to \windows\system32 when being processed
// by the daemon. This SHOULD be treated as absolute from a docker processing
// perspective.
func IsAbs(path string) bool {
	if !filepath.IsAbs(path) {
		if !strings.HasPrefix(path, string(os.PathSeparator)) {
			return false
		}
	}
	return true
}

// The origin of the functions below here are the golang OS and syscall packages,
// slightly modified to only cope with files, not directories due to the
// specific use case.
//
// The alteration is to allow a file on Windows to be opened with
// FILE_FLAG_SEQUENTIAL_SCAN (particular for docker load), to avoid eating
// the standby list, particularly when accessing large files such as layer.tar.

// CreateSequential creates the named file with mode 0666 (before umask), truncating
// it if it already exists. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.
func CreateSequential(name string) (*os.File, error) {
	return OpenFileSequential(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
}

// OpenSequential opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func OpenSequential(name string) (*os.File, error) {
	return OpenFileSequential(name, os.O_RDONLY, 0)
}

// OpenFileSequential is the generalized open call; most users will use Open
// or Create instead.
// If there is an error, it will be of type *PathError.
func OpenFileSequential(name string, flag int, _ os.FileMode) (*os.File, error) {
	if name == "" {
		return nil, &os.PathError{Op: "open", Path: name, Err: syscall.ENOENT}
	}
	r, errf := syscallOpenFileSequential(name, flag, 0)
	if errf == nil {
		return r, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: errf}
}

func syscallOpenFileSequential(name string, flag int, _ os.FileMode) (file *os.File, err error) {
	r, e := syscallOpenSequential(name, flag|syscall.O_CLOEXEC, 0)
	if e != nil {
		return nil, e
	}
	return os.NewFile(uintptr(r), name), nil
}

func makeInheritSa() *syscall.SecurityAttributes {
	var sa syscall.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	return &sa
}

func syscallOpenSequential(path string, mode int, _ uint32) (fd syscall.Handle, err error) {
	if len(path) == 0 {
		return syscall.InvalidHandle, syscall.ERROR_FILE_NOT_FOUND
	}
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return syscall.InvalidHandle, err
	}
	var access uint32
	switch mode & (syscall.O_RDONLY | syscall.O_WRONLY | syscall.O_RDWR) {
	case syscall.O_RDONLY:
		access = syscall.GENERIC_READ
	case syscall.O_WRONLY:
		access = syscall.GENERIC_WRITE
	case syscall.O_RDWR:
		access = syscall.GENERIC_READ | syscall.GENERIC_WRITE
	}
	if mode&syscall.O_CREAT != 0 {
		access |= syscall.GENERIC_WRITE
	}
	if mode&syscall.O_APPEND != 0 {
		access &^= syscall.GENERIC_WRITE
		access |= syscall.FILE_APPEND_DATA
	}
	sharemode := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE)
	var sa *syscall.SecurityAttributes
	if mode&syscall.O_CLOEXEC == 0 {
		sa = makeInheritSa()
	}
	var createmode uint32
	switch {
	case mode&(syscall.O_CREAT|syscall.O_EXCL) == (syscall.O_CREAT | syscall.O_EXCL):
		createmode = syscall.CREATE_NEW
	case mode&(syscall.O_CREAT|syscall.O_TRUNC) == (syscall.O_CREAT | syscall.O_TRUNC):
		createmode = syscall.CREATE_ALWAYS
	case mode&syscall.O_CREAT == syscall.O_CREAT:
		createmode = syscall.OPEN_ALWAYS
	case mode&syscall.O_TRUNC == syscall.O_TRUNC:
		createmode = syscall.TRUNCATE_EXISTING
	default:
		createmode = syscall.OPEN_EXISTING
	}
	// Use FILE_FLAG_SEQUENTIAL_SCAN rather than FILE_ATTRIBUTE_NORMAL as implemented in golang.
	//https://msdn.microsoft.com/en-us/library/windows/desktop/aa363858(v=vs.85).aspx
	const fileFlagSequentialScan = 0x08000000 // FILE_FLAG_SEQUENTIAL_SCAN
	h, e := syscall.CreateFile(pathp, access, sharemode, sa, createmode, fileFlagSequentialScan, 0)
	return h, e
}

// Helpers for TempFileSequential
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}
func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFileSequential is a copy of ioutil.TempFile, modified to use sequential
// file access. Below is the original comment from golang:
// TempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, opens the file for reading
// and writing, and returns the resulting *os.File.
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func TempFileSequential(dir, prefix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix())
		f, err = OpenFileSequential(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}
