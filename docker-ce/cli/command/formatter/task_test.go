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

func TestTaskContextWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
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
		{
			Context{Format: NewTaskFormat("table", true)},
			`taskID1
taskID2
`,
		},
		{
			Context{Format: NewTaskFormat("table {{.Name}} {{.Node}} {{.Ports}}", false)},
			`NAME                NODE                PORTS
foobar_baz foo1 
foobar_bar foo2 
`,
		},
		{
			Context{Format: NewTaskFormat("table {{.Name}}", true)},
			`NAME
foobar_baz
foobar_bar
`,
		},
		{
			Context{Format: NewTaskFormat("raw", true)},
			`id: taskID1
id: taskID2
`,
		},
		{
			Context{Format: NewTaskFormat("{{.Name}} {{.Node}}", false)},
			`foobar_baz foo1
foobar_bar foo2
`,
		},
	}

	for _, testcase := range cases {
		tasks := []swarm.Task{
			{ID: "taskID1"},
			{ID: "taskID2"},
		}
		names := map[string]string{
			"taskID1": "foobar_baz",
			"taskID2": "foobar_bar",
		}
		nodes := map[string]string{
			"taskID1": "foo1",
			"taskID2": "foo2",
		}
		out := bytes.NewBufferString("")
		testcase.context.Output = out
		err := TaskWrite(testcase.context, tasks, names, nodes)
		if err != nil {
			assert.Error(t, err, testcase.expected)
		} else {
			assert.Equal(t, out.String(), testcase.expected)
		}
	}
}

func TestTaskContextWriteJSONField(t *testing.T) {
	tasks := []swarm.Task{
		{ID: "taskID1"},
		{ID: "taskID2"},
	}
	names := map[string]string{
		"taskID1": "foobar_baz",
		"taskID2": "foobar_bar",
	}
	out := bytes.NewBufferString("")
	err := TaskWrite(Context{Format: "{{json .ID}}", Output: out}, tasks, names, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		var s string
		if err := json.Unmarshal([]byte(line), &s); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, s, tasks[i].ID)
	}
}
