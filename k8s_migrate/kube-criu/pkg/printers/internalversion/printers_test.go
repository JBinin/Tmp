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
/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internalversion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ghodss/yaml"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	genericprinters "k8s.io/cli-runtime/pkg/genericclioptions/printers"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/api/testapi"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/coordination"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/policy"
	"k8s.io/kubernetes/pkg/apis/storage"
	"k8s.io/kubernetes/pkg/printers"
)

func init() {
	legacyscheme.Scheme.AddKnownTypes(schema.GroupVersion{Group: "", Version: runtime.APIVersionInternal}, &TestPrintType{})
	legacyscheme.Scheme.AddKnownTypes(schema.GroupVersion{Group: "", Version: "v1"}, &TestPrintType{})
}

var testData = TestStruct{
	TypeMeta:   metav1.TypeMeta{APIVersion: "foo/bar", Kind: "TestStruct"},
	Key:        "testValue",
	Map:        map[string]int{"TestSubkey": 1},
	StringList: []string{"a", "b", "c"},
	IntList:    []int{1, 2, 3},
}

type TestStruct struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Key               string         `json:"Key"`
	Map               map[string]int `json:"Map"`
	StringList        []string       `json:"StringList"`
	IntList           []int          `json:"IntList"`
}

func (in *TestStruct) DeepCopyObject() runtime.Object {
	panic("never called")
}

func TestPrintUnstructuredObject(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Test",
			"dummy1":     "present",
			"dummy2":     "present",
			"metadata": map[string]interface{}{
				"name":              "MyName",
				"namespace":         "MyNamespace",
				"creationTimestamp": "2017-04-01T00:00:00Z",
				"resourceVersion":   123,
				"uid":               "00000000-0000-0000-0000-000000000001",
				"dummy3":            "present",
				"labels":            map[string]interface{}{"test": "other"},
			},
			/*"items": []interface{}{
				map[string]interface{}{
					"itemBool": true,
					"itemInt":  42,
				},
			},*/
			"url":    "http://localhost",
			"status": "ok",
		},
	}

	tests := []struct {
		expected string
		options  printers.PrintOptions
		object   runtime.Object
	}{
		{
			expected: "NAME\\s+AGE\nMyName\\s+\\d+",
			object:   obj,
		},
		{
			options: printers.PrintOptions{
				WithNamespace: true,
			},
			expected: "NAMESPACE\\s+NAME\\s+AGE\nMyNamespace\\s+MyName\\s+\\d+",
			object:   obj,
		},
		{
			options: printers.PrintOptions{
				ShowLabels:    true,
				WithNamespace: true,
			},
			expected: "NAMESPACE\\s+NAME\\s+AGE\\s+LABELS\nMyNamespace\\s+MyName\\s+\\d+\\w+\\s+test\\=other",
			object:   obj,
		},
		{
			expected: "NAME\\s+AGE\nMyName\\s+\\d+\\w+\nMyName2\\s+\\d+",
			object: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Test",
					"dummy1":     "present",
					"dummy2":     "present",
					"items": []interface{}{
						map[string]interface{}{
							"metadata": map[string]interface{}{
								"name":              "MyName",
								"namespace":         "MyNamespace",
								"creationTimestamp": "2017-04-01T00:00:00Z",
								"resourceVersion":   123,
								"uid":               "00000000-0000-0000-0000-000000000001",
								"dummy3":            "present",
								"labels":            map[string]interface{}{"test": "other"},
							},
						},
						map[string]interface{}{
							"metadata": map[string]interface{}{
								"name":              "MyName2",
								"namespace":         "MyNamespace",
								"creationTimestamp": "2017-04-01T00:00:00Z",
								"resourceVersion":   123,
								"uid":               "00000000-0000-0000-0000-000000000001",
								"dummy3":            "present",
								"labels":            "badlabel",
							},
						},
					},
					"url":    "http://localhost",
					"status": "ok",
				},
			},
		},
	}
	out := bytes.NewBuffer([]byte{})

	for _, test := range tests {
		out.Reset()
		printer := printers.NewHumanReadablePrinter(nil, test.options).With(AddDefaultHandlers)
		printer.PrintObj(test.object, out)

		matches, err := regexp.MatchString(test.expected, out.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !matches {
			t.Errorf("wanted:\n%s\ngot:\n%s", test.expected, out)
		}
	}
}

type TestPrintType struct {
	Data string
}

func (obj *TestPrintType) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }
func (obj *TestPrintType) DeepCopyObject() runtime.Object {
	if obj == nil {
		return nil
	}
	clone := *obj
	return &clone
}

type TestUnknownType struct{}

func (obj *TestUnknownType) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }
func (obj *TestUnknownType) DeepCopyObject() runtime.Object {
	if obj == nil {
		return nil
	}
	clone := *obj
	return &clone
}

func testPrinter(t *testing.T, printer printers.ResourcePrinter, unmarshalFunc func(data []byte, v interface{}) error) {
	buf := bytes.NewBuffer([]byte{})

	err := printer.PrintObj(&testData, buf)
	if err != nil {
		t.Fatal(err)
	}
	var poutput TestStruct
	// Verify that given function runs without error.
	err = unmarshalFunc(buf.Bytes(), &poutput)
	if err != nil {
		t.Fatal(err)
	}
	// Use real decode function to undo the versioning process.
	poutput = TestStruct{}
	s := yamlserializer.NewDecodingSerializer(testapi.Default.Codec())
	if err := runtime.DecodeInto(s, buf.Bytes(), &poutput); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(testData, poutput) {
		t.Errorf("Test data and unmarshaled data are not equal: %v", diff.ObjectDiff(poutput, testData))
	}

	obj := &v1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
	}
	buf.Reset()
	printer.PrintObj(obj, buf)
	var objOut v1.Pod
	// Verify that given function runs without error.
	err = unmarshalFunc(buf.Bytes(), &objOut)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}
	// Use real decode function to undo the versioning process.
	objOut = v1.Pod{}
	if err := runtime.DecodeInto(s, buf.Bytes(), &objOut); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(obj, &objOut) {
		t.Errorf("Unexpected inequality:\n%v", diff.ObjectDiff(obj, &objOut))
	}
}

func TestYAMLPrinter(t *testing.T) {
	testPrinter(t, genericprinters.NewTypeSetter(legacyscheme.Scheme).ToPrinter(&genericprinters.YAMLPrinter{}), yaml.Unmarshal)
}

func TestJSONPrinter(t *testing.T) {
	testPrinter(t, genericprinters.NewTypeSetter(legacyscheme.Scheme).ToPrinter(&genericprinters.JSONPrinter{}), json.Unmarshal)
}

func TestFormatResourceName(t *testing.T) {
	tests := []struct {
		kind schema.GroupKind
		name string
		want string
	}{
		{schema.GroupKind{}, "", ""},
		{schema.GroupKind{}, "name", "name"},
		{schema.GroupKind{Kind: "Kind"}, "", "kind/"}, // should not happen in practice
		{schema.GroupKind{Kind: "Kind"}, "name", "kind/name"},
		{schema.GroupKind{Group: "group", Kind: "Kind"}, "name", "kind.group/name"},
	}
	for _, tt := range tests {
		if got := printers.FormatResourceName(tt.kind, tt.name, true); got != tt.want {
			t.Errorf("formatResourceName(%q, %q) = %q, want %q", tt.kind, tt.name, got, tt.want)
		}
	}
}

func PrintCustomType(obj *TestPrintType, w io.Writer, options printers.PrintOptions) error {
	data := obj.Data
	kind := options.Kind
	if options.WithKind {
		data = kind.String() + "/" + data
	}
	_, err := fmt.Fprintf(w, "%s", data)
	return err
}

func ErrorPrintHandler(obj *TestPrintType, w io.Writer, options printers.PrintOptions) error {
	return fmt.Errorf("ErrorPrintHandler error")
}

func TestCustomTypePrinting(t *testing.T) {
	columns := []string{"Data"}
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{})
	printer.Handler(columns, nil, PrintCustomType)

	obj := TestPrintType{"test object"}
	buffer := &bytes.Buffer{}
	err := printer.PrintObj(&obj, buffer)
	if err != nil {
		t.Fatalf("An error occurred printing the custom type: %#v", err)
	}
	expectedOutput := "DATA\ntest object"
	if buffer.String() != expectedOutput {
		t.Errorf("The data was not printed as expected. Expected:\n%s\nGot:\n%s", expectedOutput, buffer.String())
	}
}

func TestPrintHandlerError(t *testing.T) {
	columns := []string{"Data"}
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{})
	printer.Handler(columns, nil, ErrorPrintHandler)
	obj := TestPrintType{"test object"}
	buffer := &bytes.Buffer{}
	err := printer.PrintObj(&obj, buffer)
	if err == nil || err.Error() != "ErrorPrintHandler error" {
		t.Errorf("Did not get the expected error: %#v", err)
	}
}

func TestUnknownTypePrinting(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{})
	buffer := &bytes.Buffer{}
	err := printer.PrintObj(&TestUnknownType{}, buffer)
	if err == nil {
		t.Errorf("An error was expected from printing unknown type")
	}
}

