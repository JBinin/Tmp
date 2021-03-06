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
package oss

import (
	"bytes"
	"fmt"
	"os"
	"time"

	. "gopkg.in/check.v1"
)

type OssDownloadSuite struct {
	client *Client
	bucket *Bucket
}

var _ = Suite(&OssDownloadSuite{})

// SetUpSuite runs once when the suite starts running
func (s *OssDownloadSuite) SetUpSuite(c *C) {
	client, err := New(endpoint, accessID, accessKey)
	c.Assert(err, IsNil)
	s.client = client

	s.client.CreateBucket(bucketName)
	time.Sleep(5 * time.Second)

	bucket, err := s.client.Bucket(bucketName)
	c.Assert(err, IsNil)
	s.bucket = bucket

	testLogger.Println("test download started")
}

// TearDownSuite runs before each test or benchmark starts running
func (s *OssDownloadSuite) TearDownSuite(c *C) {
	// Delete part
	lmur, err := s.bucket.ListMultipartUploads()
	c.Assert(err, IsNil)

	for _, upload := range lmur.Uploads {
		var imur = InitiateMultipartUploadResult{Bucket: s.bucket.BucketName,
			Key: upload.Key, UploadID: upload.UploadID}
		err = s.bucket.AbortMultipartUpload(imur)
		c.Assert(err, IsNil)
	}

	// Delete objects
	lor, err := s.bucket.ListObjects()
	c.Assert(err, IsNil)

	for _, object := range lor.Objects {
		err = s.bucket.DeleteObject(object.Key)
		c.Assert(err, IsNil)
	}

	testLogger.Println("test download completed")
}

// SetUpTest runs after each test or benchmark runs
func (s *OssDownloadSuite) SetUpTest(c *C) {
	err := removeTempFiles("../oss", ".jpg")
	c.Assert(err, IsNil)
}

// TearDownTest runs once after all tests or benchmarks have finished running
func (s *OssDownloadSuite) TearDownTest(c *C) {
	err := removeTempFiles("../oss", ".jpg")
	c.Assert(err, IsNil)

	err = removeTempFiles("../oss", ".temp")
	c.Assert(err, IsNil)
}

// TestDownloadRoutineWithoutRecovery multipart downloads without checkpoint
func (s *OssDownloadSuite) TestDownloadRoutineWithoutRecovery(c *C) {
	objectName := objectNamePrefix + "tdrwr"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file.jpg"

	// Upload a file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	// Download the file by default
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024)
	c.Assert(err, IsNil)

	// Check
	eq, err := compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Use 2 coroutines to download the file and total parts count is 5
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(2))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Use 5 coroutines to download the file and the total parts count is 5.
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(5))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Use 10 coroutines to download the file and the total parts count is 5.
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(10))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	err = s.bucket.DeleteObject(objectName)
	c.Assert(err, IsNil)
}

// DownErrorHooker requests hook by downloadPart
func DownErrorHooker(part downloadPart) error {
	if part.Index == 4 {
		time.Sleep(time.Second)
		return fmt.Errorf("ErrorHooker")
	}
	return nil
}

