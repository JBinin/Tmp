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

package generated

import "github.com/golang/glog"

/*
ReadOrDie reads a file from gobindata.
Relies heavily on the successful generation of build artifacts as per the go:generate directives above.
*/
func ReadOrDie(filePath string) []byte {

	fileBytes, err := Asset(filePath)
	if err != nil {
		gobindataMsg := "An error occurred, possibly gobindata doesn't know about the file you're opening. For questions on maintaining gobindata, contact the sig-testing group."
		glog.Infof("Available gobindata files: %v ", AssetNames())
		glog.Fatalf("Failed opening %v , with error %v.  %v.", filePath, err, gobindataMsg)
	}
	return fileBytes
}
