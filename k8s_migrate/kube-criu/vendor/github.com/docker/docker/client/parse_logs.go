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

// parse_logs.go contains utility helpers for getting information out of docker
// log lines. really, it only contains ParseDetails right now. maybe in the
// future there will be some desire to parse log messages back into a struct?
// that would go here if we did

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// ParseLogDetails takes a details string of key value pairs in the form
// "k=v,l=w", where the keys and values are url query escaped, and each pair
// is separated by a comma, returns a map. returns an error if the details
// string is not in a valid format
// the exact form of details encoding is implemented in
// api/server/httputils/write_log_stream.go
func ParseLogDetails(details string) (map[string]string, error) {
	pairs := strings.Split(details, ",")
	detailsMap := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		p := strings.SplitN(pair, "=", 2)
		// if there is no equals sign, we will only get 1 part back
		if len(p) != 2 {
			return nil, errors.New("invalid details format")
		}
		k, err := url.QueryUnescape(p[0])
		if err != nil {
			return nil, err
		}
		v, err := url.QueryUnescape(p[1])
		if err != nil {
			return nil, err
		}
		detailsMap[k] = v
	}
	return detailsMap, nil
}
