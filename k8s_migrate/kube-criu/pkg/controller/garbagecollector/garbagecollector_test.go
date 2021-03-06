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
Copyright 2016 The Kubernetes Authors.

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

package garbagecollector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "k8s.io/kubernetes/pkg/apis/core/install"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

type testRESTMapper struct {
	meta.RESTMapper
}

func (_ *testRESTMapper) Reset() {}

func TestGarbageCollectorConstruction(t *testing.T) {
	config := &restclient.Config{}
	tweakableRM := meta.NewDefaultRESTMapper(nil)
	rm := &testRESTMapper{meta.MultiRESTMapper{tweakableRM, testrestmapper.TestOnlyStaticRESTMapper(legacyscheme.Scheme)}}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	podResource := map[schema.GroupVersionResource]struct{}{
		{Version: "v1", Resource: "pods"}: {},
	}
	twoResources := map[schema.GroupVersionResource]struct{}{
		{Version: "v1", Resource: "pods"}:                     {},
		{Group: "tpr.io", Version: "v1", Resource: "unknown"}: {},
	}
	client := fake.NewSimpleClientset()
	sharedInformers := informers.NewSharedInformerFactory(client, 0)

	// No monitor will be constructed for the non-core resource, but the GC
	// construction will not fail.
	alwaysStarted := make(chan struct{})
	close(alwaysStarted)
	gc, err := NewGarbageCollector(dynamicClient, rm, twoResources, map[schema.GroupResource]struct{}{}, sharedInformers, alwaysStarted)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(gc.dependencyGraphBuilder.monitors))

	// Make sure resource monitor syncing creates and stops resource monitors.
	tweakableRM.Add(schema.GroupVersionKind{Group: "tpr.io", Version: "v1", Kind: "unknown"}, nil)
	err = gc.resyncMonitors(twoResources)
	if err != nil {
		t.Errorf("Failed adding a monitor: %v", err)
	}
	assert.Equal(t, 2, len(gc.dependencyGraphBuilder.monitors))

	err = gc.resyncMonitors(podResource)
	if err != nil {
		t.Errorf("Failed removing a monitor: %v", err)
	}
	assert.Equal(t, 1, len(gc.dependencyGraphBuilder.monitors))

	// Make sure the syncing mechanism also works after Run() has been called
	stopCh := make(chan struct{})
	defer close(stopCh)
	go gc.Run(1, stopCh)

	err = gc.resyncMonitors(twoResources)
	if err != nil {
		t.Errorf("Failed adding a monitor: %v", err)
	}
	assert.Equal(t, 2, len(gc.dependencyGraphBuilder.monitors))

	err = gc.resyncMonitors(podResource)
	if err != nil {
		t.Errorf("Failed removing a monitor: %v", err)
	}
	assert.Equal(t, 1, len(gc.dependencyGraphBuilder.monitors))
}

// fakeAction records information about requests to aid in testing.
type fakeAction struct {
	method string
	path   string
	query  string
}

// String returns method=path to aid in testing
func (f *fakeAction) String() string {
	return strings.Join([]string{f.method, f.path}, "=")
}

type FakeResponse struct {
	statusCode int
	content    []byte
}

// fakeActionHandler holds a list of fakeActions received
type fakeActionHandler struct {
	// statusCode and content returned by this handler for different method + path.
	response map[string]FakeResponse

	lock    sync.Mutex
	actions []fakeAction
}