func TestTemplatePanic(t *testing.T) {
	tmpl := `{{and ((index .currentState.info "foo").state.running.startedAt) .currentState.info.net.state.running.startedAt}}`
	printer, err := genericprinters.NewGoTemplatePrinter([]byte(tmpl))
	if err != nil {
		t.Fatalf("tmpl fail: %v", err)
	}
	buffer := &bytes.Buffer{}
	err = printer.PrintObj(&v1.Pod{}, buffer)
	if err == nil {
		t.Fatalf("expected that template to crash")
	}
	if buffer.String() == "" {
		t.Errorf("no debugging info was printed")
	}
}

func TestNamePrinter(t *testing.T) {
	tests := map[string]struct {
		obj    runtime.Object
		expect string
	}{
		"singleObject": {
			&v1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind: "Pod",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
			"pod/foo\n"},
		"List": {
			&unstructured.UnstructuredList{
				Object: map[string]interface{}{
					"kind":       "List",
					"apiVersion": "v1",
				},
				Items: []unstructured.Unstructured{
					{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"metadata": map[string]interface{}{
								"name": "bar",
							},
						},
					},
				},
			},
			"pod/bar\n"},
	}

	printFlags := genericclioptions.NewPrintFlags("").WithTypeSetter(legacyscheme.Scheme).WithDefaultOutput("name")
	printer, err := printFlags.ToPrinter()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	for name, item := range tests {
		buff := &bytes.Buffer{}
		err := printer.PrintObj(item.obj, buff)
		if err != nil {
			t.Errorf("%v: unexpected err: %v", name, err)
			continue
		}
		got := buff.String()
		if item.expect != got {
			t.Errorf("%v: expected %v, got %v", name, item.expect, got)
		}
	}
}

func TestTemplateStrings(t *testing.T) {
	// This unit tests the "exists" function as well as the template from update.sh
	table := map[string]struct {
		pod    v1.Pod
		expect string
	}{
		"nilInfo":   {v1.Pod{}, "false"},
		"emptyInfo": {v1.Pod{Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{}}}, "false"},
		"fooExists": {
			v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "foo",
						},
					},
				},
			},
			"false",
		},
		"barExists": {
			v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "bar",
						},
					},
				},
			},
			"false",
		},
		"bothExist": {
			v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "foo",
						},
						{
							Name: "bar",
						},
					},
				},
			},
			"false",
		},
		"barValid": {
			v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "foo",
						},
						{
							Name: "bar",
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{
									StartedAt: metav1.Time{},
								},
							},
						},
					},
				},
			},
			"false",
		},
		"bothValid": {
			v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name: "foo",
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{
									StartedAt: metav1.Time{},
								},
							},
						},
						{
							Name: "bar",
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{
									StartedAt: metav1.Time{},
								},
							},
						},
					},
				},
			},
			"true",
		},
	}
	// The point of this test is to verify that the below template works.
	tmpl := `{{if (exists . "status" "containerStatuses")}}{{range .status.containerStatuses}}{{if (and (eq .name "foo") (exists . "state" "running"))}}true{{end}}{{end}}{{end}}`
	printer, err := genericprinters.NewGoTemplatePrinter([]byte(tmpl))
	if err != nil {
		t.Fatalf("tmpl fail: %v", err)
	}

	for name, item := range table {
		buffer := &bytes.Buffer{}
		err = printer.PrintObj(&item.pod, buffer)
		if err != nil {
			t.Errorf("%v: unexpected err: %v", name, err)
			continue
		}
		actual := buffer.String()
		if len(actual) == 0 {
			actual = "false"
		}
		if e := item.expect; e != actual {
			t.Errorf("%v: expected %v, got %v", name, e, actual)
		}
	}
}

func TestPrinters(t *testing.T) {
	om := func(name string) metav1.ObjectMeta { return metav1.ObjectMeta{Name: name} }

	var (
		err              error
		templatePrinter  printers.ResourcePrinter
		templatePrinter2 printers.ResourcePrinter
		jsonpathPrinter  printers.ResourcePrinter
	)

	templatePrinter, err = genericprinters.NewGoTemplatePrinter([]byte("{{.name}}"))
	if err != nil {
		t.Fatal(err)
	}

	templatePrinter2, err = genericprinters.NewGoTemplatePrinter([]byte("{{len .items}}"))
	if err != nil {
		t.Fatal(err)
	}

	jsonpathPrinter, err = genericprinters.NewJSONPathPrinter("{.metadata.name}")
	if err != nil {
		t.Fatal(err)
	}

	genericPrinters := map[string]printers.ResourcePrinter{
		// TODO(juanvallejo): move "generic printer" tests to pkg/kubectl/genericclioptions/printers
		"json":      genericprinters.NewTypeSetter(legacyscheme.Scheme).ToPrinter(&genericprinters.JSONPrinter{}),
		"yaml":      genericprinters.NewTypeSetter(legacyscheme.Scheme).ToPrinter(&genericprinters.YAMLPrinter{}),
		"template":  templatePrinter,
		"template2": templatePrinter2,
		"jsonpath":  jsonpathPrinter,
	}
	objects := map[string]runtime.Object{
		"pod":             &v1.Pod{ObjectMeta: om("pod")},
		"emptyPodList":    &v1.PodList{},
		"nonEmptyPodList": &v1.PodList{Items: []v1.Pod{{}}},
		"endpoints": &v1.Endpoints{
			Subsets: []v1.EndpointSubset{{
				Addresses: []v1.EndpointAddress{{IP: "127.0.0.1"}, {IP: "localhost"}},
				Ports:     []v1.EndpointPort{{Port: 8080}},
			}}},
	}
	// map of printer name to set of objects it should fail on.
	expectedErrors := map[string]sets.String{
		"template2": sets.NewString("pod", "emptyPodList", "endpoints"),
		"jsonpath":  sets.NewString("emptyPodList", "nonEmptyPodList", "endpoints"),
	}

	for pName, p := range genericPrinters {
		for oName, obj := range objects {
			b := &bytes.Buffer{}
			if err := p.PrintObj(obj, b); err != nil {
				if set, found := expectedErrors[pName]; found && set.Has(oName) {
					// expected error
					continue
				}
				t.Errorf("printer '%v', object '%v'; error: '%v'", pName, oName, err)
			}
		}
	}

	// a humanreadable printer deals with internal-versioned objects
	humanReadablePrinter := map[string]printers.ResourcePrinter{
		"humanReadable": printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
			NoHeaders: true,
		}),
		"humanReadableHeaders": printers.NewHumanReadablePrinter(nil, printers.PrintOptions{}),
	}
	AddHandlers((humanReadablePrinter["humanReadable"]).(*printers.HumanReadablePrinter))
	AddHandlers((humanReadablePrinter["humanReadableHeaders"]).(*printers.HumanReadablePrinter))
	for pName, p := range humanReadablePrinter {
		for oName, obj := range objects {
			b := &bytes.Buffer{}
			if err := p.PrintObj(obj, b); err != nil {
				if set, found := expectedErrors[pName]; found && set.Has(oName) {
					// expected error
					continue
				}
				t.Errorf("printer '%v', object '%v'; error: '%v'", pName, oName, err)
			}
		}
	}
}

func TestPrintEventsResultSorted(t *testing.T) {
	// Arrange
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{})
	AddHandlers(printer)

	obj := api.EventList{
		Items: []api.Event{
			{
				Source:         api.EventSource{Component: "kubelet"},
				Message:        "Item 1",
				FirstTimestamp: metav1.NewTime(time.Date(2014, time.January, 15, 0, 0, 0, 0, time.UTC)),
				LastTimestamp:  metav1.NewTime(time.Date(2014, time.January, 15, 0, 0, 0, 0, time.UTC)),
				Count:          1,
				Type:           api.EventTypeNormal,
			},
			{
				Source:         api.EventSource{Component: "scheduler"},
				Message:        "Item 2",
				FirstTimestamp: metav1.NewTime(time.Date(1987, time.June, 17, 0, 0, 0, 0, time.UTC)),
				LastTimestamp:  metav1.NewTime(time.Date(1987, time.June, 17, 0, 0, 0, 0, time.UTC)),
				Count:          1,
				Type:           api.EventTypeNormal,
			},
			{
				Source:         api.EventSource{Component: "kubelet"},
				Message:        "Item 3",
				FirstTimestamp: metav1.NewTime(time.Date(2002, time.December, 25, 0, 0, 0, 0, time.UTC)),
				LastTimestamp:  metav1.NewTime(time.Date(2002, time.December, 25, 0, 0, 0, 0, time.UTC)),
				Count:          1,
				Type:           api.EventTypeNormal,
			},
		},
	}
	buffer := &bytes.Buffer{}

	// Act
	err := printer.PrintObj(&obj, buffer)

	// Assert
	if err != nil {
		t.Fatalf("An error occurred printing the EventList: %#v", err)
	}
	out := buffer.String()
	VerifyDatesInOrder(out, "\n" /* rowDelimiter */, "  " /* columnDelimiter */, t)
}

