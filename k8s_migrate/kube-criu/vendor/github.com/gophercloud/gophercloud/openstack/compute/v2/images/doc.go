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
/*
Package images provides information and interaction with the images through
the OpenStack Compute service.

This API is deprecated and will be removed from a future version of the Nova
API service.

An image is a collection of files used to create or rebuild a server.
Operators provide a number of pre-built OS images by default. You may also
create custom images from cloud servers you have launched.

Example to List Images

	listOpts := images.ListOpts{
		Limit: 2,
	}

	allPages, err := images.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		panic(err)
	}

	for _, image := range allImages {
		fmt.Printf("%+v\n", image)
	}
*/
package images
