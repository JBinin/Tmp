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
package log

import (
	"path"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	// G is an alias for GetLogger.
	//
	// We may want to define this locally to a package to get package tagged log
	// messages.
	G = GetLogger

	// L is an alias for the the standard logger.
	L = logrus.NewEntry(logrus.StandardLogger())
)

type (
	loggerKey struct{}
	moduleKey struct{}
)

// WithLogger returns a new context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithFields returns a new context with added fields to logger.
func WithFields(ctx context.Context, fields logrus.Fields) context.Context {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		logger = L
	}
	return WithLogger(ctx, logger.(*logrus.Entry).WithFields(fields))
}

// WithField is convenience wrapper around WithFields.
func WithField(ctx context.Context, key, value string) context.Context {
	return WithFields(ctx, logrus.Fields{key: value})
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		return L
	}

	return logger.(*logrus.Entry)
}

// WithModule adds the module to the context, appending it with a slash if a
// module already exists. A module is just a roughly correlated defined by the
// call tree for a given context.
//
// As an example, we might have a "node" module already part of a context. If
// this function is called with "tls", the new value of module will be
// "node/tls".
//
// Modules represent the call path. If the new module and last module are the
// same, a new module entry will not be created. If the new module and old
// older module are the same but separated by other modules, the cycle will be
// represented by the module path.
func WithModule(ctx context.Context, module string) context.Context {
	parent := GetModulePath(ctx)

	if parent != "" {
		// don't re-append module when module is the same.
		if path.Base(parent) == module {
			return ctx
		}

		module = path.Join(parent, module)
	}

	ctx = WithLogger(ctx, GetLogger(ctx).WithField("module", module))
	return context.WithValue(ctx, moduleKey{}, module)
}

// GetModulePath returns the module path for the provided context. If no module
// is set, an empty string is returned.
func GetModulePath(ctx context.Context) string {
	module := ctx.Value(moduleKey{})
	if module == nil {
		return ""
	}

	return module.(string)
}
