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
package plugins

import (
	"errors"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// regression test for deadlock in handlers
func TestPluginAddHandler(t *testing.T) {
	// make a plugin which is pre-activated
	p := &Plugin{activateWait: sync.NewCond(&sync.Mutex{})}
	p.Manifest = &Manifest{Implements: []string{"bananas"}}
	storage.plugins["qwerty"] = p

	testActive(t, p)
	Handle("bananas", func(_ string, _ *Client) {})
	testActive(t, p)
}

func TestPluginWaitBadPlugin(t *testing.T) {
	p := &Plugin{activateWait: sync.NewCond(&sync.Mutex{})}
	p.activateErr = errors.New("some junk happened")
	testActive(t, p)
}

func testActive(t *testing.T, p *Plugin) {
	done := make(chan struct{})
	go func() {
		p.waitActive()
		close(done)
	}()

	select {
	case <-time.After(100 * time.Millisecond):
		_, f, l, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: deadlock in waitActive", filepath.Base(f), l)
	case <-done:
	}

}
