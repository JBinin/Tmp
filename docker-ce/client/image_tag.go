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
package client

import (
	"net/url"

	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// ImageTag tags an image in the docker host
func (cli *Client) ImageTag(ctx context.Context, source, target string) error {
	if _, err := reference.ParseNormalizedNamed(source); err != nil {
		return errors.Wrapf(err, "Error parsing reference: %q is not a valid repository/tag", source)
	}

	ref, err := reference.ParseNormalizedNamed(target)
	if err != nil {
		return errors.Wrapf(err, "Error parsing reference: %q is not a valid repository/tag", target)
	}

	if _, isCanonical := ref.(reference.Canonical); isCanonical {
		return errors.New("refusing to create a tag with a digest reference")
	}

	ref = reference.TagNameOnly(ref)

	query := url.Values{}
	query.Set("repo", reference.FamiliarName(ref))
	if tagged, ok := ref.(reference.Tagged); ok {
		query.Set("tag", tagged.Tag())
	}

	resp, err := cli.post(ctx, "/images/"+source+"/tag", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