func TestPrintNodeStatus(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{})
	AddHandlers(printer)
	table := []struct {
		node   api.Node
		status string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{Type: api.NodeReady, Status: api.ConditionTrue}}},
			},
			status: "Ready",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
				Spec:       api.NodeSpec{Unschedulable: true},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{Type: api.NodeReady, Status: api.ConditionTrue}}},
			},
			status: "Ready,SchedulingDisabled",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo3"},
				Status: api.NodeStatus{Conditions: []api.NodeCondition{
					{Type: api.NodeReady, Status: api.ConditionTrue},
					{Type: api.NodeReady, Status: api.ConditionTrue}}},
			},
			status: "Ready",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo4"},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{Type: api.NodeReady, Status: api.ConditionFalse}}},
			},
			status: "NotReady",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo5"},
				Spec:       api.NodeSpec{Unschedulable: true},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{Type: api.NodeReady, Status: api.ConditionFalse}}},
			},
			status: "NotReady,SchedulingDisabled",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo6"},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{Type: "InvalidValue", Status: api.ConditionTrue}}},
			},
			status: "Unknown",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo7"},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{}}},
			},
			status: "Unknown",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo8"},
				Spec:       api.NodeSpec{Unschedulable: true},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{Type: "InvalidValue", Status: api.ConditionTrue}}},
			},
			status: "Unknown,SchedulingDisabled",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo9"},
				Spec:       api.NodeSpec{Unschedulable: true},
				Status:     api.NodeStatus{Conditions: []api.NodeCondition{{}}},
			},
			status: "Unknown,SchedulingDisabled",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.status) {
			t.Fatalf("Expect printing node %s with status %#v, got: %#v", test.node.Name, test.status, buffer.String())
		}
	}
}

func TestPrintNodeRole(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{})
	AddHandlers(printer)
	table := []struct {
		node     api.Node
		expected string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo9"},
			},
			expected: "<none>",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "foo10",
					Labels: map[string]string{"node-role.kubernetes.io/master": "", "node-role.kubernetes.io/proxy": "", "kubernetes.io/role": "node"},
				},
			},
			expected: "master,node,proxy",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "foo11",
					Labels: map[string]string{"kubernetes.io/role": "node"},
				},
			},
			expected: "node",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.expected) {
			t.Fatalf("Expect printing node %s with role %#v, got: %#v", test.node.Name, test.expected, buffer.String())
		}
	}
}

func TestPrintNodeOSImage(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
		ColumnLabels: []string{},
		Wide:         true,
	})
	AddHandlers(printer)

	table := []struct {
		node    api.Node
		osImage string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
				Status: api.NodeStatus{
					NodeInfo:  api.NodeSystemInfo{OSImage: "fake-os-image"},
					Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}},
				},
			},
			osImage: "fake-os-image",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
				Status: api.NodeStatus{
					NodeInfo:  api.NodeSystemInfo{KernelVersion: "fake-kernel-version"},
					Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}},
				},
			},
			osImage: "<unknown>",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.osImage) {
			t.Fatalf("Expect printing node %s with os image %#v, got: %#v", test.node.Name, test.osImage, buffer.String())
		}
	}
}

func TestPrintNodeKernelVersion(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
		ColumnLabels: []string{},
		Wide:         true,
	})
	AddHandlers(printer)

	table := []struct {
		node          api.Node
		kernelVersion string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
				Status: api.NodeStatus{
					NodeInfo:  api.NodeSystemInfo{KernelVersion: "fake-kernel-version"},
					Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}},
				},
			},
			kernelVersion: "fake-kernel-version",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
				Status: api.NodeStatus{
					NodeInfo:  api.NodeSystemInfo{OSImage: "fake-os-image"},
					Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}},
				},
			},
			kernelVersion: "<unknown>",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.kernelVersion) {
			t.Fatalf("Expect printing node %s with kernel version %#v, got: %#v", test.node.Name, test.kernelVersion, buffer.String())
		}
	}
}

func TestPrintNodeContainerRuntimeVersion(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
		ColumnLabels: []string{},
		Wide:         true,
	})
	AddHandlers(printer)

	table := []struct {
		node                    api.Node
		containerRuntimeVersion string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
				Status: api.NodeStatus{
					NodeInfo:  api.NodeSystemInfo{ContainerRuntimeVersion: "foo://1.2.3"},
					Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}},
				},
			},
			containerRuntimeVersion: "foo://1.2.3",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
				Status: api.NodeStatus{
					NodeInfo:  api.NodeSystemInfo{},
					Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}},
				},
			},
			containerRuntimeVersion: "<unknown>",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.containerRuntimeVersion) {
			t.Fatalf("Expect printing node %s with kernel version %#v, got: %#v", test.node.Name, test.containerRuntimeVersion, buffer.String())
		}
	}
}

func TestPrintNodeName(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
		Wide: true,
	})
	AddHandlers(printer)
	table := []struct {
		node api.Node
		Name string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "127.0.0.1"},
				Status:     api.NodeStatus{},
			},
			Name: "127.0.0.1",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: ""},
				Status:     api.NodeStatus{},
			},
			Name: "<unknown>",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.Name) {
			t.Fatalf("Expect printing node %s with node name %#v, got: %#v", test.node.Name, test.Name, buffer.String())
		}
	}
}

func TestPrintNodeExternalIP(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
		Wide: true,
	})
	AddHandlers(printer)
	table := []struct {
		node       api.Node
		externalIP string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
				Status:     api.NodeStatus{Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}}},
			},
			externalIP: "1.1.1.1",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
				Status:     api.NodeStatus{Addresses: []api.NodeAddress{{Type: api.NodeInternalIP, Address: "1.1.1.1"}}},
			},
			externalIP: "<none>",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo3"},
				Status: api.NodeStatus{Addresses: []api.NodeAddress{
					{Type: api.NodeExternalIP, Address: "2.2.2.2"},
					{Type: api.NodeInternalIP, Address: "3.3.3.3"},
					{Type: api.NodeExternalIP, Address: "4.4.4.4"},
				}},
			},
			externalIP: "2.2.2.2",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.externalIP) {
			t.Fatalf("Expect printing node %s with external ip %#v, got: %#v", test.node.Name, test.externalIP, buffer.String())
		}
	}
}

func TestPrintNodeInternalIP(t *testing.T) {
	printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
		Wide: true,
	})
	AddHandlers(printer)
	table := []struct {
		node       api.Node
		internalIP string
	}{
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
				Status:     api.NodeStatus{Addresses: []api.NodeAddress{{Type: api.NodeInternalIP, Address: "1.1.1.1"}}},
			},
			internalIP: "1.1.1.1",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
				Status:     api.NodeStatus{Addresses: []api.NodeAddress{{Type: api.NodeExternalIP, Address: "1.1.1.1"}}},
			},
			internalIP: "<none>",
		},
		{
			node: api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "foo3"},
				Status: api.NodeStatus{Addresses: []api.NodeAddress{
					{Type: api.NodeInternalIP, Address: "2.2.2.2"},
					{Type: api.NodeExternalIP, Address: "3.3.3.3"},
					{Type: api.NodeInternalIP, Address: "4.4.4.4"},
				}},
			},
			internalIP: "2.2.2.2",
		},
	}

	for _, test := range table {
		buffer := &bytes.Buffer{}
		err := printer.PrintObj(&test.node, buffer)
		if err != nil {
			t.Fatalf("An error occurred printing Node: %#v", err)
		}
		if !contains(strings.Fields(buffer.String()), test.internalIP) {
			t.Fatalf("Expect printing node %s with internal ip %#v, got: %#v", test.node.Name, test.internalIP, buffer.String())
		}
	}
}

func contains(fields []string, field string) bool {
	for _, v := range fields {
		if v == field {
			return true
		}
	}
	return false
}

func TestPrintHunmanReadableIngressWithColumnLabels(t *testing.T) {
	ingress := extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test1",
			CreationTimestamp: metav1.Time{Time: time.Now().AddDate(-10, 0, 0)},
			Labels: map[string]string{
				"app_name": "kubectl_test_ingress",
			},
		},
		Spec: extensions.IngressSpec{
			Backend: &extensions.IngressBackend{
				ServiceName: "svc",
				ServicePort: intstr.FromInt(93),
			},
		},
		Status: extensions.IngressStatus{
			LoadBalancer: api.LoadBalancerStatus{
				Ingress: []api.LoadBalancerIngress{
					{
						IP:       "2.3.4.5",
						Hostname: "localhost.localdomain",
					},
				},
			},
		},
	}
	buff := bytes.NewBuffer([]byte{})
	table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&ingress, printers.PrintOptions{ColumnLabels: []string{"app_name"}})
	if err != nil {
		t.Fatal(err)
	}
	verifyTable(t, table)
	if err := printers.PrintTable(table, buff, printers.PrintOptions{NoHeaders: true}); err != nil {
		t.Fatal(err)
	}
	output := string(buff.Bytes())
	appName := ingress.ObjectMeta.Labels["app_name"]
	if !strings.Contains(output, appName) {
		t.Errorf("expected to container app_name label value %s, but doesn't %s", appName, output)
	}
}

