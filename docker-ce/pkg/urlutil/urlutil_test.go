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
package urlutil

import "testing"

var (
	gitUrls = []string{
		"git://github.com/docker/docker",
		"git@github.com:docker/docker.git",
		"git@bitbucket.org:atlassianlabs/atlassian-docker.git",
		"https://github.com/docker/docker.git",
		"http://github.com/docker/docker.git",
		"http://github.com/docker/docker.git#branch",
		"http://github.com/docker/docker.git#:dir",
	}
	incompleteGitUrls = []string{
		"github.com/docker/docker",
	}
	invalidGitUrls = []string{
		"http://github.com/docker/docker.git:#branch",
	}
	transportUrls = []string{
		"tcp://example.com",
		"tcp+tls://example.com",
		"udp://example.com",
		"unix:///example",
		"unixgram:///example",
	}
)

func TestValidGitTransport(t *testing.T) {
	for _, url := range gitUrls {
		if IsGitTransport(url) == false {
			t.Fatalf("%q should be detected as valid Git prefix", url)
		}
	}

	for _, url := range incompleteGitUrls {
		if IsGitTransport(url) == true {
			t.Fatalf("%q should not be detected as valid Git prefix", url)
		}
	}
}

func TestIsGIT(t *testing.T) {
	for _, url := range gitUrls {
		if IsGitURL(url) == false {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}

	for _, url := range incompleteGitUrls {
		if IsGitURL(url) == false {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}

	for _, url := range invalidGitUrls {
		if IsGitURL(url) == true {
			t.Fatalf("%q should not be detected as valid Git prefix", url)
		}
	}
}

func TestIsTransport(t *testing.T) {
	for _, url := range transportUrls {
		if IsTransportURL(url) == false {
			t.Fatalf("%q should be detected as valid Transport url", url)
		}
	}
}