// ServeHTTP logs the action that occurred and always returns the associated status code
func (f *fakeActionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	func() {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.actions = append(f.actions, fakeAction{method: request.Method, path: request.URL.Path, query: request.URL.RawQuery})
		fakeResponse, ok := f.response[request.Method+request.URL.Path]
		if !ok {
			fakeResponse.statusCode = 200
			fakeResponse.content = []byte("{\"kind\": \"List\"}")
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(fakeResponse.statusCode)
		response.Write(fakeResponse.content)
	}()

	// This is to allow the fakeActionHandler to simulate a watch being opened
	if strings.Contains(request.URL.RawQuery, "watch=true") {
		hijacker, ok := response.(http.Hijacker)
		if !ok {
			return
		}
		connection, _, err := hijacker.Hijack()
		if err != nil {
			return
		}
		defer connection.Close()
		time.Sleep(30 * time.Second)
	}
}

// testServerAndClientConfig returns a server that listens and a config that can reference it
func testServerAndClientConfig(handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *restclient.Config) {
	srv := httptest.NewServer(http.HandlerFunc(handler))
	config := &restclient.Config{
		Host: srv.URL,
	}
	return srv, config
}

type garbageCollector struct {
	*GarbageCollector
	stop chan struct{}
}

func setupGC(t *testing.T, config *restclient.Config) garbageCollector {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	podResource := map[schema.GroupVersionResource]struct{}{{Version: "v1", Resource: "pods"}: {}}
	client := fake.NewSimpleClientset()
	sharedInformers := informers.NewSharedInformerFactory(client, 0)
	alwaysStarted := make(chan struct{})
	close(alwaysStarted)
	gc, err := NewGarbageCollector(dynamicClient, &testRESTMapper{testrestmapper.TestOnlyStaticRESTMapper(legacyscheme.Scheme)}, podResource, ignoredResources, sharedInformers, alwaysStarted)
	if err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	go sharedInformers.Start(stop)
	return garbageCollector{gc, stop}
}

func getPod(podName string, ownerReferences []metav1.OwnerReference) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            podName,
			Namespace:       "ns1",
			OwnerReferences: ownerReferences,
		},
	}
}

func serilizeOrDie(t *testing.T, object interface{}) []byte {
	data, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// test the attemptToDeleteItem function making the expected actions.
func TestAttemptToDeleteItem(t *testing.T) {
	pod := getPod("ToBeDeletedPod", []metav1.OwnerReference{
		{
			Kind:       "ReplicationController",
			Name:       "owner1",
			UID:        "123",
			APIVersion: "v1",
		},
	})
	testHandler := &fakeActionHandler{
		response: map[string]FakeResponse{
			"GET" + "/api/v1/namespaces/ns1/replicationcontrollers/owner1": {
				404,
				[]byte{},
			},
			"GET" + "/api/v1/namespaces/ns1/pods/ToBeDeletedPod": {
				200,
				serilizeOrDie(t, pod),
			},
		},
	}
	srv, clientConfig := testServerAndClientConfig(testHandler.ServeHTTP)
	defer srv.Close()

	gc := setupGC(t, clientConfig)
	defer close(gc.stop)

	item := &node{
		identity: objectReference{
			OwnerReference: metav1.OwnerReference{
				Kind:       pod.Kind,
				APIVersion: pod.APIVersion,
				Name:       pod.Name,
				UID:        pod.UID,
			},
			Namespace: pod.Namespace,
		},
		// owners are intentionally left empty. The attemptToDeleteItem routine should get the latest item from the server.
		owners: nil,
	}
	err := gc.attemptToDeleteItem(item)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	expectedActionSet := sets.NewString()
	expectedActionSet.Insert("GET=/api/v1/namespaces/ns1/replicationcontrollers/owner1")
	expectedActionSet.Insert("DELETE=/api/v1/namespaces/ns1/pods/ToBeDeletedPod")
	expectedActionSet.Insert("GET=/api/v1/namespaces/ns1/pods/ToBeDeletedPod")

	actualActionSet := sets.NewString()
	for _, action := range testHandler.actions {
		actualActionSet.Insert(action.String())
	}
	if !expectedActionSet.Equal(actualActionSet) {
		t.Errorf("expected actions:\n%v\n but got:\n%v\nDifference:\n%v", expectedActionSet,
			actualActionSet, expectedActionSet.Difference(actualActionSet))
	}
}

// verifyGraphInvariants verifies that all of a node's owners list the node as a
// dependent and vice versa. uidToNode has all the nodes in the graph.
func verifyGraphInvariants(scenario string, uidToNode map[types.UID]*node, t *testing.T) {
	for myUID, node := range uidToNode {
		for dependentNode := range node.dependents {
			found := false
			for _, owner := range dependentNode.owners {
				if owner.UID == myUID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("scenario: %s: node %s has node %s as a dependent, but it's not present in the latter node's owners list", scenario, node.identity, dependentNode.identity)
			}
		}

		for _, owner := range node.owners {
			ownerNode, ok := uidToNode[owner.UID]
			if !ok {
				// It's possible that the owner node doesn't exist
				continue
			}
			if _, ok := ownerNode.dependents[node]; !ok {
				t.Errorf("node %s has node %s as an owner, but it's not present in the latter node's dependents list", node.identity, ownerNode.identity)
			}
		}
	}
}

func createEvent(eventType eventType, selfUID string, owners []string) event {
	var ownerReferences []metav1.OwnerReference
	for i := 0; i < len(owners); i++ {
		ownerReferences = append(ownerReferences, metav1.OwnerReference{UID: types.UID(owners[i])})
	}
	return event{
		eventType: eventType,
		obj: &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:             types.UID(selfUID),
				OwnerReferences: ownerReferences,
			},
		},
	}
}