func TestPrintHumanReadableService(t *testing.T) {
	tests := []api.Service{
		{
			Spec: api.ServiceSpec{
				ClusterIP: "1.2.3.4",
				Type:      "LoadBalancer",
				Ports: []api.ServicePort{
					{
						Port:     80,
						Protocol: "TCP",
					},
				},
			},
			Status: api.ServiceStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{
							IP: "2.3.4.5",
						},
						{
							IP: "3.4.5.6",
						},
					},
				},
			},
		},
		{
			Spec: api.ServiceSpec{
				ClusterIP: "1.3.4.5",
				Ports: []api.ServicePort{
					{
						Port:     80,
						Protocol: "TCP",
					},
					{
						Port:     8090,
						Protocol: "UDP",
					},
					{
						Port:     8000,
						Protocol: "TCP",
					},
				},
			},
		},
		{
			Spec: api.ServiceSpec{
				ClusterIP: "1.4.5.6",
				Type:      "LoadBalancer",
				Ports: []api.ServicePort{
					{
						Port:     80,
						Protocol: "TCP",
					},
					{
						Port:     8090,
						Protocol: "UDP",
					},
					{
						Port:     8000,
						Protocol: "TCP",
					},
				},
			},
			Status: api.ServiceStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{
							IP: "2.3.4.5",
						},
					},
				},
			},
		},
		{
			Spec: api.ServiceSpec{
				ClusterIP: "1.5.6.7",
				Type:      "LoadBalancer",
				Ports: []api.ServicePort{
					{
						Port:     80,
						Protocol: "TCP",
					},
					{
						Port:     8090,
						Protocol: "UDP",
					},
					{
						Port:     8000,
						Protocol: "TCP",
					},
				},
			},
			Status: api.ServiceStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{
							IP: "2.3.4.5",
						},
						{
							IP: "3.4.5.6",
						},
						{
							IP:       "5.6.7.8",
							Hostname: "host5678",
						},
					},
				},
			},
		},
	}

	for _, svc := range tests {
		for _, wide := range []bool{false, true} {
			buff := bytes.NewBuffer([]byte{})
			table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&svc, printers.PrintOptions{Wide: wide})
			if err != nil {
				t.Fatal(err)
			}
			verifyTable(t, table)
			if err := printers.PrintTable(table, buff, printers.PrintOptions{NoHeaders: true}); err != nil {
				t.Fatal(err)
			}
			output := string(buff.Bytes())
			ip := svc.Spec.ClusterIP
			if !strings.Contains(output, ip) {
				t.Errorf("expected to contain ClusterIP %s, but doesn't: %s", ip, output)
			}

			for n, ingress := range svc.Status.LoadBalancer.Ingress {
				ip = ingress.IP
				// For non-wide output, we only guarantee the first IP to be printed
				if (n == 0 || wide) && !strings.Contains(output, ip) {
					t.Errorf("expected to contain ingress ip %s with wide=%v, but doesn't: %s", ip, wide, output)
				}
			}

			for _, port := range svc.Spec.Ports {
				portSpec := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
				if !strings.Contains(output, portSpec) {
					t.Errorf("expected to contain port: %s, but doesn't: %s", portSpec, output)
				}
			}
			// Each service should print on one line
			if 1 != strings.Count(output, "\n") {
				t.Errorf("expected a single newline, found %d", strings.Count(output, "\n"))
			}
		}
	}
}

func TestPrintHumanReadableWithNamespace(t *testing.T) {
	namespaceName := "testnamespace"
	name := "test"
	table := []struct {
		obj          runtime.Object
		isNamespaced bool
	}{
		{
			obj: &api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
			},
			isNamespaced: true,
		},
		{
			obj: &api.ReplicationController{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Spec: api.ReplicationControllerSpec{
					Replicas: 2,
					Template: &api.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"name": "foo",
								"type": "production",
							},
						},
						Spec: api.PodSpec{
							Containers: []api.Container{
								{
									Image: "foo/bar",
									TerminationMessagePath: api.TerminationMessagePathDefault,
									ImagePullPolicy:        api.PullIfNotPresent,
								},
							},
							RestartPolicy: api.RestartPolicyAlways,
							DNSPolicy:     api.DNSDefault,
							NodeSelector: map[string]string{
								"baz": "blah",
							},
						},
					},
				},
			},
			isNamespaced: true,
		},
		{
			obj: &api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Spec: api.ServiceSpec{
					ClusterIP: "1.2.3.4",
					Ports: []api.ServicePort{
						{
							Port:     80,
							Protocol: "TCP",
						},
					},
				},
				Status: api.ServiceStatus{
					LoadBalancer: api.LoadBalancerStatus{
						Ingress: []api.LoadBalancerIngress{
							{
								IP: "2.3.4.5",
							},
						},
					},
				},
			},
			isNamespaced: true,
		},
		{
			obj: &api.Endpoints{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Subsets: []api.EndpointSubset{{
					Addresses: []api.EndpointAddress{{IP: "127.0.0.1"}, {IP: "localhost"}},
					Ports:     []api.EndpointPort{{Port: 8080}},
				},
				}},
			isNamespaced: true,
		},
		{
			obj: &api.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: name},
			},
			isNamespaced: false,
		},
		{
			obj: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
			},
			isNamespaced: true,
		},
		{
			obj: &api.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Secrets:    []api.ObjectReference{},
			},
			isNamespaced: true,
		},
		{
			obj: &api.Node{
				ObjectMeta: metav1.ObjectMeta{Name: name},
				Status:     api.NodeStatus{},
			},
			isNamespaced: false,
		},
		{
			obj: &api.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Spec:       api.PersistentVolumeSpec{},
			},
			isNamespaced: false,
		},
		{
			obj: &api.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Spec:       api.PersistentVolumeClaimSpec{},
			},
			isNamespaced: true,
		},
		{
			obj: &api.Event{
				ObjectMeta:     metav1.ObjectMeta{Name: name, Namespace: namespaceName},
				Source:         api.EventSource{Component: "kubelet"},
				Message:        "Item 1",
				FirstTimestamp: metav1.NewTime(time.Date(2014, time.January, 15, 0, 0, 0, 0, time.UTC)),
				LastTimestamp:  metav1.NewTime(time.Date(2014, time.January, 15, 0, 0, 0, 0, time.UTC)),
				Count:          1,
				Type:           api.EventTypeNormal,
			},
			isNamespaced: true,
		},
		{
			obj: &api.LimitRange{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
			},
			isNamespaced: true,
		},
		{
			obj: &api.ResourceQuota{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceName},
			},
			isNamespaced: true,
		},
		{
			obj: &api.ComponentStatus{
				Conditions: []api.ComponentCondition{
					{Type: api.ComponentHealthy, Status: api.ConditionTrue, Message: "ok", Error: ""},
				},
			},
			isNamespaced: false,
		},
	}

	for i, test := range table {
		if test.isNamespaced {
			// Expect output to include namespace when requested.
			printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
				WithNamespace: true,
			})
			AddHandlers(printer)
			buffer := &bytes.Buffer{}
			err := printer.PrintObj(test.obj, buffer)
			if err != nil {
				t.Fatalf("An error occurred printing object: %#v", err)
			}
			matched := contains(strings.Fields(buffer.String()), fmt.Sprintf("%s", namespaceName))
			if !matched {
				t.Errorf("%d: Expect printing object to contain namespace: %#v", i, test.obj)
			}
		} else {
			// Expect error when trying to get all namespaces for un-namespaced object.
			printer := printers.NewHumanReadablePrinter(nil, printers.PrintOptions{
				WithNamespace: true,
			})
			buffer := &bytes.Buffer{}
			err := printer.PrintObj(test.obj, buffer)
			if err == nil {
				t.Errorf("Expected error when printing un-namespaced type")
			}
		}
	}
}

