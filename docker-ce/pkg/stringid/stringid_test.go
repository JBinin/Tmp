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
package stringid

import (
	"strings"
	"testing"
)

func TestGenerateRandomID(t *testing.T) {
	id := GenerateRandomID()

	if len(id) != 64 {
		t.Fatalf("Id returned is incorrect: %s", id)
	}
}

func TestGenerateNonCryptoID(t *testing.T) {
	id := GenerateNonCryptoID()

	if len(id) != 64 {
		t.Fatalf("Id returned is incorrect: %s", id)
	}
}

func TestShortenId(t *testing.T) {
	id := "90435eec5c4e124e741ef731e118be2fc799a68aba0466ec17717f24ce2ae6a2"
	truncID := TruncateID(id)
	if truncID != "90435eec5c4e" {
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	}
}

func TestShortenSha256Id(t *testing.T) {
	id := "sha256:4e38e38c8ce0b8d9041a9c4fefe786631d1416225e13b0bfe8cfa2321aec4bba"
	truncID := TruncateID(id)
	if truncID != "4e38e38c8ce0" {
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	}
}

func TestShortenIdEmpty(t *testing.T) {
	id := ""
	truncID := TruncateID(id)
	if len(truncID) > len(id) {
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	}
}

func TestShortenIdInvalid(t *testing.T) {
	id := "1234"
	truncID := TruncateID(id)
	if len(truncID) != len(id) {
		t.Fatalf("Id returned is incorrect: truncate on %s returned %s", id, truncID)
	}
}

func TestIsShortIDNonHex(t *testing.T) {
	id := "some non-hex value"
	if IsShortID(id) {
		t.Fatalf("%s is not a short ID", id)
	}
}

func TestIsShortIDNotCorrectSize(t *testing.T) {
	id := strings.Repeat("a", shortLen+1)
	if IsShortID(id) {
		t.Fatalf("%s is not a short ID", id)
	}
	id = strings.Repeat("a", shortLen-1)
	if IsShortID(id) {
		t.Fatalf("%s is not a short ID", id)
	}
}
