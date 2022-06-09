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
package dockerfile

import (
	"strings"

	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder/dockerfile/parser"
)

func TestBuildProcessLabels(t *testing.T) {
	dockerfile := "FROM scratch"
	d := parser.Directive{}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &d)
	n, err := parser.Parse(strings.NewReader(dockerfile), &d)
	if err != nil {
		t.Fatalf("Error when parsing Dockerfile: %s", err)
	}

	options := &types.ImageBuildOptions{
		Labels: map[string]string{
			"org.e": "cli-e",
			"org.d": "cli-d",
			"org.c": "cli-c",
			"org.b": "cli-b",
			"org.a": "cli-a",
		},
	}
	b := &Builder{
		runConfig:  &container.Config{},
		options:    options,
		directive:  d,
		dockerfile: n,
	}
	err = b.processLabels()
	if err != nil {
		t.Fatalf("Error when processing labels: %s", err)
	}

	expected := []string{
		"FROM scratch",
		`LABEL "org.a"='cli-a' "org.b"='cli-b' "org.c"='cli-c' "org.d"='cli-d' "org.e"='cli-e'`,
	}
	if len(b.dockerfile.Children) != 2 {
		t.Fatalf("Expect 2, got %d", len(b.dockerfile.Children))
	}
	for i, v := range b.dockerfile.Children {
		if v.Original != expected[i] {
			t.Fatalf("Expect '%s' for %dth children, got, '%s'", expected[i], i, v.Original)
		}
	}
}