func TestPrintPodTable(t *testing.T) {
	runningPod := &api.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test1", Labels: map[string]string{"a": "1", "b": "2"}},
		Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
		Status: api.PodStatus{
			Phase: "Running",
			ContainerStatuses: []api.ContainerStatus{
				{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
				{RestartCount: 3},
			},
		},
	}
	failedPod := &api.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test2", Labels: map[string]string{"b": "2"}},
		Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
		Status: api.PodStatus{
			Phase: "Failed",
			ContainerStatuses: []api.ContainerStatus{
				{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
				{RestartCount: 3},
			},
		},
	}
	tests := []struct {
		obj          runtime.Object
		opts         printers.PrintOptions
		expect       string
		ignoreLegacy bool
	}{
		{
			obj: runningPod, opts: printers.PrintOptions{},
			expect: "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\ntest1\t1/2\tRunning\t6\t<unknown>\n",
		},
		{
			obj: runningPod, opts: printers.PrintOptions{WithKind: true, Kind: schema.GroupKind{Kind: "Pod"}},
			expect: "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\npod/test1\t1/2\tRunning\t6\t<unknown>\n",
		},
		{
			obj: runningPod, opts: printers.PrintOptions{ShowLabels: true},
			expect: "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\tLABELS\ntest1\t1/2\tRunning\t6\t<unknown>\ta=1,b=2\n",
		},
		{
			obj: &api.PodList{Items: []api.Pod{*runningPod, *failedPod}}, opts: printers.PrintOptions{ColumnLabels: []string{"a"}},
			expect: "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\tA\ntest1\t1/2\tRunning\t6\t<unknown>\t1\ntest2\t1/2\tFailed\t6\t<unknown>\t\n",
		},
		{
			obj: runningPod, opts: printers.PrintOptions{NoHeaders: true},
			expect: "test1\t1/2\tRunning\t6\t<unknown>\n",
		},
		{
			obj: failedPod, opts: printers.PrintOptions{},
			expect:       "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\ntest2\t1/2\tFailed\t6\t<unknown>\n",
			ignoreLegacy: true, // filtering is not done by the printer in the legacy path
		},
		{
			obj: failedPod, opts: printers.PrintOptions{},
			expect: "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\ntest2\t1/2\tFailed\t6\t<unknown>\n",
		},
	}

	for i, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(test.obj, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		buf := &bytes.Buffer{}
		p := printers.NewHumanReadablePrinter(nil, test.opts).With(AddHandlers).AddTabWriter(false)
		if err := p.PrintObj(table, buf); err != nil {
			t.Fatal(err)
		}
		if test.expect != buf.String() {
			t.Errorf("%d mismatch:\n%s\n%s", i, strconv.Quote(test.expect), strconv.Quote(buf.String()))
		}
		if test.ignoreLegacy {
			continue
		}

		buf.Reset()
		if err := p.PrintObj(test.obj, buf); err != nil {
			t.Fatal(err)
		}
		if test.expect != buf.String() {
			t.Errorf("%d legacy mismatch:\n%s\n%s", i, strconv.Quote(test.expect), strconv.Quote(buf.String()))
		}
	}
}

func TestPrintPod(t *testing.T) {
	tests := []struct {
		pod    api.Pod
		expect []metav1beta1.TableRow
	}{
		{
			// Test name, num of containers, restarts, container ready status
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test1"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "podPhase", int64(6), "<unknown>"}}},
		},
		{
			// Test container error overwrites pod phase
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test2"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{State: api.ContainerState{Waiting: &api.ContainerStateWaiting{Reason: "ContainerWaitingReason"}}, RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test2", "1/2", "ContainerWaitingReason", int64(6), "<unknown>"}}},
		},
		{
			// Test the same as the above but with Terminated state and the first container overwrites the rest
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test3"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{State: api.ContainerState{Waiting: &api.ContainerStateWaiting{Reason: "ContainerWaitingReason"}}, RestartCount: 3},
						{State: api.ContainerState{Terminated: &api.ContainerStateTerminated{Reason: "ContainerTerminatedReason"}}, RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test3", "0/2", "ContainerWaitingReason", int64(6), "<unknown>"}}},
		},
		{
			// Test ready is not enough for reporting running
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test4"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{Ready: true, RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test4", "1/2", "podPhase", int64(6), "<unknown>"}}},
		},
		{
			// Test ready is not enough for reporting running
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test5"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Reason: "podReason",
					Phase:  "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{Ready: true, RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test5", "1/2", "podReason", int64(6), "<unknown>"}}},
		},
		{
			// Test pod has 2 containers, one is running and the other is completed.
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test6"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase:  "Running",
					Reason: "",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Terminated: &api.ContainerStateTerminated{Reason: "Completed", ExitCode: 0}}},
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test6", "1/2", "Running", int64(6), "<unknown>"}}},
		},
	}

	for i, test := range tests {
		rows, err := printPod(&test.pod, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		for i := range rows {
			rows[i].Object.Object = nil
		}
		if !reflect.DeepEqual(test.expect, rows) {
			t.Errorf("%d mismatch: %s", i, diff.ObjectReflectDiff(test.expect, rows))
		}
	}
}

func TestPrintPodwide(t *testing.T) {
	tests := []struct {
		pod    api.Pod
		expect []metav1beta1.TableRow
	}{
		{
			// Test when the NodeName and PodIP are not none
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test1"},
				Spec: api.PodSpec{
					Containers: make([]api.Container, 2),
					NodeName:   "test1",
				},
				Status: api.PodStatus{
					Phase: "podPhase",
					PodIP: "1.1.1.1",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
					NominatedNodeName: "node1",
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "podPhase", int64(6), "<unknown>", "1.1.1.1", "test1", "node1"}}},
		},
		{
			// Test when the NodeName and PodIP are none
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test2"},
				Spec: api.PodSpec{
					Containers: make([]api.Container, 2),
					NodeName:   "",
				},
				Status: api.PodStatus{
					Phase: "podPhase",
					PodIP: "",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{State: api.ContainerState{Waiting: &api.ContainerStateWaiting{Reason: "ContainerWaitingReason"}}, RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test2", "1/2", "ContainerWaitingReason", int64(6), "<unknown>", "<none>", "<none>", "<none>"}}},
		},
	}

	for i, test := range tests {
		rows, err := printPod(&test.pod, printers.PrintOptions{Wide: true})
		if err != nil {
			t.Fatal(err)
		}
		for i := range rows {
			rows[i].Object.Object = nil
		}
		if !reflect.DeepEqual(test.expect, rows) {
			t.Errorf("%d mismatch: %s", i, diff.ObjectReflectDiff(test.expect, rows))
		}
	}
}

func TestPrintPodList(t *testing.T) {
	tests := []struct {
		pods   api.PodList
		expect []metav1beta1.TableRow
	}{
		// Test podList's pod: name, num of containers, restarts, container ready status
		{
			api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "test1"},
						Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
						Status: api.PodStatus{
							Phase: "podPhase",
							ContainerStatuses: []api.ContainerStatus{
								{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
								{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "test2"},
						Spec:       api.PodSpec{Containers: make([]api.Container, 1)},
						Status: api.PodStatus{
							Phase: "podPhase",
							ContainerStatuses: []api.ContainerStatus{
								{Ready: true, RestartCount: 1, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
							},
						},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "2/2", "podPhase", int64(6), "<unknown>"}}, {Cells: []interface{}{"test2", "1/1", "podPhase", int64(1), "<unknown>"}}},
		},
	}

	for _, test := range tests {
		rows, err := printPodList(&test.pods, printers.PrintOptions{})

		if err != nil {
			t.Fatal(err)
		}
		for i := range rows {
			rows[i].Object.Object = nil
		}
		if !reflect.DeepEqual(test.expect, rows) {
			t.Errorf("mismatch: %s", diff.ObjectReflectDiff(test.expect, rows))
		}
	}
}

func TestPrintNonTerminatedPod(t *testing.T) {
	tests := []struct {
		pod    api.Pod
		expect []metav1beta1.TableRow
	}{
		{
			// Test pod phase Running should be printed
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test1"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: api.PodRunning,
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "Running", int64(6), "<unknown>"}}},
		},
		{
			// Test pod phase Pending should be printed
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test2"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: api.PodPending,
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test2", "1/2", "Pending", int64(6), "<unknown>"}}},
		},
		{
			// Test pod phase Unknown should be printed
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test3"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: api.PodUnknown,
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test3", "1/2", "Unknown", int64(6), "<unknown>"}}},
		},
		{
			// Test pod phase Succeeded shouldn't be printed
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test4"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: api.PodSucceeded,
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test4", "1/2", "Succeeded", int64(6), "<unknown>"}, Conditions: podSuccessConditions}},
		},
		{
			// Test pod phase Failed shouldn't be printed
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test5"},
				Spec:       api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: api.PodFailed,
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{Ready: true, RestartCount: 3},
					},
				},
			},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test5", "1/2", "Failed", int64(6), "<unknown>"}, Conditions: podFailedConditions}},
		},
	}

	for i, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.pod, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		rows := table.Rows
		for i := range rows {
			rows[i].Object.Object = nil
		}
		if !reflect.DeepEqual(test.expect, rows) {
			t.Errorf("%d mismatch: %s", i, diff.ObjectReflectDiff(test.expect, rows))
		}
	}
}

func TestPrintPodWithLabels(t *testing.T) {
	tests := []struct {
		pod          api.Pod
		labelColumns []string
		expect       []metav1beta1.TableRow
	}{
		{
			// Test name, num of containers, restarts, container ready status
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "test1",
					Labels: map[string]string{"col1": "asd", "COL2": "zxc"},
				},
				Spec: api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]string{"col1", "COL2"},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "podPhase", int64(6), "<unknown>", "asd", "zxc"}}},
		},
		{
			// Test name, num of containers, restarts, container ready status
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "test1",
					Labels: map[string]string{"col1": "asd", "COL2": "zxc"},
				},
				Spec: api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			[]string{},
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "podPhase", int64(6), "<unknown>"}}},
		},
	}

	for i, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.pod, printers.PrintOptions{ColumnLabels: test.labelColumns})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		rows := table.Rows
		for i := range rows {
			rows[i].Object.Object = nil
		}
		if !reflect.DeepEqual(test.expect, rows) {
			t.Errorf("%d mismatch: %s", i, diff.ObjectReflectDiff(test.expect, rows))
		}
	}
}

type stringTestList []struct {
	name, got, exp string
}

