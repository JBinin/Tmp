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
package aaparser

import (
	"testing"
)

type versionExpected struct {
	output  string
	version int
}

func TestParseVersion(t *testing.T) {
	versions := []versionExpected{
		{
			output: `AppArmor parser version 2.10
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 210000,
		},
		{
			output: `AppArmor parser version 2.8
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 208000,
		},
		{
			output: `AppArmor parser version 2.20
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 220000,
		},
		{
			output: `AppArmor parser version 2.05
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 205000,
		},
		{
			output: `AppArmor parser version 2.9.95
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 209095,
		},
		{
			output: `AppArmor parser version 3.14.159
Copyright (C) 1999-2008 Novell Inc.
Copyright 2009-2012 Canonical Ltd.

`,
			version: 314159,
		},
	}

	for _, v := range versions {
		version, err := parseVersion(v.output)
		if err != nil {
			t.Fatalf("expected error to be nil for %#v, got: %v", v, err)
		}
		if version != v.version {
			t.Fatalf("expected version to be %d, was %d, for: %#v\n", v.version, version, v)
		}
	}
}
