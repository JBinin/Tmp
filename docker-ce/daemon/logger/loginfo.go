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
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Info provides enough information for a logging driver to do its function.
type Info struct {
	Config              map[string]string
	ContainerID         string
	ContainerName       string
	ContainerEntrypoint string
	ContainerArgs       []string
	ContainerImageID    string
	ContainerImageName  string
	ContainerCreated    time.Time
	ContainerEnv        []string
	ContainerLabels     map[string]string
	LogPath             string
	DaemonName          string
}

// ExtraAttributes returns the user-defined extra attributes (labels,
// environment variables) in key-value format. This can be used by log drivers
// that support metadata to add more context to a log.
func (info *Info) ExtraAttributes(keyMod func(string) string) map[string]string {
	extra := make(map[string]string)
	labels, ok := info.Config["labels"]
	if ok && len(labels) > 0 {
		for _, l := range strings.Split(labels, ",") {
			if v, ok := info.ContainerLabels[l]; ok {
				if keyMod != nil {
					l = keyMod(l)
				}
				extra[l] = v
			}
		}
	}

	env, ok := info.Config["env"]
	if ok && len(env) > 0 {
		envMapping := make(map[string]string)
		for _, e := range info.ContainerEnv {
			if kv := strings.SplitN(e, "=", 2); len(kv) == 2 {
				envMapping[kv[0]] = kv[1]
			}
		}
		for _, l := range strings.Split(env, ",") {
			if v, ok := envMapping[l]; ok {
				if keyMod != nil {
					l = keyMod(l)
				}
				extra[l] = v
			}
		}
	}

	return extra
}

// Hostname returns the hostname from the underlying OS.
func (info *Info) Hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("logger: can not resolve hostname: %v", err)
	}
	return hostname, nil
}

// Command returns the command that the container being logged was
// started with. The Entrypoint is prepended to the container
// arguments.
func (info *Info) Command() string {
	terms := []string{info.ContainerEntrypoint}
	terms = append(terms, info.ContainerArgs...)
	command := strings.Join(terms, " ")
	return command
}

// ID Returns the Container ID shortened to 12 characters.
func (info *Info) ID() string {
	return info.ContainerID[:12]
}

// FullID is an alias of ContainerID.
func (info *Info) FullID() string {
	return info.ContainerID
}

// Name returns the ContainerName without a preceding '/'.
func (info *Info) Name() string {
	return strings.TrimPrefix(info.ContainerName, "/")
}

// ImageID returns the ContainerImageID shortened to 12 characters.
func (info *Info) ImageID() string {
	return info.ContainerImageID[:12]
}

// ImageFullID is an alias of ContainerImageID.
func (info *Info) ImageFullID() string {
	return info.ContainerImageID
}

// ImageName is an alias of ContainerImageName
func (info *Info) ImageName() string {
	return info.ContainerImageName
}
