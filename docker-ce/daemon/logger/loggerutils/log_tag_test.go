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
package loggerutils

import (
	"testing"

	"github.com/docker/docker/daemon/logger"
)

func TestParseLogTagDefaultTag(t *testing.T) {
	info := buildContext(map[string]string{})
	tag, e := ParseLogTag(info, "{{.ID}}")
	assertTag(t, e, tag, info.ID())
}

func TestParseLogTag(t *testing.T) {
	info := buildContext(map[string]string{"tag": "{{.ImageName}}/{{.Name}}/{{.ID}}"})
	tag, e := ParseLogTag(info, "{{.ID}}")
	assertTag(t, e, tag, "test-image/test-container/container-ab")
}

func TestParseLogTagEmptyTag(t *testing.T) {
	info := buildContext(map[string]string{})
	tag, e := ParseLogTag(info, "{{.DaemonName}}/{{.ID}}")
	assertTag(t, e, tag, "test-dockerd/container-ab")
}

// Helpers

func buildContext(cfg map[string]string) logger.Info {
	return logger.Info{
		ContainerID:        "container-abcdefghijklmnopqrstuvwxyz01234567890",
		ContainerName:      "/test-container",
		ContainerImageID:   "image-abcdefghijklmnopqrstuvwxyz01234567890",
		ContainerImageName: "test-image",
		Config:             cfg,
		DaemonName:         "test-dockerd",
	}
}

func assertTag(t *testing.T, e error, tag string, expected string) {
	if e != nil {
		t.Fatalf("Error generating tag: %q", e)
	}
	if tag != expected {
		t.Fatalf("Wrong tag: %q, should be %q", tag, expected)
	}
}