// TestDownloadRoutineWithRecovery multi-routine resumable download
func (s *OssDownloadSuite) TestDownloadRoutineWithRecovery(c *C) {
	objectName := objectNamePrefix + "tdrtr"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file-2.jpg"

	// Upload a file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	// Download a file with default checkpoint
	downloadPartHooker = DownErrorHooker
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, ""))
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "ErrorHooker")
	downloadPartHooker = defaultDownloadPartHook

	// Check
	dcp := downloadCheckpoint{}
	err = dcp.load(newFile + ".cp")
	c.Assert(err, IsNil)
	c.Assert(dcp.Magic, Equals, downloadCpMagic)
	c.Assert(len(dcp.MD5), Equals, len("LC34jZU5xK4hlxi3Qn3XGQ=="))
	c.Assert(dcp.FilePath, Equals, newFile)
	c.Assert(dcp.ObjStat.Size, Equals, int64(482048))
	c.Assert(len(dcp.ObjStat.LastModified) > 0, Equals, true)
	c.Assert(dcp.ObjStat.Etag, Equals, "\"2351E662233817A7AE974D8C5B0876DD-5\"")
	c.Assert(dcp.Object, Equals, objectName)
	c.Assert(len(dcp.Parts), Equals, 5)
	c.Assert(len(dcp.todoParts()), Equals, 1)

	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, ""))
	c.Assert(err, IsNil)

	err = dcp.load(newFile + ".cp")
	c.Assert(err, NotNil)

	eq, err := compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Resumable download with checkpoint
	os.Remove(newFile)
	downloadPartHooker = DownErrorHooker
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, objectName+".cp"))
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "ErrorHooker")
	downloadPartHooker = defaultDownloadPartHook

	// Check
	dcp = downloadCheckpoint{}
	err = dcp.load(objectName + ".cp")
	c.Assert(err, IsNil)
	c.Assert(dcp.Magic, Equals, downloadCpMagic)
	c.Assert(len(dcp.MD5), Equals, len("LC34jZU5xK4hlxi3Qn3XGQ=="))
	c.Assert(dcp.FilePath, Equals, newFile)
	c.Assert(dcp.ObjStat.Size, Equals, int64(482048))
	c.Assert(len(dcp.ObjStat.LastModified) > 0, Equals, true)
	c.Assert(dcp.ObjStat.Etag, Equals, "\"2351E662233817A7AE974D8C5B0876DD-5\"")
	c.Assert(dcp.Object, Equals, objectName)
	c.Assert(len(dcp.Parts), Equals, 5)
	c.Assert(len(dcp.todoParts()), Equals, 1)

	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, objectName+".cp"))
	c.Assert(err, IsNil)

	err = dcp.load(objectName + ".cp")
	c.Assert(err, NotNil)

	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Resumable download with checkpoint at a time. No error is expected in the download procedure.
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, ""))
	c.Assert(err, IsNil)

	err = dcp.load(newFile + ".cp")
	c.Assert(err, NotNil)

	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Resumable download with checkpoint at a time. No error is expected in the download procedure.
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(10), Checkpoint(true, ""))
	c.Assert(err, IsNil)

	err = dcp.load(newFile + ".cp")
	c.Assert(err, NotNil)

	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	err = s.bucket.DeleteObject(objectName)
	c.Assert(err, IsNil)
}

// TestDownloadOption options
func (s *OssDownloadSuite) TestDownloadOption(c *C) {
	objectName := objectNamePrefix + "tdmo"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file-3.jpg"

	// Upload the file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	meta, err := s.bucket.GetObjectDetailedMeta(objectName)
	c.Assert(err, IsNil)

	// IfMatch
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(3), IfMatch(meta.Get("Etag")))
	c.Assert(err, IsNil)

	eq, err := compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// IfNoneMatch
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(3), IfNoneMatch(meta.Get("Etag")))
	c.Assert(err, NotNil)

	// IfMatch
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(3), Checkpoint(true, ""), IfMatch(meta.Get("Etag")))
	c.Assert(err, IsNil)

	eq, err = compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// IfNoneMatch
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(3), Checkpoint(true, ""), IfNoneMatch(meta.Get("Etag")))
	c.Assert(err, NotNil)
}

// TestDownloadObjectChange tests the file is updated during the upload
func (s *OssDownloadSuite) TestDownloadObjectChange(c *C) {
	objectName := objectNamePrefix + "tdloc"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file-4.jpg"

	// Upload a file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	// Download with default checkpoint
	downloadPartHooker = DownErrorHooker
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, ""))
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "ErrorHooker")
	downloadPartHooker = defaultDownloadPartHook

	err = s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Checkpoint(true, ""))
	c.Assert(err, IsNil)

	eq, err := compareFiles(fileName, newFile)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)
}

