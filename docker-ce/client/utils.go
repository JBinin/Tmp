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
	"github.com/docker/docker/api/types/filters"
	"net/url"
	"regexp"
)

var headerRegexp = regexp.MustCompile(`\ADocker/.+\s\((.+)\)\z`)

// getDockerOS returns the operating system based on the server header from the daemon.
func getDockerOS(serverHeader string) string {
	var osType string
	matches := headerRegexp.FindStringSubmatch(serverHeader)
	if len(matches) > 0 {
		osType = matches[1]
	}
	return osType
}

// getFiltersQuery returns a url query with "filters" query term, based on the
// filters provided.
func getFiltersQuery(f filters.Args) (url.Values, error) {
	query := url.Values{}
	if f.Len() > 0 {
		filterJSON, err := filters.ToParam(f)
		if err != nil {
			return query, err
		}
		query.Set("filters", filterJSON)
	}
	return query, nil
}
