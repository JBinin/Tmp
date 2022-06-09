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
package build

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"golang.org/x/net/context"
)

// Backend abstracts an image builder whose only purpose is to build an image referenced by an imageID.
type Backend interface {
	// Build builds a Docker image referenced by an imageID string.
	//
	// Note: Tagging an image should not be done by a Builder, it should instead be done
	// by the caller.
	//
	// TODO: make this return a reference instead of string
	BuildFromContext(ctx context.Context, src io.ReadCloser, remote string, buildOptions *types.ImageBuildOptions, pg backend.ProgressWriter) (string, error)
}
