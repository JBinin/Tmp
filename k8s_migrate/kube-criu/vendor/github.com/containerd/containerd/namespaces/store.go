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
package namespaces

import "context"

// Store provides introspection about namespaces.
//
// Note that these are slightly different than other objects, which are record
// oriented. A namespace is really just a name and a set of labels. Objects
// that belong to a namespace are returned when the namespace is assigned to a
// given context.
//
//
type Store interface {
	Create(ctx context.Context, namespace string, labels map[string]string) error
	Labels(ctx context.Context, namespace string) (map[string]string, error)
	SetLabel(ctx context.Context, namespace, key, value string) error
	List(ctx context.Context) ([]string, error)

	// Delete removes the namespace. The namespace must be empty to be deleted.
	Delete(ctx context.Context, namespace string) error
}
