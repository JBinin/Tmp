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

package i18n

import (
	"os"
	"testing"
)

func TestTranslation(t *testing.T) {
	err := LoadTranslations("test", func() string { return "default" })
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := T("test_string")
	if result != "foo" {
		t.Errorf("expected: %s, saw: %s", "foo", result)
	}
}

func TestTranslationPlural(t *testing.T) {
	err := LoadTranslations("test", func() string { return "default" })
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := T("test_plural", 3)
	if result != "there were 3 items" {
		t.Errorf("expected: %s, saw: %s", "there were 3 items", result)
	}

	result = T("test_plural", 1)
	if result != "there was 1 item" {
		t.Errorf("expected: %s, saw: %s", "there was 1 item", result)
	}
}

func TestTranslationEnUSEnv(t *testing.T) {
	os.Setenv("LANG", "en_US.UTF-8")
	err := LoadTranslations("test", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := T("test_string")
	if result != "baz" {
		t.Errorf("expected: %s, saw: %s", "baz", result)
	}
}
