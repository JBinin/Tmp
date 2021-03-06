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
// Package endpoints validates regional endpoints for services.
package endpoints

//go:generate go run -tags codegen ../model/cli/gen-endpoints/main.go endpoints.json endpoints_map.go
//go:generate gofmt -s -w endpoints_map.go

import (
	"fmt"
	"regexp"
	"strings"
)

// NormalizeEndpoint takes and endpoint and service API information to return a
// normalized endpoint and signing region.  If the endpoint is not an empty string
// the service name and region will be used to look up the service's API endpoint.
// If the endpoint is provided the scheme will be added if it is not present.
func NormalizeEndpoint(endpoint, serviceName, region string, disableSSL, useDualStack bool) (normEndpoint, signingRegion string) {
	if endpoint == "" {
		return EndpointForRegion(serviceName, region, disableSSL, useDualStack)
	}

	return AddScheme(endpoint, disableSSL), ""
}

// EndpointForRegion returns an endpoint and its signing region for a service and region.
// if the service and region pair are not found endpoint and signingRegion will be empty.
func EndpointForRegion(svcName, region string, disableSSL, useDualStack bool) (endpoint, signingRegion string) {
	dualStackField := ""
	if useDualStack {
		dualStackField = "/dualstack"
	}

	derivedKeys := []string{
		region + "/" + svcName + dualStackField,
		region + "/*" + dualStackField,
		"*/" + svcName + dualStackField,
		"*/*" + dualStackField,
	}

	for _, key := range derivedKeys {
		if val, ok := endpointsMap.Endpoints[key]; ok {
			ep := val.Endpoint
			ep = strings.Replace(ep, "{region}", region, -1)
			ep = strings.Replace(ep, "{service}", svcName, -1)

			endpoint = ep
			signingRegion = val.SigningRegion
			break
		}
	}

	return AddScheme(endpoint, disableSSL), signingRegion
}

// Regular expression to determine if the endpoint string is prefixed with a scheme.
var schemeRE = regexp.MustCompile("^([^:]+)://")

// AddScheme adds the HTTP or HTTPS schemes to a endpoint URL if there is no
// scheme. If disableSSL is true HTTP will be added instead of the default HTTPS.
func AddScheme(endpoint string, disableSSL bool) string {
	if endpoint != "" && !schemeRE.MatchString(endpoint) {
		scheme := "https"
		if disableSSL {
			scheme = "http"
		}
		endpoint = fmt.Sprintf("%s://%s", scheme, endpoint)
	}

	return endpoint
}
