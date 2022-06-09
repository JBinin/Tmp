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
package plugin

import (
	"github.com/gogo/protobuf/proto"
	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// DeepcopyEnabled returns true if deepcopy is enabled for the descriptor.
func DeepcopyEnabled(options *google_protobuf.MessageOptions) bool {
	return proto.GetBoolExtension(options, E_Deepcopy, true)
}
