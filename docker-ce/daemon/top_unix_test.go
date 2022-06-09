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
//+build !windows

package daemon

import (
	"testing"
)

func TestContainerTopValidatePSArgs(t *testing.T) {
	tests := map[string]bool{
		"ae -o uid=PID":             true,
		"ae -o \"uid= PID\"":        true,  // ascii space (0x20)
		"ae -o \"uid= PID\"":        false, // unicode space (U+2003, 0xe2 0x80 0x83)
		"ae o uid=PID":              true,
		"aeo uid=PID":               true,
		"ae -O uid=PID":             true,
		"ae -o pid=PID2 -o uid=PID": true,
		"ae -o pid=PID":             false,
		"ae -o pid=PID -o uid=PIDX": true, // FIXME: we do not need to prohibit this
		"aeo pid=PID":               false,
		"ae":                        false,
		"":                          false,
	}
	for psArgs, errExpected := range tests {
		err := validatePSArgs(psArgs)
		t.Logf("tested %q, got err=%v", psArgs, err)
		if errExpected && err == nil {
			t.Fatalf("expected error, got %v (%q)", err, psArgs)
		}
		if !errExpected && err != nil {
			t.Fatalf("expected nil, got %v (%q)", err, psArgs)
		}
	}
}

func TestContainerTopParsePSOutput(t *testing.T) {
	tests := []struct {
		output      []byte
		pids        []int
		errExpected bool
	}{
		{[]byte(`  PID COMMAND
   42 foo
   43 bar
  100 baz
`), []int{42, 43}, false},
		{[]byte(`  UID COMMAND
   42 foo
   43 bar
  100 baz
`), []int{42, 43}, true},
		// unicode space (U+2003, 0xe2 0x80 0x83)
		{[]byte(` PID COMMAND
   42 foo
   43 bar
  100 baz
`), []int{42, 43}, true},
		// the first space is U+2003, the second one is ascii.
		{[]byte(` PID COMMAND
   42 foo
   43 bar
  100 baz
`), []int{42, 43}, true},
	}

	for _, f := range tests {
		_, err := parsePSOutput(f.output, f.pids)
		t.Logf("tested %q, got err=%v", string(f.output), err)
		if f.errExpected && err == nil {
			t.Fatalf("expected error, got %v (%q)", err, string(f.output))
		}
		if !f.errExpected && err != nil {
			t.Fatalf("expected nil, got %v (%q)", err, string(f.output))
		}
	}
}
