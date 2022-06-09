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
package daemon

import (
	"io"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/distribution"
	progressutils "github.com/docker/docker/distribution/utils"
	"github.com/docker/docker/pkg/progress"
	"golang.org/x/net/context"
)

// PushImage initiates a push operation on the repository named localName.
func (daemon *Daemon) PushImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return err
	}
	if tag != "" {
		// Push by digest is not supported, so only tags are supported.
		ref, err = reference.WithTag(ref, tag)
		if err != nil {
			return err
		}
	}

	// Include a buffer so that slow client connections don't affect
	// transfer performance.
	progressChan := make(chan progress.Progress, 100)

	writesDone := make(chan struct{})

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() {
		progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
		close(writesDone)
	}()

	imagePushConfig := &distribution.ImagePushConfig{
		Config: distribution.Config{
			MetaHeaders:      metaHeaders,
			AuthConfig:       authConfig,
			ProgressOutput:   progress.ChanOutput(progressChan),
			RegistryService:  daemon.RegistryService,
			ImageEventLogger: daemon.LogImageEvent,
			MetadataStore:    daemon.distributionMetadataStore,
			ImageStore:       distribution.NewImageConfigStoreFromStore(daemon.imageStore),
			ReferenceStore:   daemon.referenceStore,
		},
		ConfigMediaType: schema2.MediaTypeImageConfig,
		LayerStore:      distribution.NewLayerProviderFromStore(daemon.layerStore),
		TrustKey:        daemon.trustKey,
		UploadManager:   daemon.uploadManager,
	}

	err = distribution.Push(ctx, ref, imagePushConfig)
	close(progressChan)
	<-writesDone
	return err
}
