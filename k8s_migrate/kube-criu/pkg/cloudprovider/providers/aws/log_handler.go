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
Copyright 2015 The Kubernetes Authors.

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
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/golang/glog"
)

// Handler for aws-sdk-go that logs all requests
func awsHandlerLogger(req *request.Request) {
	service, name := awsServiceAndName(req)
	glog.V(4).Infof("AWS request: %s %s", service, name)
}

func awsSendHandlerLogger(req *request.Request) {
	service, name := awsServiceAndName(req)
	glog.V(4).Infof("AWS API Send: %s %s %v %v", service, name, req.Operation, req.Params)
}

func awsValidateResponseHandlerLogger(req *request.Request) {
	service, name := awsServiceAndName(req)
	glog.V(4).Infof("AWS API ValidateResponse: %s %s %v %v %s", service, name, req.Operation, req.Params, req.HTTPResponse.Status)
}

func awsServiceAndName(req *request.Request) (string, string) {
	service := req.ClientInfo.ServiceName

	name := "?"
	if req.Operation != nil {
		name = req.Operation.Name
	}
	return service, name
}
