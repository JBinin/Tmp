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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/testutil/assert"
	"golang.org/x/net/context"
)

func TestImagesPruneError(t *testing.T) {
	client := &Client{
		client:  newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
		version: "1.25",
	}

	filters := filters.NewArgs()

	_, err := client.ImagesPrune(context.Background(), filters)
	assert.Error(t, err, "Error response from daemon: Server error")
}

func TestImagesPrune(t *testing.T) {
	expectedURL := "/v1.25/images/prune"

	danglingFilters := filters.NewArgs()
	danglingFilters.Add("dangling", "true")

	noDanglingFilters := filters.NewArgs()
	noDanglingFilters.Add("dangling", "false")

	listCases := []struct {
		filters             filters.Args
		expectedQueryParams map[string]string
	}{
		{
			filters: filters.Args{},
			expectedQueryParams: map[string]string{
				"until":   "",
				"filter":  "",
				"filters": "",
			},
		},
		{
			filters: danglingFilters,
			expectedQueryParams: map[string]string{
				"until":   "",
				"filter":  "",
				"filters": `{"dangling":{"true":true}}`,
			},
		},
		{
			filters: noDanglingFilters,
			expectedQueryParams: map[string]string{
				"until":   "",
				"filter":  "",
				"filters": `{"dangling":{"false":true}}`,
			},
		},
	}
	for _, listCase := range listCases {
		client := &Client{
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				if !strings.HasPrefix(req.URL.Path, expectedURL) {
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				}
				query := req.URL.Query()
				for key, expected := range listCase.expectedQueryParams {
					actual := query.Get(key)
					assert.Equal(t, actual, expected)
				}
				content, err := json.Marshal(types.ImagesPruneReport{
					ImagesDeleted: []types.ImageDeleteResponseItem{
						{
							Deleted: "image_id1",
						},
						{
							Deleted: "image_id2",
						},
					},
					SpaceReclaimed: 9999,
				})
				if err != nil {
					return nil, err
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(content)),
				}, nil
			}),
			version: "1.25",
		}

		report, err := client.ImagesPrune(context.Background(), listCase.filters)
		assert.NilError(t, err)
		assert.Equal(t, len(report.ImagesDeleted), 2)
		assert.Equal(t, report.SpaceReclaimed, uint64(9999))
	}
}
