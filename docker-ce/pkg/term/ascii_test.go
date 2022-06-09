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
package term

import "testing"

func TestToBytes(t *testing.T) {
	codes, err := ToBytes("ctrl-a,a")
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 2 {
		t.Fatalf("Expected 2 codes, got %d", len(codes))
	}
	if codes[0] != 1 || codes[1] != 97 {
		t.Fatalf("Expected '1' '97', got '%d' '%d'", codes[0], codes[1])
	}

	codes, err = ToBytes("shift-z")
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	codes, err = ToBytes("ctrl-@,ctrl-[,~,ctrl-o")
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 4 {
		t.Fatalf("Expected 4 codes, got %d", len(codes))
	}
	if codes[0] != 0 || codes[1] != 27 || codes[2] != 126 || codes[3] != 15 {
		t.Fatalf("Expected '0' '27' '126', '15', got '%d' '%d' '%d' '%d'", codes[0], codes[1], codes[2], codes[3])
	}

	codes, err = ToBytes("DEL,+")
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 2 {
		t.Fatalf("Expected 2 codes, got %d", len(codes))
	}
	if codes[0] != 127 || codes[1] != 43 {
		t.Fatalf("Expected '127 '43'', got '%d' '%d'", codes[0], codes[1])
	}
}
