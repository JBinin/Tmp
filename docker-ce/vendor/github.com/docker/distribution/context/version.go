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
package context

// WithVersion stores the application version in the context. The new context
// gets a logger to ensure log messages are marked with the application
// version.
func WithVersion(ctx Context, version string) Context {
	ctx = WithValue(ctx, "version", version)
	// push a new logger onto the stack
	return WithLogger(ctx, GetLogger(ctx, "version"))
}

// GetVersion returns the application version from the context. An empty
// string may returned if the version was not set on the context.
func GetVersion(ctx Context) string {
	return GetStringValue(ctx, "version")
}
