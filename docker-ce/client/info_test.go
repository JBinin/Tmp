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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func TestInfoServerError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}
	_, err := client.Info(context.Background())
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestInfoInvalidResponseJSONError(t *testing.T) {
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("invalid json"))),
			}, nil
		}),
	}
	_, err := client.Info(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("expected a 'invalid character' error, got %v", err)
	}
}

func TestInfo(t *testing.T) {
	expectedURL := "/info"
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			info := &types.Info{
				ID:         "daemonID",
				Containers: 3,
			}
			b, err := json.Marshal(info)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			}, nil
		}),
	}

	info, err := client.Info(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if info.ID != "daemonID" {
		t.Fatalf("expected daemonID, got %s", info.ID)
	}

	if info.Containers != 3 {
		t.Fatalf("expected 3 containers, got %d", info.Containers)
	}
}
