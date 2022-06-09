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
// Package urlutil provides helper function to check urls kind.
// It supports http urls, git urls and transport url (tcp://, …)
package urlutil

import (
	"regexp"
	"strings"
)

var (
	validPrefixes = map[string][]string{
		"url":       {"http://", "https://"},
		"git":       {"git://", "github.com/", "git@"},
		"transport": {"tcp://", "tcp+tls://", "udp://", "unix://", "unixgram://"},
	}
	urlPathWithFragmentSuffix = regexp.MustCompile(".git(?:#.+)?$")
)

// IsURL returns true if the provided str is an HTTP(S) URL.
func IsURL(str string) bool {
	return checkURL(str, "url")
}

// IsGitURL returns true if the provided str is a git repository URL.
func IsGitURL(str string) bool {
	if IsURL(str) && urlPathWithFragmentSuffix.MatchString(str) {
		return true
	}
	return checkURL(str, "git")
}

// IsGitTransport returns true if the provided str is a git transport by inspecting
// the prefix of the string for known protocols used in git.
func IsGitTransport(str string) bool {
	return IsURL(str) || strings.HasPrefix(str, "git://") || strings.HasPrefix(str, "git@")
}

// IsTransportURL returns true if the provided str is a transport (tcp, tcp+tls, udp, unix) URL.
func IsTransportURL(str string) bool {
	return checkURL(str, "transport")
}

func checkURL(str, kind string) bool {
	for _, prefix := range validPrefixes[kind] {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}