func TestProcessEvent(t *testing.T) {
	var testScenarios = []struct {
		name string
		// a series of events that will be supplied to the
		// GraphBuilder.graphChanges.
		events []event
	}{
		{
			name: "test1",
			events: []event{
				createEvent(addEvent, "1", []string{}),
				createEvent(addEvent, "2", []string{"1"}),
				createEvent(addEvent, "3", []string{"1", "2"}),
			},
		},
		{
			name: "test2",
			events: []event{
				createEvent(addEvent, "1", []string{}),
				createEvent(addEvent, "2", []string{"1"}),
				createEvent(addEvent, "3", []string{"1", "2"}),
				createEvent(addEvent, "4", []string{"2"}),
				createEvent(deleteEvent, "2", []string{"doesn't matter"}),
			},
		},
		{
			name: "test3",
			events: []event{
				createEvent(addEvent, "1", []string{}),
				createEvent(addEvent, "2", []string{"1"}),
				createEvent(addEvent, "3", []string{"1", "2"}),
				createEvent(addEvent, "4", []string{"3"}),
				createEvent(updateEvent, "2", []string{"4"}),
			},
		},
		{
			name: "reverse test2",
			events: []event{
				createEvent(addEvent, "4", []string{"2"}),
				createEvent(addEvent, "3", []string{"1", "2"}),
				createEvent(addEvent, "2", []string{"1"}),
				createEvent(addEvent, "1", []string{}),
				createEvent(deleteEvent, "2", []string{"doesn't matter"}),
			},
		},
	}

	alwaysStarted := make(chan struct{})
	close(alwaysStarted)
	for _, scenario := range testScenarios {
		dependencyGraphBuilder := &GraphBuilder{
			informersStarted: alwaysStarted,
			graphChanges:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
			uidToNode: &concurrentUIDToNode{
				uidToNodeLock: sync.RWMutex{},
				uidToNode:     make(map[types.UID]*node),
			},
			attemptToDelete:  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
			absentOwnerCache: NewUIDCache(2),
		}
		for i := 0; i < len(scenario.events); i++ {
			dependencyGraphBuilder.graphChanges.Add(&scenario.events[i])
			dependencyGraphBuilder.processGraphChanges()
			verifyGraphInvariants(scenario.name, dependencyGraphBuilder.uidToNode.uidToNode, t)
		}
	}
}

// TestDependentsRace relies on golang's data race detector to check if there is
// data race among in the dependents field.
func TestDependentsRace(t *testing.T) {
	gc := setupGC(t, &restclient.Config{})
	defer close(gc.stop)

	const updates = 100
	owner := &node{dependents: make(map[*node]struct{})}
	ownerUID := types.UID("owner")
	gc.dependencyGraphBuilder.uidToNode.Write(owner)
	go func() {
		for i := 0; i < updates; i++ {
			dependent := &node{}
			gc.dependencyGraphBuilder.addDependentToOwners(dependent, []metav1.OwnerReference{{UID: ownerUID}})
			gc.dependencyGraphBuilder.removeDependentFromOwners(dependent, []metav1.OwnerReference{{UID: ownerUID}})
		}
	}()
	go func() {
		gc.attemptToOrphan.Add(owner)
		for i := 0; i < updates; i++ {
			gc.attemptToOrphanWorker()
		}
	}()
}

