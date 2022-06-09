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

package journald

import (
	"testing"
)

func TestSanitizeKeyMod(t *testing.T) {
	entries := map[string]string{
		"io.kubernetes.pod.name":      "IO_KUBERNETES_POD_NAME",
		"io?.kubernetes.pod.name":     "IO__KUBERNETES_POD_NAME",
		"?io.kubernetes.pod.name":     "IO_KUBERNETES_POD_NAME",
		"io123.kubernetes.pod.name":   "IO123_KUBERNETES_POD_NAME",
		"_io123.kubernetes.pod.name":  "IO123_KUBERNETES_POD_NAME",
		"__io123_kubernetes.pod.name": "IO123_KUBERNETES_POD_NAME",
	}
	for k, v := range entries {
		if sanitizeKeyMod(k) != v {
			t.Fatalf("Failed to sanitize %s, got %s, expected %s", k, sanitizeKeyMod(k), v)
		}
	}
}