func TestTranslateTimestampSince(t *testing.T) {
	tl := stringTestList{
		{"a while from now", translateTimestampSince(metav1.Time{Time: time.Now().Add(2.1e9)}), "<invalid>"},
		{"almost now", translateTimestampSince(metav1.Time{Time: time.Now().Add(1.9e9)}), "0s"},
		{"now", translateTimestampSince(metav1.Time{Time: time.Now()}), "0s"},
		{"unknown", translateTimestampSince(metav1.Time{}), "<unknown>"},
		{"30 seconds ago", translateTimestampSince(metav1.Time{Time: time.Now().Add(-3e10)}), "30s"},
		{"5 minutes ago", translateTimestampSince(metav1.Time{Time: time.Now().Add(-3e11)}), "5m"},
		{"an hour ago", translateTimestampSince(metav1.Time{Time: time.Now().Add(-6e12)}), "100m"},
		{"2 days ago", translateTimestampSince(metav1.Time{Time: time.Now().UTC().AddDate(0, 0, -2)}), "2d"},
		{"months ago", translateTimestampSince(metav1.Time{Time: time.Now().UTC().AddDate(0, 0, -90)}), "90d"},
		{"10 years ago", translateTimestampSince(metav1.Time{Time: time.Now().UTC().AddDate(-10, 0, 0)}), "10y"},
	}
	for _, test := range tl {
		if test.got != test.exp {
			t.Errorf("On %v, expected '%v', but got '%v'",
				test.name, test.exp, test.got)
		}
	}
}

func TestTranslateTimestampUntil(t *testing.T) {
	// Since this method compares the time with time.Now() internally,
	// small buffers of 0.1 seconds are added on comparing times to consider method call overhead.
	// Otherwise, the output strings become shorter than expected.
	const buf = 1e8
	tl := stringTestList{
		{"a while ago", translateTimestampUntil(metav1.Time{Time: time.Now().Add(-2.1e9)}), "<invalid>"},
		{"almost now", translateTimestampUntil(metav1.Time{Time: time.Now().Add(-1.9e9)}), "0s"},
		{"now", translateTimestampUntil(metav1.Time{Time: time.Now()}), "0s"},
		{"unknown", translateTimestampUntil(metav1.Time{}), "<unknown>"},
		{"in 30 seconds", translateTimestampUntil(metav1.Time{Time: time.Now().Add(3e10 + buf)}), "30s"},
		{"in 5 minutes", translateTimestampUntil(metav1.Time{Time: time.Now().Add(3e11 + buf)}), "5m"},
		{"in an hour", translateTimestampUntil(metav1.Time{Time: time.Now().Add(6e12 + buf)}), "100m"},
		{"in 2 days", translateTimestampUntil(metav1.Time{Time: time.Now().UTC().AddDate(0, 0, 2).Add(buf)}), "2d"},
		{"in months", translateTimestampUntil(metav1.Time{Time: time.Now().UTC().AddDate(0, 0, 90).Add(buf)}), "90d"},
		{"in 10 years", translateTimestampUntil(metav1.Time{Time: time.Now().UTC().AddDate(10, 0, 0).Add(buf)}), "10y"},
	}
	for _, test := range tl {
		if test.got != test.exp {
			t.Errorf("On %v, expected '%v', but got '%v'",
				test.name, test.exp, test.got)
		}
	}
}

func TestPrintDeployment(t *testing.T) {
	tests := []struct {
		deployment extensions.Deployment
		expect     string
		wideExpect string
	}{
		{
			extensions.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: extensions.DeploymentSpec{
					Replicas: 5,
					Template: api.PodTemplateSpec{
						Spec: api.PodSpec{
							Containers: []api.Container{
								{
									Name:  "fake-container1",
									Image: "fake-image1",
								},
								{
									Name:  "fake-container2",
									Image: "fake-image2",
								},
							},
						},
					},
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
				Status: extensions.DeploymentStatus{
					Replicas:            10,
					UpdatedReplicas:     2,
					AvailableReplicas:   1,
					UnavailableReplicas: 4,
				},
			},
			"test1\t5\t10\t2\t1\t0s\n",
			"test1\t5\t10\t2\t1\t0s\tfake-container1,fake-container2\tfake-image1,fake-image2\tfoo=bar\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.deployment, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()
		table, err = printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.deployment, printers.PrintOptions{Wide: true})
		verifyTable(t, table)
		// print deployment with '-o wide' option
		if err := printers.PrintTable(table, buf, printers.PrintOptions{Wide: true, NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.wideExpect {
			t.Fatalf("Expected: %s, got: %s", test.wideExpect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintDaemonSet(t *testing.T) {
	tests := []struct {
		ds         extensions.DaemonSet
		startsWith string
	}{
		{
			extensions.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: extensions.DaemonSetSpec{
					Template: api.PodTemplateSpec{
						Spec: api.PodSpec{Containers: make([]api.Container, 2)},
					},
				},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 2,
					DesiredNumberScheduled: 3,
					NumberReady:            1,
					UpdatedNumberScheduled: 2,
					NumberAvailable:        0,
				},
			},
			"test1\t3\t2\t1\t2\t0\t<none>\t0s\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.ds, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(buf.String(), test.startsWith) {
			t.Fatalf("Expected to start with %s but got %s", test.startsWith, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintJob(t *testing.T) {
	now := time.Now()
	completions := int32(2)
	tests := []struct {
		job    batch.Job
		expect string
	}{
		{
			batch.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "job1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: batch.JobSpec{
					Completions: &completions,
				},
				Status: batch.JobStatus{
					Succeeded: 1,
				},
			},
			"job1\t1/2\t\t0s\n",
		},
		{
			batch.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "job2",
					CreationTimestamp: metav1.Time{Time: time.Now().AddDate(-10, 0, 0)},
				},
				Spec: batch.JobSpec{
					Completions: nil,
				},
				Status: batch.JobStatus{
					Succeeded: 0,
				},
			},
			"job2\t0/1\t\t10y\n",
		},
		{
			batch.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "job3",
					CreationTimestamp: metav1.Time{Time: time.Now().AddDate(-10, 0, 0)},
				},
				Spec: batch.JobSpec{
					Completions: nil,
				},
				Status: batch.JobStatus{
					Succeeded:      0,
					StartTime:      &metav1.Time{Time: now.Add(time.Minute)},
					CompletionTime: &metav1.Time{Time: now.Add(31 * time.Minute)},
				},
			},
			"job3\t0/1\t30m\t10y\n",
		},
		{
			batch.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "job4",
					CreationTimestamp: metav1.Time{Time: time.Now().AddDate(-10, 0, 0)},
				},
				Spec: batch.JobSpec{
					Completions: nil,
				},
				Status: batch.JobStatus{
					Succeeded: 0,
					StartTime: &metav1.Time{Time: time.Now().Add(-20 * time.Minute)},
				},
			},
			"job4\t0/1\t20m\t10y\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.job, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintHPA(t *testing.T) {
	minReplicasVal := int32(2)
	targetUtilizationVal := int32(80)
	currentUtilizationVal := int32(50)
	tests := []struct {
		hpa      autoscaling.HorizontalPodAutoscaler
		expected string
	}{
		// minReplicas unset
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MaxReplicas: 10,
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<none>\t<unset>\t10\t4\t<unknown>\n",
		},
		// external source type, target average value (no current)
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								MetricName:         "some-external-metric",
								TargetAverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<unknown>/100m (avg)\t2\t10\t4\t<unknown>\n",
		},
		// external source type, target average value
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								MetricName:         "some-external-metric",
								TargetAverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricStatus{
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								MetricName:          "some-external-metric",
								CurrentAverageValue: resource.NewMilliQuantity(50, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50m/100m (avg)\t2\t10\t4\t<unknown>\n",
		},
		// external source type, target value (no current)
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								MetricName:  "some-service-metric",
								TargetValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<unknown>/100m\t2\t10\t4\t<unknown>\n",
		},
		// external source type, target value
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricSource{
								MetricSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"label": "value",
									},
								},
								MetricName:  "some-external-metric",
								TargetValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.ExternalMetricSourceType,
							External: &autoscaling.ExternalMetricStatus{
								MetricName:   "some-external-metric",
								CurrentValue: *resource.NewMilliQuantity(50, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50m/100m\t2\t10\t4\t<unknown>\n",
		},
		// pods source type (no current)
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricSource{
								MetricName:         "some-pods-metric",
								TargetAverageValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<unknown>/100m\t2\t10\t4\t<unknown>\n",
		},
		// pods source type
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricSource{
								MetricName:         "some-pods-metric",
								TargetAverageValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricStatus{
								MetricName:          "some-pods-metric",
								CurrentAverageValue: *resource.NewMilliQuantity(50, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50m/100m\t2\t10\t4\t<unknown>\n",
		},
		// object source type (no current)
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ObjectMetricSourceType,
							Object: &autoscaling.ObjectMetricSource{
								Target: autoscaling.CrossVersionObjectReference{
									Name: "some-service",
									Kind: "Service",
								},
								MetricName:  "some-service-metric",
								TargetValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<unknown>/100m\t2\t10\t4\t<unknown>\n",
		},
		// object source type
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ObjectMetricSourceType,
							Object: &autoscaling.ObjectMetricSource{
								Target: autoscaling.CrossVersionObjectReference{
									Name: "some-service",
									Kind: "Service",
								},
								MetricName:  "some-service-metric",
								TargetValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.ObjectMetricSourceType,
							Object: &autoscaling.ObjectMetricStatus{
								Target: autoscaling.CrossVersionObjectReference{
									Name: "some-service",
									Kind: "Service",
								},
								MetricName:   "some-service-metric",
								CurrentValue: *resource.NewMilliQuantity(50, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50m/100m\t2\t10\t4\t<unknown>\n",
		},
		// resource source type, targetVal (no current)
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name:               api.ResourceCPU,
								TargetAverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<unknown>/100m\t2\t10\t4\t<unknown>\n",
		},
		// resource source type, targetVal
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name:               api.ResourceCPU,
								TargetAverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricStatus{
								Name:                api.ResourceCPU,
								CurrentAverageValue: *resource.NewMilliQuantity(50, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50m/100m\t2\t10\t4\t<unknown>\n",
		},
		// resource source type, targetUtil (no current)
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: &targetUtilizationVal,
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
				},
			},
			"some-hpa\tReplicationController/some-rc\t<unknown>/80%\t2\t10\t4\t<unknown>\n",
		},
		// resource source type, targetUtil
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: &targetUtilizationVal,
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricStatus{
								Name: api.ResourceCPU,
								CurrentAverageUtilization: &currentUtilizationVal,
								CurrentAverageValue:       *resource.NewMilliQuantity(40, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50%/80%\t2\t10\t4\t<unknown>\n",
		},
		// multiple specs
		{
			autoscaling.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Name: "some-hpa"},
				Spec: autoscaling.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscaling.CrossVersionObjectReference{
						Name: "some-rc",
						Kind: "ReplicationController",
					},
					MinReplicas: &minReplicasVal,
					MaxReplicas: 10,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricSource{
								MetricName:         "some-pods-metric",
								TargetAverageValue: *resource.NewMilliQuantity(100, resource.DecimalSI),
							},
						},
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name: api.ResourceCPU,
								TargetAverageUtilization: &targetUtilizationVal,
							},
						},
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricSource{
								MetricName:         "other-pods-metric",
								TargetAverageValue: *resource.NewMilliQuantity(400, resource.DecimalSI),
							},
						},
					},
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					CurrentReplicas: 4,
					DesiredReplicas: 5,
					CurrentMetrics: []autoscaling.MetricStatus{
						{
							Type: autoscaling.PodsMetricSourceType,
							Pods: &autoscaling.PodsMetricStatus{
								MetricName:          "some-pods-metric",
								CurrentAverageValue: *resource.NewMilliQuantity(50, resource.DecimalSI),
							},
						},
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricStatus{
								Name: api.ResourceCPU,
								CurrentAverageUtilization: &currentUtilizationVal,
								CurrentAverageValue:       *resource.NewMilliQuantity(40, resource.DecimalSI),
							},
						},
					},
				},
			},
			"some-hpa\tReplicationController/some-rc\t50m/100m, 50%/80% + 1 more...\t2\t10\t4\t<unknown>\n",
		},
	}

	buff := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.hpa, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buff, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buff.String() != test.expected {
			t.Errorf("expected %q, got %q", test.expected, buff.String())
		}

		buff.Reset()
	}
}