// test the list and watch functions correctly converts the ListOptions
func TestGCListWatcher(t *testing.T) {
	testHandler := &fakeActionHandler{}
	srv, clientConfig := testServerAndClientConfig(testHandler.ServeHTTP)
	defer srv.Close()
	podResource := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	lw := listWatcher(dynamicClient, podResource)
	lw.DisableChunking = true
	if _, err := lw.Watch(metav1.ListOptions{ResourceVersion: "1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := lw.List(metav1.ListOptions{ResourceVersion: "1"}); err != nil {
		t.Fatal(err)
	}
	if e, a := 2, len(testHandler.actions); e != a {
		t.Errorf("expect %d requests, got %d", e, a)
	}
	if e, a := "resourceVersion=1&watch=true", testHandler.actions[0].query; e != a {
		t.Errorf("expect %s, got %s", e, a)
	}
	if e, a := "resourceVersion=1", testHandler.actions[1].query; e != a {
		t.Errorf("expect %s, got %s", e, a)
	}
}

func podToGCNode(pod *v1.Pod) *node {
	return &node{
		identity: objectReference{
			OwnerReference: metav1.OwnerReference{
				Kind:       pod.Kind,
				APIVersion: pod.APIVersion,
				Name:       pod.Name,
				UID:        pod.UID,
			},
			Namespace: pod.Namespace,
		},
		// owners are intentionally left empty. The attemptToDeleteItem routine should get the latest item from the server.
		owners: nil,
	}
}

func TestAbsentUIDCache(t *testing.T) {
	rc1Pod1 := getPod("rc1Pod1", []metav1.OwnerReference{
		{
			Kind:       "ReplicationController",
			Name:       "rc1",
			UID:        "1",
			APIVersion: "v1",
		},
	})
	rc1Pod2 := getPod("rc1Pod2", []metav1.OwnerReference{
		{
			Kind:       "ReplicationController",
			Name:       "rc1",
			UID:        "1",
			APIVersion: "v1",
		},
	})
	rc2Pod1 := getPod("rc2Pod1", []metav1.OwnerReference{
		{
			Kind:       "ReplicationController",
			Name:       "rc2",
			UID:        "2",
			APIVersion: "v1",
		},
	})
	rc3Pod1 := getPod("rc3Pod1", []metav1.OwnerReference{
		{
			Kind:       "ReplicationController",
			Name:       "rc3",
			UID:        "3",
			APIVersion: "v1",
		},
	})
	testHandler := &fakeActionHandler{
		response: map[string]FakeResponse{
			"GET" + "/api/v1/namespaces/ns1/pods/rc1Pod1": {
				200,
				serilizeOrDie(t, rc1Pod1),
			},
			"GET" + "/api/v1/namespaces/ns1/pods/rc1Pod2": {
				200,
				serilizeOrDie(t, rc1Pod2),
			},
			"GET" + "/api/v1/namespaces/ns1/pods/rc2Pod1": {
				200,
				serilizeOrDie(t, rc2Pod1),
			},
			"GET" + "/api/v1/namespaces/ns1/pods/rc3Pod1": {
				200,
				serilizeOrDie(t, rc3Pod1),
			},
			"GET" + "/api/v1/namespaces/ns1/replicationcontrollers/rc1": {
				404,
				[]byte{},
			},
			"GET" + "/api/v1/namespaces/ns1/replicationcontrollers/rc2": {
				404,
				[]byte{},
			},
			"GET" + "/api/v1/namespaces/ns1/replicationcontrollers/rc3": {
				404,
				[]byte{},
			},
		},
	}
	srv, clientConfig := testServerAndClientConfig(testHandler.ServeHTTP)
	defer srv.Close()
	gc := setupGC(t, clientConfig)
	defer close(gc.stop)
	gc.absentOwnerCache = NewUIDCache(2)
	gc.attemptToDeleteItem(podToGCNode(rc1Pod1))
	gc.attemptToDeleteItem(podToGCNode(rc2Pod1))
	// rc1 should already be in the cache, no request should be sent. rc1 should be promoted in the UIDCache
	gc.attemptToDeleteItem(podToGCNode(rc1Pod2))
	// after this call, rc2 should be evicted from the UIDCache
	gc.attemptToDeleteItem(podToGCNode(rc3Pod1))
	// check cache
	if !gc.absentOwnerCache.Has(types.UID("1")) {
		t.Errorf("expected rc1 to be in the cache")
	}
	if gc.absentOwnerCache.Has(types.UID("2")) {
		t.Errorf("expected rc2 to not exist in the cache")
	}
	if !gc.absentOwnerCache.Has(types.UID("3")) {
		t.Errorf("expected rc3 to be in the cache")
	}
	// check the request sent to the server
	count := 0
	for _, action := range testHandler.actions {
		if action.String() == "GET=/api/v1/namespaces/ns1/replicationcontrollers/rc1" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected only 1 GET rc1 request, got %d", count)
	}
}

func TestDeleteOwnerRefPatch(t *testing.T) {
	original := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			UID: "100",
			OwnerReferences: []metav1.OwnerReference{
				{UID: "1"},
				{UID: "2"},
				{UID: "3"},
			},
		},
	}
	originalData := serilizeOrDie(t, original)
	expected := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			UID: "100",
			OwnerReferences: []metav1.OwnerReference{
				{UID: "1"},
			},
		},
	}
	patch := deleteOwnerRefStrategicMergePatch("100", "2", "3")
	patched, err := strategicpatch.StrategicMergePatch(originalData, patch, v1.Pod{})
	if err != nil {
		t.Fatal(err)
	}
	var got v1.Pod
	if err := json.Unmarshal(patched, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected: %#v,\ngot: %#v", expected, got)
	}
}

