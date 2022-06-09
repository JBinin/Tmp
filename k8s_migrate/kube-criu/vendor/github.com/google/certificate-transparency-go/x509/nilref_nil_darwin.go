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
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cgo,!arm,!arm64,!ios,!go1.10

package x509

/*
#cgo CFLAGS: -mmacosx-version-min=10.6 -D__MAC_OS_X_VERSION_MAX_ALLOWED=1080
#cgo LDFLAGS: -framework CoreFoundation -framework Security

#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

// For Go versions before 1.10, nil values for Apple's CoreFoundation
// CF*Ref types were represented by nil.  See:
//   https://github.com/golang/go/commit/b868616b63a8
func setNilCFRef(v *C.CFDataRef) {
	*v = nil
}

func isNilCFRef(v C.CFDataRef) bool {
	return v == nil
}
