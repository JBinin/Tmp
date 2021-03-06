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
package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func TestImageTagError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	err := client.ImageTag(context.Background(), "image_id", "repo:tag")
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

// Note: this is not testing all the InvalidReference as it's the responsibility
// of distribution/reference package.
func TestImageTagInvalidReference(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	err := client.ImageTag(context.Background(), "image_id", "aa/asdf$$^/aa")
	if err == nil || err.Error() != `Error parsing reference: "aa/asdf$$^/aa" is not a valid repository/tag: invalid reference format` {
		t.Fatalf("expected ErrReferenceInvalidFormat, got %v", err)
	}
}

func TestImageTagInvalidSourceImageName(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	err := client.ImageTag(context.Background(), "invalid_source_image_name_", "repo:tag")
	if err == nil || err.Error() != "Error parsing reference: \"invalid_source_image_name_\" is not a valid repository/tag: invalid reference format" {
		t.Fatalf("expected Parsing Reference Error, got %v", err)
	}
}

func TestImageTag(t *testing.T) {
	expectedURL := "/images/image_id/tag"
	tagCases := []struct {
		reference           string
		expectedQueryParams map[string]string
	}{
		{
			reference: "repository:tag1",
			expectedQueryParams: map[string]string{
				"repo": "repository",
				"tag":  "tag1",
			},
		}, {
			reference: "another_repository:latest",
			expectedQueryParams: map[string]string{
				"repo": "another_repository",
				"tag":  "latest",
			},
		}, {
			reference: "another_repository",
			expectedQueryParams: map[string]string{
				"repo": "another_repository",
				"tag":  "latest",
			},
		}, {
			reference: "test/another_repository",
			expectedQueryParams: map[string]string{
				"repo": "test/another_repository",
				"tag":  "latest",
			},
		}, {
			reference: "test/another_repository:tag1",
			expectedQueryParams: map[string]string{
				"repo": "test/another_repository",
				"tag":  "tag1",
			},
		}, {
			reference: "test/test/another_repository:tag1",
			expectedQueryParams: map[string]string{
				"repo": "test/test/another_repository",
				"tag":  "tag1",
			},
		}, {
			reference: "test:5000/test/another_repository:tag1",
			expectedQueryParams: map[string]string{
				"repo": "test:5000/test/another_repository",
				"tag":  "tag1",
			},
		}, {
			reference: "test:5000/test/another_repository",
			expectedQueryParams: map[string]string{
				"repo": "test:5000/test/another_repository",
				"tag":  "latest",
			},
		},
	}
	for _, tagCase := range tagCases {
		client := &Client{
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				if !strings.HasPrefix(req.URL.Path, expectedURL) {
					return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
				}
				if req.Method != "POST" {
					return nil, fmt.Errorf("expected POST method, got %s", req.Method)
				}
				query := req.URL.Query()
				for key, expected := range tagCase.expectedQueryParams {
					actual := query.Get(key)
					if actual != expected {
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					}
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			}),
		}
		err := client.ImageTag(context.Background(), "image_id", tagCase.reference)
		if err != nil {
			t.Fatal(err)
		}
	}
}