func TestUnblockOwnerReference(t *testing.T) {
	trueVar := true
	falseVar := false
	original := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			UID: "100",
			OwnerReferences: []metav1.OwnerReference{
				{UID: "1", BlockOwnerDeletion: &trueVar},
				{UID: "2", BlockOwnerDeletion: &falseVar},
				{UID: "3"},
			},
		},
	}
	originalData := serilizeOrDie(t, original)
	expected := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			UID: "100",
			OwnerReferences: []metav1.OwnerReference{
				{UID: "1", BlockOwnerDeletion: &falseVar},
				{UID: "2", BlockOwnerDeletion: &falseVar},
				{UID: "3"},
			},
		},
	}
	accessor, err := meta.Accessor(&original)
	if err != nil {
		t.Fatal(err)
	}
	n := node{
		owners: accessor.GetOwnerReferences(),
	}
	patch, err := n.unblockOwnerReferencesStrategicMergePatch()
	if err != nil {
		t.Fatal(err)
	}
	patched, err := strategicpatch.StrategicMergePatch(originalData, patch, v1.Pod{})
	if err != nil {
		t.Fatal(err)
	}
	var got v1.Pod
	if err := json.Unmarshal(patched, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected: %#v,\ngot: %#v", expected, got)
		t.Errorf("expected: %#v,\ngot: %#v", expected.OwnerReferences, got.OwnerReferences)
		for _, ref := range got.OwnerReferences {
			t.Errorf("ref.UID=%s, ref.BlockOwnerDeletion=%v", ref.UID, *ref.BlockOwnerDeletion)
		}
	}
}

func TestOrphanDependentsFailure(t *testing.T) {
	testHandler := &fakeActionHandler{
		response: map[string]FakeResponse{
			"PATCH" + "/api/v1/namespaces/ns1/pods/pod": {
				409,
				[]byte{},
			},
		},
	}
	srv, clientConfig := testServerAndClientConfig(testHandler.ServeHTTP)
	defer srv.Close()

	gc := setupGC(t, clientConfig)
	defer close(gc.stop)

	dependents := []*node{
		{
			identity: objectReference{
				OwnerReference: metav1.OwnerReference{
					Kind:       "Pod",
					APIVersion: "v1",
					Name:       "pod",
				},
				Namespace: "ns1",
			},
		},
	}
	err := gc.orphanDependents(objectReference{}, dependents)
	expected := `the server reported a conflict`
	if err == nil || !strings.Contains(err.Error(), expected) {
		if err != nil {
			t.Errorf("expected error contains text %q, got %q", expected, err.Error())
		} else {
			t.Errorf("expected error contains text %q, got nil", expected)
		}
	}
}

