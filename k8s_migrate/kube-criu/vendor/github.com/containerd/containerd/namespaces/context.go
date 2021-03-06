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

import (
	"os"

	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	// NamespaceEnvVar is the environment variable key name
	NamespaceEnvVar = "CONTAINERD_NAMESPACE"
	// Default is the name of the default namespace
	Default = "default"
)

type namespaceKey struct{}

// WithNamespace sets a given namespace on the context
func WithNamespace(ctx context.Context, namespace string) context.Context {
	ctx = context.WithValue(ctx, namespaceKey{}, namespace) // set our key for namespace

	// also store on the grpc headers so it gets picked up by any clients that
	// are using this.
	return withGRPCNamespaceHeader(ctx, namespace)
}

// NamespaceFromEnv uses the namespace defined in CONTAINERD_NAMESPACE or
// default
func NamespaceFromEnv(ctx context.Context) context.Context {
	namespace := os.Getenv(NamespaceEnvVar)
	if namespace == "" {
		namespace = Default
	}
	return WithNamespace(ctx, namespace)
}

// Namespace returns the namespace from the context.
//
// The namespace is not guaranteed to be valid.
func Namespace(ctx context.Context) (string, bool) {
	namespace, ok := ctx.Value(namespaceKey{}).(string)
	if !ok {
		return fromGRPCHeader(ctx)
	}

	return namespace, ok
}

// NamespaceRequired returns the valid namepace from the context or an error.
func NamespaceRequired(ctx context.Context) (string, error) {
	namespace, ok := Namespace(ctx)
	if !ok || namespace == "" {
		return "", errors.Wrapf(errdefs.ErrFailedPrecondition, "namespace is required")
	}

	if err := Validate(namespace); err != nil {
		return "", errors.Wrap(err, "namespace validation")
	}

	return namespace, nil
}
