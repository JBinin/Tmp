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
package debug

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestEnable(t *testing.T) {
	defer func() {
		os.Setenv("DEBUG", "")
		logrus.SetLevel(logrus.InfoLevel)
	}()
	Enable()
	if os.Getenv("DEBUG") != "1" {
		t.Fatalf("expected DEBUG=1, got %s\n", os.Getenv("DEBUG"))
	}
	if logrus.GetLevel() != logrus.DebugLevel {
		t.Fatalf("expected log level %v, got %v\n", logrus.DebugLevel, logrus.GetLevel())
	}
}

func TestDisable(t *testing.T) {
	Disable()
	if os.Getenv("DEBUG") != "" {
		t.Fatalf("expected DEBUG=\"\", got %s\n", os.Getenv("DEBUG"))
	}
	if logrus.GetLevel() != logrus.InfoLevel {
		t.Fatalf("expected log level %v, got %v\n", logrus.InfoLevel, logrus.GetLevel())
	}
}

func TestEnabled(t *testing.T) {
	Enable()
	if !IsEnabled() {
		t.Fatal("expected debug enabled, got false")
	}
	Disable()
	if IsEnabled() {
		t.Fatal("expected debug disabled, got true")
	}
}
