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
// +build linux

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

// Reads the pod configuration from file or a directory of files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/exp/inotify"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/flowcontrol"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	retryPeriod    = 1 * time.Second
	maxRetryPeriod = 20 * time.Second
)

type retryableError struct {
	message string
}

func (e *retryableError) Error() string {
	return e.message
}

func (s *sourceFile) startWatch() {
	backOff := flowcontrol.NewBackOff(retryPeriod, maxRetryPeriod)
	backOffId := "watch"

	go wait.Forever(func() {
		if backOff.IsInBackOffSinceUpdate(backOffId, time.Now()) {
			return
		}

		if err := s.doWatch(); err != nil {
			glog.Errorf("Unable to read config path %q: %v", s.path, err)
			if _, retryable := err.(*retryableError); !retryable {
				backOff.Next(backOffId, time.Now())
			}
		}
	}, retryPeriod)
}

func (s *sourceFile) doWatch() error {
	_, err := os.Stat(s.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// Emit an update with an empty PodList to allow FileSource to be marked as seen
		s.updates <- kubetypes.PodUpdate{Pods: []*v1.Pod{}, Op: kubetypes.SET, Source: kubetypes.FileSource}
		return &retryableError{"path does not exist, ignoring"}
	}

	w, err := inotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("unable to create inotify: %v", err)
	}
	defer w.Close()

	err = w.AddWatch(s.path, inotify.IN_DELETE_SELF|inotify.IN_CREATE|inotify.IN_MOVED_TO|inotify.IN_MODIFY|inotify.IN_MOVED_FROM|inotify.IN_DELETE|inotify.IN_ATTRIB)
	if err != nil {
		return fmt.Errorf("unable to create inotify for path %q: %v", s.path, err)
	}

	for {
		select {
		case event := <-w.Event:
			if err = s.produceWatchEvent(event); err != nil {
				return fmt.Errorf("error while processing inotify event (%+v): %v", event, err)
			}
		case err = <-w.Error:
			return fmt.Errorf("error while watching %q: %v", s.path, err)
		}
	}
}

func (s *sourceFile) produceWatchEvent(e *inotify.Event) error {
	// Ignore file start with dots
	if strings.HasPrefix(filepath.Base(e.Name), ".") {
		glog.V(4).Infof("Ignored pod manifest: %s, because it starts with dots", e.Name)
		return nil
	}
	var eventType podEventType
	switch {
	case (e.Mask & inotify.IN_ISDIR) > 0:
		glog.Errorf("Not recursing into manifest path %q", s.path)
		return nil
	case (e.Mask & inotify.IN_CREATE) > 0:
		eventType = podAdd
	case (e.Mask & inotify.IN_MOVED_TO) > 0:
		eventType = podAdd
	case (e.Mask & inotify.IN_MODIFY) > 0:
		eventType = podModify
	case (e.Mask & inotify.IN_ATTRIB) > 0:
		eventType = podModify
	case (e.Mask & inotify.IN_DELETE) > 0:
		eventType = podDelete
	case (e.Mask & inotify.IN_MOVED_FROM) > 0:
		eventType = podDelete
	case (e.Mask & inotify.IN_DELETE_SELF) > 0:
		return fmt.Errorf("the watched path is deleted")
	default:
		// Ignore rest events
		return nil
	}

	s.watchEvents <- &watchEvent{e.Name, eventType}
	return nil
}

func (s *sourceFile) consumeWatchEvent(e *watchEvent) error {
	switch e.eventType {
	case podAdd, podModify:
		if pod, err := s.extractFromFile(e.fileName); err != nil {
			return fmt.Errorf("can't process config file %q: %v", e.fileName, err)
		} else {
			return s.store.Add(pod)
		}
	case podDelete:
		if objKey, keyExist := s.fileKeyMapping[e.fileName]; keyExist {
			pod, podExist, err := s.store.GetByKey(objKey)
			if err != nil {
				return err
			} else if !podExist {
				return fmt.Errorf("the pod with key %s doesn't exist in cache", objKey)
			} else {
				if err = s.store.Delete(pod); err != nil {
					return fmt.Errorf("failed to remove deleted pod from cache: %v", err)
				} else {
					delete(s.fileKeyMapping, e.fileName)
				}
			}
		}
	}
	return nil
}
