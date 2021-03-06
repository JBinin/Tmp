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
Copyright 2017 The Kubernetes Authors.

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

package timer

import (
	"time"

	"bytes"
	"fmt"

	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/perftype"
	"sync"
)

var now = time.Now

// Represents a phase of a test. Phases can overlap.
type Phase struct {
	sequenceNumber int
	name           string
	startTime      time.Time
	endTime        time.Time
}

func (phase *Phase) ended() bool {
	return !phase.endTime.IsZero()
}

// End marks the phase as ended, unless it had already been ended before.
func (phase *Phase) End() {
	if !phase.ended() {
		phase.endTime = now()
	}
}

func (phase *Phase) label() string {
	return fmt.Sprintf("%03d-%s", phase.sequenceNumber, phase.name)
}

func (phase *Phase) duration() time.Duration {
	endTime := phase.endTime
	if !phase.ended() {
		endTime = now()
	}
	return endTime.Sub(phase.startTime)
}

func (phase *Phase) humanReadable() string {
	if phase.ended() {
		return fmt.Sprintf("Phase %s: %v\n", phase.label(), phase.duration())
	} else {
		return fmt.Sprintf("Phase %s: %v so far\n", phase.label(), phase.duration())
	}
}

// A TestPhaseTimer groups phases and provides a way to export their measurements as JSON or human-readable text.
// It is safe to use concurrently.
type TestPhaseTimer struct {
	lock   sync.Mutex
	phases []*Phase
}

// NewTestPhaseTimer creates a new TestPhaseTimer.
func NewTestPhaseTimer() *TestPhaseTimer {
	return &TestPhaseTimer{}
}

// StartPhase starts a new phase.
// sequenceNumber is an integer prepended to phaseName in the output, such that lexicographic sorting
// of phases in perfdash reconstructs the order of execution. Unfortunately it needs to be
// provided manually, since a simple incrementing counter would have the effect that inserting
// a new phase would renumber subsequent phases, breaking the continuity of historical records.
func (timer *TestPhaseTimer) StartPhase(sequenceNumber int, phaseName string) *Phase {
	timer.lock.Lock()
	defer timer.lock.Unlock()
	newPhase := &Phase{sequenceNumber: sequenceNumber, name: phaseName, startTime: now()}
	timer.phases = append(timer.phases, newPhase)
	return newPhase
}

func (timer *TestPhaseTimer) SummaryKind() string {
	return "TestPhaseTimer"
}

func (timer *TestPhaseTimer) PrintHumanReadable() string {
	buf := bytes.Buffer{}
	timer.lock.Lock()
	defer timer.lock.Unlock()
	for _, phase := range timer.phases {
		buf.WriteString(phase.humanReadable())
	}
	return buf.String()
}

func (timer *TestPhaseTimer) PrintJSON() string {
	data := perftype.PerfData{
		Version: "v1",
		DataItems: []perftype.DataItem{{
			Unit:   "s",
			Labels: map[string]string{"test": "phases"},
			Data:   make(map[string]float64)}}}
	timer.lock.Lock()
	defer timer.lock.Unlock()
	for _, phase := range timer.phases {
		data.DataItems[0].Data[phase.label()] = phase.duration().Seconds()
		if !phase.ended() {
			data.DataItems[0].Labels["ended"] = "false"
		}
	}
	return framework.PrettyPrintJSON(data)
}
