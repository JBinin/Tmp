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
package container

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/signal"
)

func TestContainerStopSignal(t *testing.T) {
	c := &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{},
		},
	}

	def, err := signal.ParseSignal(signal.DefaultStopSignal)
	if err != nil {
		t.Fatal(err)
	}

	s := c.StopSignal()
	if s != int(def) {
		t.Fatalf("Expected %v, got %v", def, s)
	}

	c = &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{StopSignal: "SIGKILL"},
		},
	}
	s = c.StopSignal()
	if s != 9 {
		t.Fatalf("Expected 9, got %v", s)
	}
}

func TestContainerStopTimeout(t *testing.T) {
	c := &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{},
		},
	}

	s := c.StopTimeout()
	if s != DefaultStopTimeout {
		t.Fatalf("Expected %v, got %v", DefaultStopTimeout, s)
	}

	stopTimeout := 15
	c = &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{StopTimeout: &stopTimeout},
		},
	}
	s = c.StopSignal()
	if s != 15 {
		t.Fatalf("Expected 15, got %v", s)
	}
}
