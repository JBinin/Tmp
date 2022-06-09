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
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/image"
)

// TagImage creates the tag specified by newTag, pointing to the image named
// imageName (alternatively, imageName can also be an image ID).
func (daemon *Daemon) TagImage(imageName, repository, tag string) error {
	imageID, err := daemon.GetImageID(imageName)
	if err != nil {
		return err
	}

	newTag, err := reference.ParseNormalizedNamed(repository)
	if err != nil {
		return err
	}
	if tag != "" {
		if newTag, err = reference.WithTag(reference.TrimNamed(newTag), tag); err != nil {
			return err
		}
	}

	return daemon.TagImageWithReference(imageID, newTag)
}

// TagImageWithReference adds the given reference to the image ID provided.
func (daemon *Daemon) TagImageWithReference(imageID image.ID, newTag reference.Named) error {
	if err := daemon.referenceStore.AddTag(newTag, imageID.Digest(), true); err != nil {
		return err
	}

	daemon.LogImageEvent(imageID.String(), reference.FamiliarString(newTag), "tag")
	return nil
}