func TestPrintPodShowLabels(t *testing.T) {
	tests := []struct {
		pod        api.Pod
		showLabels bool
		expect     []metav1beta1.TableRow
	}{
		{
			// Test name, num of containers, restarts, container ready status
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "test1",
					Labels: map[string]string{"col1": "asd", "COL2": "zxc"},
				},
				Spec: api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			true,
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "podPhase", int64(6), "<unknown>", "COL2=zxc,col1=asd"}}},
		},
		{
			// Test name, num of containers, restarts, container ready status
			api.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "test1",
					Labels: map[string]string{"col3": "asd", "COL4": "zxc"},
				},
				Spec: api.PodSpec{Containers: make([]api.Container, 2)},
				Status: api.PodStatus{
					Phase: "podPhase",
					ContainerStatuses: []api.ContainerStatus{
						{Ready: true, RestartCount: 3, State: api.ContainerState{Running: &api.ContainerStateRunning{}}},
						{RestartCount: 3},
					},
				},
			},
			false,
			[]metav1beta1.TableRow{{Cells: []interface{}{"test1", "1/2", "podPhase", int64(6), "<unknown>"}}},
		},
	}

	for i, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.pod, printers.PrintOptions{ShowLabels: test.showLabels})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		rows := table.Rows
		for i := range rows {
			rows[i].Object.Object = nil
		}
		if !reflect.DeepEqual(test.expect, rows) {
			t.Errorf("%d mismatch: %s", i, diff.ObjectReflectDiff(test.expect, rows))
		}
	}
}

func TestPrintService(t *testing.T) {
	single_ExternalIP := []string{"80.11.12.10"}
	mul_ExternalIP := []string{"80.11.12.10", "80.11.12.11"}
	tests := []struct {
		service api.Service
		expect  string
	}{
		{
			// Test name, cluster ip, port with protocol
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test1"},
				Spec: api.ServiceSpec{
					Type: api.ServiceTypeClusterIP,
					Ports: []api.ServicePort{
						{
							Protocol: "tcp",
							Port:     2233,
						},
					},
					ClusterIP: "10.9.8.7",
				},
			},
			"test1\tClusterIP\t10.9.8.7\t<none>\t2233/tcp\t<unknown>\n",
		},
		{
			// Test NodePort service
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test2"},
				Spec: api.ServiceSpec{
					Type: api.ServiceTypeNodePort,
					Ports: []api.ServicePort{
						{
							Protocol: "tcp",
							Port:     8888,
							NodePort: 9999,
						},
					},
					ClusterIP: "10.9.8.7",
				},
			},
			"test2\tNodePort\t10.9.8.7\t<none>\t8888:9999/tcp\t<unknown>\n",
		},
		{
			// Test LoadBalancer service
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test3"},
				Spec: api.ServiceSpec{
					Type: api.ServiceTypeLoadBalancer,
					Ports: []api.ServicePort{
						{
							Protocol: "tcp",
							Port:     8888,
						},
					},
					ClusterIP: "10.9.8.7",
				},
			},
			"test3\tLoadBalancer\t10.9.8.7\t<pending>\t8888/tcp\t<unknown>\n",
		},
		{
			// Test LoadBalancer service with single ExternalIP and no LoadBalancerStatus
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test4"},
				Spec: api.ServiceSpec{
					Type: api.ServiceTypeLoadBalancer,
					Ports: []api.ServicePort{
						{
							Protocol: "tcp",
							Port:     8888,
						},
					},
					ClusterIP:   "10.9.8.7",
					ExternalIPs: single_ExternalIP,
				},
			},
			"test4\tLoadBalancer\t10.9.8.7\t80.11.12.10\t8888/tcp\t<unknown>\n",
		},
		{
			// Test LoadBalancer service with single ExternalIP
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test5"},
				Spec: api.ServiceSpec{
					Type: api.ServiceTypeLoadBalancer,
					Ports: []api.ServicePort{
						{
							Protocol: "tcp",
							Port:     8888,
						},
					},
					ClusterIP:   "10.9.8.7",
					ExternalIPs: single_ExternalIP,
				},
				Status: api.ServiceStatus{
					LoadBalancer: api.LoadBalancerStatus{
						Ingress: []api.LoadBalancerIngress{
							{
								IP:       "3.4.5.6",
								Hostname: "test.cluster.com",
							},
						},
					},
				},
			},
			"test5\tLoadBalancer\t10.9.8.7\t3.4.5.6,80.11.12.10\t8888/tcp\t<unknown>\n",
		},
		{
			// Test LoadBalancer service with mul ExternalIPs
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test6"},
				Spec: api.ServiceSpec{
					Type: api.ServiceTypeLoadBalancer,
					Ports: []api.ServicePort{
						{
							Protocol: "tcp",
							Port:     8888,
						},
					},
					ClusterIP:   "10.9.8.7",
					ExternalIPs: mul_ExternalIP,
				},
				Status: api.ServiceStatus{
					LoadBalancer: api.LoadBalancerStatus{
						Ingress: []api.LoadBalancerIngress{
							{
								IP:       "2.3.4.5",
								Hostname: "test.cluster.local",
							},
							{
								IP:       "3.4.5.6",
								Hostname: "test.cluster.com",
							},
						},
					},
				},
			},
			"test6\tLoadBalancer\t10.9.8.7\t2.3.4.5,3.4.5.6,80.11.12.10,80.11.12.11\t8888/tcp\t<unknown>\n",
		},
		{
			// Test ExternalName service
			api.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test7"},
				Spec: api.ServiceSpec{
					Type:         api.ServiceTypeExternalName,
					ExternalName: "my.database.example.com",
				},
			},
			"test7\tExternalName\t<none>\tmy.database.example.com\t<none>\t<unknown>\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.service, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		// We ignore time
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, but got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintPodDisruptionBudget(t *testing.T) {
	minAvailable := intstr.FromInt(22)
	maxUnavailable := intstr.FromInt(11)
	tests := []struct {
		pdb    policy.PodDisruptionBudget
		expect string
	}{
		{
			policy.PodDisruptionBudget{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:         "ns1",
					Name:              "pdb1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: policy.PodDisruptionBudgetSpec{
					MinAvailable: &minAvailable,
				},
				Status: policy.PodDisruptionBudgetStatus{
					PodDisruptionsAllowed: 5,
				},
			},
			"pdb1\t22\tN/A\t5\t0s\n",
		},
		{
			policy.PodDisruptionBudget{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:         "ns2",
					Name:              "pdb2",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: policy.PodDisruptionBudgetSpec{
					MaxUnavailable: &maxUnavailable,
				},
				Status: policy.PodDisruptionBudgetStatus{
					PodDisruptionsAllowed: 5,
				},
			},
			"pdb2\tN/A\t11\t5\t0s\n",
		}}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.pdb, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintControllerRevision(t *testing.T) {
	tests := []struct {
		history apps.ControllerRevision
		expect  string
	}{
		{
			apps.ControllerRevision{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
					OwnerReferences: []metav1.OwnerReference{
						{
							Controller: boolP(true),
							APIVersion: "apps/v1",
							Kind:       "DaemonSet",
							Name:       "foo",
						},
					},
				},
				Revision: 1,
			},
			"test1\tdaemonset.apps/foo\t1\t0s\n",
		},
		{
			apps.ControllerRevision{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test2",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
					OwnerReferences: []metav1.OwnerReference{
						{
							Controller: boolP(false),
							Kind:       "ABC",
							Name:       "foo",
						},
					},
				},
				Revision: 2,
			},
			"test2\t<none>\t2\t0s\n",
		},
		{
			apps.ControllerRevision{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test3",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
					OwnerReferences:   []metav1.OwnerReference{},
				},
				Revision: 3,
			},
			"test3\t<none>\t3\t0s\n",
		},
		{
			apps.ControllerRevision{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test4",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
					OwnerReferences:   nil,
				},
				Revision: 4,
			},
			"test4\t<none>\t4\t0s\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.history, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, but got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func boolP(b bool) *bool {
	return &b
}

