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
package manifest

// Versioned provides a struct with the manifest schemaVersion and mediaType.
// Incoming content with unknown schema version can be decoded against this
// struct to check the version.
type Versioned struct {
	// SchemaVersion is the image manifest schema that this image follows
	SchemaVersion int `json:"schemaVersion"`

	// MediaType is the media type of this schema.
	MediaType string `json:"mediaType,omitempty"`
}
