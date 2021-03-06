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

package stats

import (
	"sync"
	"sync/atomic"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/golang/glog"
)

// Map to PodVolumeStats pointers since the addresses for map values are not constant and can cause pain
// if we need ever to get a pointer to one of the values (e.g. you can't)
type Cache map[types.UID]*volumeStatCalculator

// fsResourceAnalyzerInterface is for embedding fs functions into ResourceAnalyzer
type fsResourceAnalyzerInterface interface {
	GetPodVolumeStats(uid types.UID) (PodVolumeStats, bool)
}

// fsResourceAnalyzer provides stats about fs resource usage
type fsResourceAnalyzer struct {
	statsProvider     StatsProvider
	calcPeriod        time.Duration
	cachedVolumeStats atomic.Value
	startOnce         sync.Once
}

var _ fsResourceAnalyzerInterface = &fsResourceAnalyzer{}

// newFsResourceAnalyzer returns a new fsResourceAnalyzer implementation
func newFsResourceAnalyzer(statsProvider StatsProvider, calcVolumePeriod time.Duration) *fsResourceAnalyzer {
	r := &fsResourceAnalyzer{
		statsProvider: statsProvider,
		calcPeriod:    calcVolumePeriod,
	}
	r.cachedVolumeStats.Store(make(Cache))
	return r
}

// Start eager background caching of volume stats.
func (s *fsResourceAnalyzer) Start() {
	s.startOnce.Do(func() {
		if s.calcPeriod <= 0 {
			glog.Info("Volume stats collection disabled.")
			return
		}
		glog.Info("Starting FS ResourceAnalyzer")
		go wait.Forever(func() { s.updateCachedPodVolumeStats() }, s.calcPeriod)
	})
}

// updateCachedPodVolumeStats calculates and caches the PodVolumeStats for every Pod known to the kubelet.
func (s *fsResourceAnalyzer) updateCachedPodVolumeStats() {
	oldCache := s.cachedVolumeStats.Load().(Cache)
	newCache := make(Cache)

	// Copy existing entries to new map, creating/starting new entries for pods missing from the cache
	for _, pod := range s.statsProvider.GetPods() {
		if value, found := oldCache[pod.GetUID()]; !found {
			newCache[pod.GetUID()] = newVolumeStatCalculator(s.statsProvider, s.calcPeriod, pod).StartOnce()
		} else {
			newCache[pod.GetUID()] = value
		}
	}

	// Stop entries for pods that have been deleted
	for uid, entry := range oldCache {
		if _, found := newCache[uid]; !found {
			entry.StopOnce()
		}
	}

	// Update the cache reference
	s.cachedVolumeStats.Store(newCache)
}

// GetPodVolumeStats returns the PodVolumeStats for a given pod.  Results are looked up from a cache that
// is eagerly populated in the background, and never calculated on the fly.
func (s *fsResourceAnalyzer) GetPodVolumeStats(uid types.UID) (PodVolumeStats, bool) {
	cache := s.cachedVolumeStats.Load().(Cache)
	if statCalc, found := cache[uid]; !found {
		// TODO: Differentiate between stats being empty
		// See issue #20679
		return PodVolumeStats{}, false
	} else {
		return statCalc.GetLatest()
	}
}
