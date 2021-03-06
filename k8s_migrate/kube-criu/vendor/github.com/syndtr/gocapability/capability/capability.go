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
// Copyright (c) 2013, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package capability provides utilities for manipulating POSIX capabilities.
package capability

type Capabilities interface {
	// Get check whether a capability present in the given
	// capabilities set. The 'which' value should be one of EFFECTIVE,
	// PERMITTED, INHERITABLE, BOUNDING or AMBIENT.
	Get(which CapType, what Cap) bool

	// Empty check whether all capability bits of the given capabilities
	// set are zero. The 'which' value should be one of EFFECTIVE,
	// PERMITTED, INHERITABLE, BOUNDING or AMBIENT.
	Empty(which CapType) bool

	// Full check whether all capability bits of the given capabilities
	// set are one. The 'which' value should be one of EFFECTIVE,
	// PERMITTED, INHERITABLE, BOUNDING or AMBIENT.
	Full(which CapType) bool

	// Set sets capabilities of the given capabilities sets. The
	// 'which' value should be one or combination (OR'ed) of EFFECTIVE,
	// PERMITTED, INHERITABLE, BOUNDING or AMBIENT.
	Set(which CapType, caps ...Cap)

	// Unset unsets capabilities of the given capabilities sets. The
	// 'which' value should be one or combination (OR'ed) of EFFECTIVE,
	// PERMITTED, INHERITABLE, BOUNDING or AMBIENT.
	Unset(which CapType, caps ...Cap)

	// Fill sets all bits of the given capabilities kind to one. The
	// 'kind' value should be one or combination (OR'ed) of CAPS,
	// BOUNDS or AMBS.
	Fill(kind CapType)

	// Clear sets all bits of the given capabilities kind to zero. The
	// 'kind' value should be one or combination (OR'ed) of CAPS,
	// BOUNDS or AMBS.
	Clear(kind CapType)

	// String return current capabilities state of the given capabilities
	// set as string. The 'which' value should be one of EFFECTIVE,
	// PERMITTED, INHERITABLE BOUNDING or AMBIENT
	StringCap(which CapType) string

	// String return current capabilities state as string.
	String() string

	// Load load actual capabilities value. This will overwrite all
	// outstanding changes.
	Load() error

	// Apply apply the capabilities settings, so all changes will take
	// effect.
	Apply(kind CapType) error
}

// NewPid create new initialized Capabilities object for given pid when it
// is nonzero, or for the current pid if pid is 0
func NewPid(pid int) (Capabilities, error) {
	return newPid(pid)
}

// NewFile create new initialized Capabilities object for given named file.
func NewFile(name string) (Capabilities, error) {
	return newFile(name)
}