// TestGetDeletableResources ensures GetDeletableResources always returns
// something usable regardless of discovery output.
func TestGetDeletableResources(t *testing.T) {
	tests := map[string]struct {
		serverResources    []*metav1.APIResourceList
		err                error
		deletableResources map[schema.GroupVersionResource]struct{}
	}{
		"no error": {
			serverResources: []*metav1.APIResourceList{
				{
					// Valid GroupVersion
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{Name: "pods", Namespaced: true, Kind: "Pod", Verbs: metav1.Verbs{"delete", "list", "watch"}},
						{Name: "services", Namespaced: true, Kind: "Service"},
					},
				},
				{
					// Invalid GroupVersion, should be ignored
					GroupVersion: "foo//whatever",
					APIResources: []metav1.APIResource{
						{Name: "bars", Namespaced: true, Kind: "Bar", Verbs: metav1.Verbs{"delete", "list", "watch"}},
					},
				},
				{
					// Valid GroupVersion, missing required verbs, should be ignored
					GroupVersion: "acme/v1",
					APIResources: []metav1.APIResource{
						{Name: "widgets", Namespaced: true, Kind: "Widget", Verbs: metav1.Verbs{"delete"}},
					},
				},
			},
			err: nil,
			deletableResources: map[schema.GroupVersionResource]struct{}{
				{Group: "apps", Version: "v1", Resource: "pods"}: {},
			},
		},
		"nonspecific failure, includes usable results": {
			serverResources: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{Name: "pods", Namespaced: true, Kind: "Pod", Verbs: metav1.Verbs{"delete", "list", "watch"}},
						{Name: "services", Namespaced: true, Kind: "Service"},
					},
				},
			},
			err: fmt.Errorf("internal error"),
			deletableResources: map[schema.GroupVersionResource]struct{}{
				{Group: "apps", Version: "v1", Resource: "pods"}: {},
			},
		},
		"partial discovery failure, includes usable results": {
			serverResources: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{Name: "pods", Namespaced: true, Kind: "Pod", Verbs: metav1.Verbs{"delete", "list", "watch"}},
						{Name: "services", Namespaced: true, Kind: "Service"},
					},
				},
			},
			err: &discovery.ErrGroupDiscoveryFailed{
				Groups: map[schema.GroupVersion]error{
					{Group: "foo", Version: "v1"}: fmt.Errorf("discovery failure"),
				},
			},
			deletableResources: map[schema.GroupVersionResource]struct{}{
				{Group: "apps", Version: "v1", Resource: "pods"}: {},
			},
		},
		"discovery failure, no results": {
			serverResources:    nil,
			err:                fmt.Errorf("internal error"),
			deletableResources: map[schema.GroupVersionResource]struct{}{},
		},
	}

	for name, test := range tests {
		t.Logf("testing %q", name)
		client := &fakeServerResources{
			PreferredResources: test.serverResources,
			Error:              test.err,
		}
		actual := GetDeletableResources(client)
		if !reflect.DeepEqual(test.deletableResources, actual) {
			t.Errorf("expected resources:\n%v\ngot:\n%v", test.deletableResources, actual)
		}
	}
}