func TestPrintReplicaSet(t *testing.T) {
	tests := []struct {
		replicaSet extensions.ReplicaSet
		expect     string
		wideExpect string
	}{
		{
			extensions.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: extensions.ReplicaSetSpec{
					Replicas: 5,
					Template: api.PodTemplateSpec{
						Spec: api.PodSpec{
							Containers: []api.Container{
								{
									Name:  "fake-container1",
									Image: "fake-image1",
								},
								{
									Name:  "fake-container2",
									Image: "fake-image2",
								},
							},
						},
					},
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				},
				Status: extensions.ReplicaSetStatus{
					Replicas:      5,
					ReadyReplicas: 2,
				},
			},
			"test1\t5\t5\t2\t0s\n",
			"test1\t5\t5\t2\t0s\tfake-container1,fake-container2\tfake-image1,fake-image2\tfoo=bar\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.replicaSet, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()

		table, err = printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.replicaSet, printers.PrintOptions{Wide: true})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true, Wide: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.wideExpect {
			t.Fatalf("Expected: %s, got: %s", test.wideExpect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintPersistentVolumeClaim(t *testing.T) {
	myScn := "my-scn"
	tests := []struct {
		pvc    api.PersistentVolumeClaim
		expect string
	}{
		{
			// Test name, num of containers, restarts, container ready status
			api.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test1",
				},
				Spec: api.PersistentVolumeClaimSpec{
					VolumeName: "my-volume",
				},
				Status: api.PersistentVolumeClaimStatus{
					Phase:       api.ClaimBound,
					AccessModes: []api.PersistentVolumeAccessMode{api.ReadOnlyMany},
					Capacity: map[api.ResourceName]resource.Quantity{
						api.ResourceStorage: resource.MustParse("4Gi"),
					},
				},
			},
			"test1\tBound\tmy-volume\t4Gi\tROX\t\t<unknown>\n",
		},
		{
			// Test name, num of containers, restarts, container ready status
			api.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test2",
				},
				Spec: api.PersistentVolumeClaimSpec{},
				Status: api.PersistentVolumeClaimStatus{
					Phase:       api.ClaimLost,
					AccessModes: []api.PersistentVolumeAccessMode{api.ReadOnlyMany},
					Capacity: map[api.ResourceName]resource.Quantity{
						api.ResourceStorage: resource.MustParse("4Gi"),
					},
				},
			},
			"test2\tLost\t\t\t\t\t<unknown>\n",
		},
		{
			// Test name, num of containers, restarts, container ready status
			api.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test3",
				},
				Spec: api.PersistentVolumeClaimSpec{
					VolumeName: "my-volume",
				},
				Status: api.PersistentVolumeClaimStatus{
					Phase:       api.ClaimPending,
					AccessModes: []api.PersistentVolumeAccessMode{api.ReadWriteMany},
					Capacity: map[api.ResourceName]resource.Quantity{
						api.ResourceStorage: resource.MustParse("10Gi"),
					},
				},
			},
			"test3\tPending\tmy-volume\t10Gi\tRWX\t\t<unknown>\n",
		},
		{
			// Test name, num of containers, restarts, container ready status
			api.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test4",
				},
				Spec: api.PersistentVolumeClaimSpec{
					VolumeName:       "my-volume",
					StorageClassName: &myScn,
				},
				Status: api.PersistentVolumeClaimStatus{
					Phase:       api.ClaimPending,
					AccessModes: []api.PersistentVolumeAccessMode{api.ReadWriteOnce},
					Capacity: map[api.ResourceName]resource.Quantity{
						api.ResourceStorage: resource.MustParse("10Gi"),
					},
				},
			},
			"test4\tPending\tmy-volume\t10Gi\tRWO\tmy-scn\t<unknown>\n",
		},
	}
	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.pvc, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			fmt.Println(buf.String())
			fmt.Println(test.expect)
			t.Fatalf("Expected: %s, but got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintCronJob(t *testing.T) {
	suspend := false
	tests := []struct {
		cronjob batch.CronJob
		expect  string
	}{
		{
			batch.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "cronjob1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: batch.CronJobSpec{
					Schedule: "0/5 * * * ?",
					Suspend:  &suspend,
				},
				Status: batch.CronJobStatus{
					LastScheduleTime: &metav1.Time{Time: time.Now().Add(1.9e9)},
				},
			},
			"cronjob1\t0/5 * * * ?\tFalse\t0\t0s\t0s\n",
		},
		{
			batch.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "cronjob2",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-3e11)},
				},
				Spec: batch.CronJobSpec{
					Schedule: "0/5 * * * ?",
					Suspend:  &suspend,
				},
				Status: batch.CronJobStatus{
					LastScheduleTime: &metav1.Time{Time: time.Now().Add(-3e10)},
				},
			},
			"cronjob2\t0/5 * * * ?\tFalse\t0\t30s\t5m\n",
		},
		{
			batch.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "cronjob3",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-3e11)},
				},
				Spec: batch.CronJobSpec{
					Schedule: "0/5 * * * ?",
					Suspend:  &suspend,
				},
				Status: batch.CronJobStatus{},
			},
			"cronjob3\t0/5 * * * ?\tFalse\t0\t<none>\t5m\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.cronjob, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintStorageClass(t *testing.T) {
	tests := []struct {
		sc     storage.StorageClass
		expect string
	}{
		{
			storage.StorageClass{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sc1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Provisioner: "kubernetes.io/glusterfs",
			},
			"sc1\tkubernetes.io/glusterfs\t0s\n",
		},
		{
			storage.StorageClass{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sc2",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-3e11)},
				},
				Provisioner: "kubernetes.io/nfs",
			},
			"sc2\tkubernetes.io/nfs\t5m\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.sc, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		verifyTable(t, table)
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func TestPrintLease(t *testing.T) {
	holder1 := "holder1"
	holder2 := "holder2"
	tests := []struct {
		sc     coordination.Lease
		expect string
	}{
		{
			coordination.Lease{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "lease1",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(1.9e9)},
				},
				Spec: coordination.LeaseSpec{
					HolderIdentity: &holder1,
				},
			},
			"lease1\tholder1\t0s\n",
		},
		{
			coordination.Lease{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "lease2",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-3e11)},
				},
				Spec: coordination.LeaseSpec{
					HolderIdentity: &holder2,
				},
			},
			"lease2\tholder2\t5m\n",
		},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, test := range tests {
		table, err := printers.NewTablePrinter().With(AddHandlers).PrintTable(&test.sc, printers.PrintOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if err := printers.PrintTable(table, buf, printers.PrintOptions{NoHeaders: true}); err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.expect {
			t.Fatalf("Expected: %s, got: %s", test.expect, buf.String())
		}
		buf.Reset()
	}
}

func verifyTable(t *testing.T, table *metav1beta1.Table) {
	var panicErr interface{}
	func() {
		defer func() {
			panicErr = recover()
		}()
		table.DeepCopyObject() // cells are untyped, better check that types are JSON types and can be deep copied
	}()

	if panicErr != nil {
		t.Errorf("unexpected panic during deepcopy of table %#v: %v", table, panicErr)
	}
}
