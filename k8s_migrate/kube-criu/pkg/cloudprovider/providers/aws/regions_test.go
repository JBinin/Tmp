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
/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"testing"
)

// TestRegions does basic checking of region verification / addition
func TestRegions(t *testing.T) {
	RecognizeWellKnownRegions()

	tests := []struct {
		Add            string
		Lookup         string
		ExpectIsRegion bool
	}{
		{
			Lookup:         "us-east-1",
			ExpectIsRegion: true,
		},
		{
			Lookup:         "us-east-1a",
			ExpectIsRegion: false,
		},
		{
			Add:            "us-test-1",
			Lookup:         "us-east-1",
			ExpectIsRegion: true,
		},
		{
			Lookup:         "us-test-1",
			ExpectIsRegion: true,
		},
		{
			Add:            "us-test-1",
			Lookup:         "us-test-1",
			ExpectIsRegion: true,
		},
	}

	for _, test := range tests {
		if test.Add != "" {
			RecognizeRegion(test.Add)
		}

		if test.Lookup != "" {
			if isRegionValid(test.Lookup) != test.ExpectIsRegion {
				t.Fatalf("region valid mismatch: %q", test.Lookup)
			}
		}
	}
}

// TestRecognizesNewRegion verifies that we see a region from metadata, we recognize it as valid
func TestRecognizesNewRegion(t *testing.T) {
	region := "us-testrecognizesnewregion-1"
	if isRegionValid(region) {
		t.Fatalf("region already valid: %q", region)
	}

	awsServices := NewFakeAWSServices(TestClusterId).WithAz(region + "a")
	_, err := newAWSCloud(CloudConfig{}, awsServices)
	if err != nil {
		t.Errorf("error building AWS cloud: %v", err)
	}

	if !isRegionValid(region) {
		t.Fatalf("newly discovered region not valid: %q", region)
	}
}
