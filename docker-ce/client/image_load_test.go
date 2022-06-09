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

func TestImageLoadError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	_, err := client.ImageLoad(context.Background(), nil, true)
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestImageLoad(t *testing.T) {
	expectedURL := "/images/load"
	expectedInput := "inputBody"
	expectedOutput := "outputBody"
	loadCases := []struct {
		quiet                bool
		responseContentType  string
		expectedResponseJSON bool
		expectedQueryParams  map[string]string
	}{
		{
			quiet:                false,
			responseContentType:  "text/plain",
			expectedResponseJSON: false,
			expectedQueryParams: map[string]string{
				"quiet": "0",
			},
		},
		{
			quiet:                true,
			responseContentType:  "application/json",
			expectedResponseJSON: true,
			expectedQueryParams: map[string]string{
				"quiet": "1",
			},
		},
	}
	for _, loadCase := range loadCases {
		client := &Client{
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				if !strings.HasPrefix(req.URL.Path, expectedURL) {
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				}
				contentType := req.Header.Get("Content-Type")
				if contentType != "application/x-tar" {
					return nil, fmt.Errorf("content-type not set in URL headers properly. Expected 'application/x-tar', got %s", contentType)
				}
				query := req.URL.Query()
				for key, expected := range loadCase.expectedQueryParams {
					actual := query.Get(key)
					if actual != expected {
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					}
				}
				headers := http.Header{}
				headers.Add("Content-Type", loadCase.responseContentType)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(expectedOutput))),
					Header:     headers,
				}, nil
			}),
		}

		input := bytes.NewReader([]byte(expectedInput))
		imageLoadResponse, err := client.ImageLoad(context.Background(), input, loadCase.quiet)
		if err != nil {
			t.Fatal(err)
		}
		if imageLoadResponse.JSON != loadCase.expectedResponseJSON {
			t.Fatalf("expected a JSON response, was not.")
		}
		body, err := ioutil.ReadAll(imageLoadResponse.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != expectedOutput {
			t.Fatalf("expected %s, got %s", expectedOutput, string(body))
		}
	}
}
