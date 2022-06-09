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
package layer

import (
	"io"
	"testing"

	"github.com/opencontainers/go-digest"
)

func TestEmptyLayer(t *testing.T) {
	if EmptyLayer.ChainID() != ChainID(DigestSHA256EmptyTar) {
		t.Fatal("wrong ID for empty layer")
	}

	if EmptyLayer.DiffID() != DigestSHA256EmptyTar {
		t.Fatal("wrong DiffID for empty layer")
	}

	if EmptyLayer.Parent() != nil {
		t.Fatal("expected no parent for empty layer")
	}

	if size, err := EmptyLayer.Size(); err != nil || size != 0 {
		t.Fatal("expected zero size for empty layer")
	}

	if diffSize, err := EmptyLayer.DiffSize(); err != nil || diffSize != 0 {
		t.Fatal("expected zero diffsize for empty layer")
	}

	tarStream, err := EmptyLayer.TarStream()
	if err != nil {
		t.Fatalf("error streaming tar for empty layer: %v", err)
	}

	digester := digest.Canonical.Digester()
	_, err = io.Copy(digester.Hash(), tarStream)

	if err != nil {
		t.Fatalf("error hashing empty tar layer: %v", err)
	}

	if digester.Digest() != digest.Digest(DigestSHA256EmptyTar) {
		t.Fatal("empty layer tar stream hashes to wrong value")
	}
}
