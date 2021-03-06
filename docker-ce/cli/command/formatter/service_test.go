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
package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestServiceContextWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
		// Errors
		{
			Context{Format: "{{InvalidFunction}}"},
			`Template parsing error: template: :1: function "InvalidFunction" not defined
`,
		},
		{
			Context{Format: "{{nil}}"},
			`Template parsing error: template: :1:2: executing "" at <nil>: nil is not a command
`,
		},
		// Table format
		{
			Context{Format: NewServiceListFormat("table", false)},
			`ID                  NAME                MODE                REPLICAS            IMAGE
id_baz              baz                 global              2/4                 
id_bar              bar                 replicated          2/4                 
`,
		},
		{
			Context{Format: NewServiceListFormat("table", true)},
			`id_baz
id_bar
`,
		},
		{
			Context{Format: NewServiceListFormat("table {{.Name}}", false)},
			`NAME
baz
bar
`,
		},
		{
			Context{Format: NewServiceListFormat("table {{.Name}}", true)},
			`NAME
baz
bar
`,
		},
		// Raw Format
		{
			Context{Format: NewServiceListFormat("raw", false)},
			`id: id_baz
name: baz
mode: global
replicas: 2/4
image: 

id: id_bar
name: bar
mode: replicated
replicas: 2/4
image: 

`,
		},
		{
			Context{Format: NewServiceListFormat("raw", true)},
			`id: id_baz
id: id_bar
`,
		},
		// Custom Format
		{
			Context{Format: NewServiceListFormat("{{.Name}}", false)},
			`baz
bar
`,
		},
	}

	for _, testcase := range cases {
		services := []swarm.Service{
			{ID: "id_baz", Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "baz"}}},
			{ID: "id_bar", Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "bar"}}},
		}
		info := map[string]ServiceListInfo{
			"id_baz": {
				Mode:     "global",
				Replicas: "2/4",
			},
			"id_bar": {
				Mode:     "replicated",
				Replicas: "2/4",
			},
		}
		out := bytes.NewBufferString("")
		testcase.context.Output = out
		err := ServiceListWrite(testcase.context, services, info)
		if err != nil {
			assert.Error(t, err, testcase.expected)
		} else {
			assert.Equal(t, out.String(), testcase.expected)
		}
	}
}

func TestServiceContextWriteJSON(t *testing.T) {
	services := []swarm.Service{
		{ID: "id_baz", Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "baz"}}},
		{ID: "id_bar", Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "bar"}}},
	}
	info := map[string]ServiceListInfo{
		"id_baz": {
			Mode:     "global",
			Replicas: "2/4",
		},
		"id_bar": {
			Mode:     "replicated",
			Replicas: "2/4",
		},
	}
	expectedJSONs := []map[string]interface{}{
		{"ID": "id_baz", "Name": "baz", "Mode": "global", "Replicas": "2/4", "Image": ""},
		{"ID": "id_bar", "Name": "bar", "Mode": "replicated", "Replicas": "2/4", "Image": ""},
	}

	out := bytes.NewBufferString("")
	err := ServiceListWrite(Context{Format: "{{json .}}", Output: out}, services, info)
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		t.Logf("Output: line %d: %s", i, line)
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatal(err)
		}
		assert.DeepEqual(t, m, expectedJSONs[i])
	}
}
func TestServiceContextWriteJSONField(t *testing.T) {
	services := []swarm.Service{
		{ID: "id_baz", Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "baz"}}},
		{ID: "id_bar", Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "bar"}}},
	}
	info := map[string]ServiceListInfo{
		"id_baz": {
			Mode:     "global",
			Replicas: "2/4",
		},
		"id_bar": {
			Mode:     "replicated",
			Replicas: "2/4",
		},
	}
	out := bytes.NewBufferString("")
	err := ServiceListWrite(Context{Format: "{{json .Name}}", Output: out}, services, info)
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		t.Logf("Output: line %d: %s", i, line)
		var s string
		if err := json.Unmarshal([]byte(line), &s); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, s, services[i].Spec.Name)
	}
}
