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
package swarm

import "os"

// Secret represents a secret.
type Secret struct {
	ID string
	Meta
	Spec SecretSpec
}

// SecretSpec represents a secret specification from a secret in swarm
type SecretSpec struct {
	Annotations
	Data   []byte  `json:",omitempty"`
	Driver *Driver `json:",omitempty"` // name of the secrets driver used to fetch the secret's value from an external secret store
}

// SecretReferenceFileTarget is a file target in a secret reference
type SecretReferenceFileTarget struct {
	Name string
	UID  string
	GID  string
	Mode os.FileMode
}

// SecretReference is a reference to a secret in swarm
type SecretReference struct {
	File       *SecretReferenceFileTarget
	SecretID   string
	SecretName string
}
