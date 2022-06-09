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
package filenotify

import "github.com/fsnotify/fsnotify"

// fsNotifyWatcher wraps the fsnotify package to satisfy the FileNotifer interface
type fsNotifyWatcher struct {
	*fsnotify.Watcher
}

// Events returns the fsnotify event channel receiver
func (w *fsNotifyWatcher) Events() <-chan fsnotify.Event {
	return w.Watcher.Events
}

// Errors returns the fsnotify error channel receiver
func (w *fsNotifyWatcher) Errors() <-chan error {
	return w.Watcher.Errors
}
