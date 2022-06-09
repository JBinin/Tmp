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

import "testing"

// TestCheckSystemDriveAndRemoveDriveLetter tests CheckSystemDriveAndRemoveDriveLetter
func TestCheckSystemDriveAndRemoveDriveLetter(t *testing.T) {
	// Fails if not C drive.
	path, err := CheckSystemDriveAndRemoveDriveLetter(`d:\`)
	if err == nil || (err != nil && err.Error() != "The specified path is not on the system drive (C:)") {
		t.Fatalf("Expected error for d:")
	}

	// Single character is unchanged
	if path, err = CheckSystemDriveAndRemoveDriveLetter("z"); err != nil {
		t.Fatalf("Single character should pass")
	}
	if path != "z" {
		t.Fatalf("Single character should be unchanged")
	}

	// Two characters without colon is unchanged
	if path, err = CheckSystemDriveAndRemoveDriveLetter("AB"); err != nil {
		t.Fatalf("2 characters without colon should pass")
	}
	if path != "AB" {
		t.Fatalf("2 characters without colon should be unchanged")
	}

	// Abs path without drive letter
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`\l`); err != nil {
		t.Fatalf("abs path no drive letter should pass")
	}
	if path != `\l` {
		t.Fatalf("abs path without drive letter should be unchanged")
	}

	// Abs path without drive letter, linux style
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`/l`); err != nil {
		t.Fatalf("abs path no drive letter linux style should pass")
	}
	if path != `\l` {
		t.Fatalf("abs path without drive letter linux failed %s", path)
	}

	// Drive-colon should be stripped
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`c:\`); err != nil {
		t.Fatalf("An absolute path should pass")
	}
	if path != `\` {
		t.Fatalf(`An absolute path should have been shortened to \ %s`, path)
	}

	// Verify with a linux-style path
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`c:/`); err != nil {
		t.Fatalf("An absolute path should pass")
	}
	if path != `\` {
		t.Fatalf(`A linux style absolute path should have been shortened to \ %s`, path)
	}

	// Failure on c:
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`c:`); err == nil {
		t.Fatalf("c: should fail")
	}
	if err.Error() != `No relative path specified in "c:"` {
		t.Fatalf(path, err)
	}

	// Failure on d:
	if path, err = CheckSystemDriveAndRemoveDriveLetter(`d:`); err == nil {
		t.Fatalf("c: should fail")
	}
	if err.Error() != `No relative path specified in "d:"` {
		t.Fatalf(path, err)
	}
}