// TestDownloadNegative tests downloading negative
func (s *OssDownloadSuite) TestDownloadNegative(c *C) {
	objectName := objectNamePrefix + "tdn"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file-3.jpg"

	// Upload a file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	// Worker routine error
	downloadPartHooker = DownErrorHooker
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(2))
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "ErrorHooker")
	downloadPartHooker = defaultDownloadPartHook

	// Local file does not exist
	err = s.bucket.DownloadFile(objectName, "/tmp/", 100*1024, Routines(2))
	c.Assert(err, NotNil)

	// Invalid part size
	err = s.bucket.DownloadFile(objectName, newFile, 0, Routines(2))
	c.Assert(err, NotNil)

	err = s.bucket.DownloadFile(objectName, newFile, 1024*1024*1024*100, Routines(2))
	c.Assert(err, IsNil)

	err = s.bucket.DeleteObject(objectName)
	c.Assert(err, IsNil)

	// Local file does not exist
	err = s.bucket.DownloadFile(objectName, "/tmp/", 100*1024, Checkpoint(true, ""))
	c.Assert(err, NotNil)

	err = s.bucket.DownloadFile(objectName, "/tmp/", 100*1024, Routines(2), Checkpoint(true, ""))
	c.Assert(err, NotNil)

	// Invalid part size
	err = s.bucket.DownloadFile(objectName, newFile, -1, Checkpoint(true, ""))
	c.Assert(err, NotNil)

	err = s.bucket.DownloadFile(objectName, newFile, 0, Routines(2), Checkpoint(true, ""))
	c.Assert(err, NotNil)

	err = s.bucket.DownloadFile(objectName, newFile, 1024*1024*1024*100, Checkpoint(true, ""))
	c.Assert(err, NotNil)

	err = s.bucket.DownloadFile(objectName, newFile, 1024*1024*1024*100, Routines(2), Checkpoint(true, ""))
	c.Assert(err, NotNil)
}

// TestDownloadWithRange tests concurrent downloading with range specified and checkpoint enabled
func (s *OssDownloadSuite) TestDownloadWithRange(c *C) {
	objectName := objectNamePrefix + "tdwr"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file-tdwr.jpg"
	newFileGet := "down-new-file-tdwr-2.jpg"

	// Upload a file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	fileSize, err := getFileSize(fileName)
	c.Assert(err, IsNil)

	// Download with range, from 1024 to 4096
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(3), Range(1024, 4095))
	c.Assert(err, IsNil)

	// Check
	eq, err := compareFilesWithRange(fileName, 1024, newFile, 0, 3072)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, Range(1024, 4095))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Download with range, from 1024 to 4096
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 1024, Routines(3), NormalizedRange("1024-4095"))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFilesWithRange(fileName, 1024, newFile, 0, 3072)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, NormalizedRange("1024-4095"))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Download with range, from 2048 to the end
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 1024*1024, Routines(3), NormalizedRange("2048-"))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFilesWithRange(fileName, 2048, newFile, 0, fileSize-2048)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, NormalizedRange("2048-"))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Download with range, the last 4096
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 1024, Routines(3), NormalizedRange("-4096"))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFilesWithRange(fileName, fileSize-4096, newFile, 0, 4096)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, NormalizedRange("-4096"))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	err = s.bucket.DeleteObject(objectName)
	c.Assert(err, IsNil)
}

