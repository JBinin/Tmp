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
package sample

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// SignURLSample signs URL sample
func SignURLSample() {
	// Create bucket
	bucket, err := GetTestBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Put object
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPPut, 60)
	if err != nil {
		HandleError(err)
	}

	var val = "花间一壶酒，独酌无相亲。 举杯邀明月，对影成三人。"
	err = bucket.PutObjectWithURL(signedURL, strings.NewReader(val))
	if err != nil {
		HandleError(err)
	}

	// Put object with option
	options := []oss.Option{
		oss.Meta("myprop", "mypropval"),
		oss.ContentType("image/tiff"),
	}

	signedURL, err = bucket.SignURL(objectKey, oss.HTTPPut, 60, options...)
	if err != nil {
		HandleError(err)
	}

	err = bucket.PutObjectFromFileWithURL(signedURL, localFile, options...)
	if err != nil {
		HandleError(err)
	}

	// Get object
	signedURL, err = bucket.SignURL(objectKey, oss.HTTPGet, 60)
	if err != nil {
		HandleError(err)
	}

	body, err := bucket.GetObjectWithURL(signedURL)
	if err != nil {
		HandleError(err)
	}
	// Read content
	data, err := ioutil.ReadAll(body)
	body.Close()
	data = data // use data

	err = bucket.GetObjectToFileWithURL(signedURL, "mynewfile-1.jpg")
	if err != nil {
		HandleError(err)
	}

	// Delete the object and bucket
	err = DeleteTestBucketAndObject(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("SignURLSample completed")
}