// TestGarbageCollectorSync ensures that a discovery client error
// will not cause the garbage collector to block infinitely.
func TestGarbageCollectorSync(t *testing.T) {
	serverResources := []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Name: "pods", Namespaced: true, Kind: "Pod", Verbs: metav1.Verbs{"delete", "list", "watch"}},
			},
		},
	}
	unsyncableServerResources := []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Name: "pods", Namespaced: true, Kind: "Pod", Verbs: metav1.Verbs{"delete", "list", "watch"}},
				{Name: "secrets", Namespaced: true, Kind: "Secret", Verbs: metav1.Verbs{"delete", "list", "watch"}},
			},
		},
	}
	fakeDiscoveryClient := &fakeServerResources{
		PreferredResources: serverResources,
		Error:              nil,
		Lock:               sync.Mutex{},
		InterfaceUsedCount: 0,
	}

	testHandler := &fakeActionHandler{
		response: map[string]FakeResponse{
			"GET" + "/api/v1/pods": {
				200,
				[]byte("{}"),
			},
			"GET" + "/api/v1/secrets": {
				404,
				[]byte("{}"),
			},
		},
	}
	srv, clientConfig := testServerAndClientConfig(testHandler.ServeHTTP)
	defer srv.Close()
	clientConfig.ContentConfig.NegotiatedSerializer = nil
	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	rm := &testRESTMapper{testrestmapper.TestOnlyStaticRESTMapper(legacyscheme.Scheme)}
	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	podResource := map[schema.GroupVersionResource]struct{}{
		{Group: "", Version: "v1", Resource: "pods"}: {},
	}
	sharedInformers := informers.NewSharedInformerFactory(client, 0)
	alwaysStarted := make(chan struct{})
	close(alwaysStarted)
	gc, err := NewGarbageCollector(dynamicClient, rm, podResource, map[schema.GroupResource]struct{}{}, sharedInformers, alwaysStarted)
	if err != nil {
		t.Fatal(err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	go gc.Run(1, stopCh)
	go gc.Sync(fakeDiscoveryClient, 10*time.Millisecond, stopCh)

	// Wait until the sync discovers the initial resources
	fmt.Printf("Test output")
	time.Sleep(1 * time.Second)

	err = expectSyncNotBlocked(fakeDiscoveryClient, &gc.workerLock)
	if err != nil {
		t.Fatalf("Expected garbagecollector.Sync to be running but it is blocked: %v", err)
	}

	// Simulate the discovery client returning an error
	fakeDiscoveryClient.setPreferredResources(nil)
	fakeDiscoveryClient.setError(fmt.Errorf("Error calling discoveryClient.ServerPreferredResources()"))

	// Wait until sync discovers the change
	time.Sleep(1 * time.Second)

	// Remove the error from being returned and see if the garbage collector sync is still working
	fakeDiscoveryClient.setPreferredResources(serverResources)
	fakeDiscoveryClient.setError(nil)

	err = expectSyncNotBlocked(fakeDiscoveryClient, &gc.workerLock)
	if err != nil {
		t.Fatalf("Expected garbagecollector.Sync to still be running but it is blocked: %v", err)
	}

	// Simulate the discovery client returning a resource the restmapper can resolve, but will not sync caches
	fakeDiscoveryClient.setPreferredResources(unsyncableServerResources)
	fakeDiscoveryClient.setError(nil)

	// Wait until sync discovers the change
	time.Sleep(1 * time.Second)

	// Put the resources back to normal and ensure garbage collector sync recovers
	fakeDiscoveryClient.setPreferredResources(serverResources)
	fakeDiscoveryClient.setError(nil)

	err = expectSyncNotBlocked(fakeDiscoveryClient, &gc.workerLock)
	if err != nil {
		t.Fatalf("Expected garbagecollector.Sync to still be running but it is blocked: %v", err)
	}
}

func expectSyncNotBlocked(fakeDiscoveryClient *fakeServerResources, workerLock *sync.RWMutex) error {
	before := fakeDiscoveryClient.getInterfaceUsedCount()
	t := 1 * time.Second
	time.Sleep(t)
	after := fakeDiscoveryClient.getInterfaceUsedCount()
	if before == after {
		return fmt.Errorf("discoveryClient.ServerPreferredResources() called %d times over %v", after-before, t)
	}

	workerLockAcquired := make(chan struct{})
	go func() {
		workerLock.Lock()
		workerLock.Unlock()
		close(workerLockAcquired)
	}()
	select {
	case <-workerLockAcquired:
		return nil
	case <-time.After(t):
		return fmt.Errorf("workerLock blocked for at least %v", t)
	}
}

type fakeServerResources struct {
	PreferredResources []*metav1.APIResourceList
	Error              error
	Lock               sync.Mutex
	InterfaceUsedCount int
}

func (_ *fakeServerResources) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	return nil, nil
}

func (_ *fakeServerResources) ServerResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (f *fakeServerResources) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	f.InterfaceUsedCount++
	return f.PreferredResources, f.Error
}

func (f *fakeServerResources) setPreferredResources(resources []*metav1.APIResourceList) {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	f.PreferredResources = resources
}

func (f *fakeServerResources) setError(err error) {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	f.Error = err
}

func (f *fakeServerResources) getInterfaceUsedCount() int {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	return f.InterfaceUsedCount
}

func (_ *fakeServerResources) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}