// TestDownloadWithCheckoutAndRange tests concurrent downloading with range specified and checkpoint enabled
func (s *OssDownloadSuite) TestDownloadWithCheckoutAndRange(c *C) {
	objectName := objectNamePrefix + "tdwcr"
	fileName := "../sample/BingWallpaper-2015-11-07.jpg"
	newFile := "down-new-file-tdwcr.jpg"
	newFileGet := "down-new-file-tdwcr-2.jpg"

	// Upload a file
	err := s.bucket.UploadFile(objectName, fileName, 100*1024, Routines(3))
	c.Assert(err, IsNil)

	fileSize, err := getFileSize(fileName)
	c.Assert(err, IsNil)

	// Download with range, from 1024 to 4096
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 100*1024, Routines(3), Checkpoint(true, ""), Range(1024, 4095))
	c.Assert(err, IsNil)

	// Check
	eq, err := compareFilesWithRange(fileName, 1024, newFile, 0, 3072)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, Range(1024, 4095))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Download with range, from 1024 to 4096
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 1024, Routines(3), Checkpoint(true, ""), NormalizedRange("1024-4095"))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFilesWithRange(fileName, 1024, newFile, 0, 3072)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, NormalizedRange("1024-4095"))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Download with range, from 2048 to the end
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 1024*1024, Routines(3), Checkpoint(true, ""), NormalizedRange("2048-"))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFilesWithRange(fileName, 2048, newFile, 0, fileSize-2048)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, NormalizedRange("2048-"))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	// Download with range, the last 4096 bytes
	os.Remove(newFile)
	err = s.bucket.DownloadFile(objectName, newFile, 1024, Routines(3), Checkpoint(true, ""), NormalizedRange("-4096"))
	c.Assert(err, IsNil)

	// Check
	eq, err = compareFilesWithRange(fileName, fileSize-4096, newFile, 0, 4096)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	os.Remove(newFileGet)
	err = s.bucket.GetObjectToFile(objectName, newFileGet, NormalizedRange("-4096"))
	c.Assert(err, IsNil)

	// Compare get and download
	eq, err = compareFiles(newFile, newFileGet)
	c.Assert(err, IsNil)
	c.Assert(eq, Equals, true)

	err = s.bucket.DeleteObject(objectName)
	c.Assert(err, IsNil)
}

// TestCombineCRCInDownloadParts tests combineCRCInParts
func (s *OssDownloadSuite) TestCombineCRCInDownloadParts(c *C) {
	crc := combineCRCInParts(nil)
	c.Assert(crc == 0, Equals, true)

	crc = combineCRCInParts(make([]downloadPart, 0))
	c.Assert(crc == 0, Equals, true)

	parts := make([]downloadPart, 1)
	parts[0].CRC64 = 10278880121275185425
	crc = combineCRCInParts(parts)
	c.Assert(crc == 10278880121275185425, Equals, true)

	parts = make([]downloadPart, 2)
	parts[0].CRC64 = 6748440630437108969
	parts[0].Start = 0
	parts[0].End = 4
	parts[1].CRC64 = 10278880121275185425
	parts[1].Start = 5
	parts[1].End = 8
	crc = combineCRCInParts(parts)
	c.Assert(crc == 11051210869376104954, Equals, true)
}

func getFileSize(fileName string) (int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

// compareFilesWithRange compares the content between fileL and fileR with specified range
func compareFilesWithRange(fileL string, offsetL int64, fileR string, offsetR int64, size int64) (bool, error) {
	finL, err := os.Open(fileL)
	if err != nil {
		return false, err
	}
	defer finL.Close()
	finL.Seek(offsetL, os.SEEK_SET)

	finR, err := os.Open(fileR)
	if err != nil {
		return false, err
	}
	defer finR.Close()
	finR.Seek(offsetR, os.SEEK_SET)

	statL, err := finL.Stat()
	if err != nil {
		return false, err
	}

	statR, err := finR.Stat()
	if err != nil {
		return false, err
	}

	if (offsetL+size > statL.Size()) || (offsetR+size > statR.Size()) {
		return false, nil
	}

	part := statL.Size() - offsetL
	if part > 16*1024 {
		part = 16 * 1024
	}

	bufL := make([]byte, part)
	bufR := make([]byte, part)
	for readN := int64(0); readN < size; {
		n, _ := finL.Read(bufL)
		if 0 == n {
			break
		}

		n, _ = finR.Read(bufR)
		if 0 == n {
			break
		}

		tailer := part
		if tailer > size-readN {
			tailer = size - readN
		}
		readN += tailer

		if !bytes.Equal(bufL[0:tailer], bufR[0:tailer]) {
			return false, nil
		}
	}

	return true, nil
}
