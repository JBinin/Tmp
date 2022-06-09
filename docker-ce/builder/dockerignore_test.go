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
package builder

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"testing"
)

const shouldStayFilename = "should_stay"

func extractFilenames(files []os.FileInfo) []string {
	filenames := make([]string, len(files), len(files))

	for i, file := range files {
		filenames[i] = file.Name()
	}

	return filenames
}

func checkDirectory(t *testing.T, dir string, expectedFiles []string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		t.Fatalf("Could not read directory: %s", err)
	}

	if len(files) != len(expectedFiles) {
		log.Fatalf("Directory should contain exactly %d file(s), got %d", len(expectedFiles), len(files))
	}

	filenames := extractFilenames(files)
	sort.Strings(filenames)
	sort.Strings(expectedFiles)

	for i, filename := range filenames {
		if filename != expectedFiles[i] {
			t.Fatalf("File %s should be in the directory, got: %s", expectedFiles[i], filename)
		}
	}
}

func executeProcess(t *testing.T, contextDir string) {
	modifiableCtx := &tarSumContext{root: contextDir}
	ctx := DockerIgnoreContext{ModifiableContext: modifiableCtx}

	err := ctx.Process([]string{DefaultDockerfileName})

	if err != nil {
		t.Fatalf("Error when executing Process: %s", err)
	}
}

func TestProcessShouldRemoveDockerfileDockerignore(t *testing.T) {
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerignore-process-test")
	defer cleanup()

	createTestTempFile(t, contextDir, shouldStayFilename, testfileContents, 0777)
	createTestTempFile(t, contextDir, dockerignoreFilename, "Dockerfile\n.dockerignore", 0777)
	createTestTempFile(t, contextDir, DefaultDockerfileName, dockerfileContents, 0777)

	executeProcess(t, contextDir)

	checkDirectory(t, contextDir, []string{shouldStayFilename})

}

func TestProcessNoDockerignore(t *testing.T) {
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerignore-process-test")
	defer cleanup()

	createTestTempFile(t, contextDir, shouldStayFilename, testfileContents, 0777)
	createTestTempFile(t, contextDir, DefaultDockerfileName, dockerfileContents, 0777)

	executeProcess(t, contextDir)

	checkDirectory(t, contextDir, []string{shouldStayFilename, DefaultDockerfileName})

}

func TestProcessShouldLeaveAllFiles(t *testing.T) {
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerignore-process-test")
	defer cleanup()

	createTestTempFile(t, contextDir, shouldStayFilename, testfileContents, 0777)
	createTestTempFile(t, contextDir, DefaultDockerfileName, dockerfileContents, 0777)
	createTestTempFile(t, contextDir, dockerignoreFilename, "input1\ninput2", 0777)

	executeProcess(t, contextDir)

	checkDirectory(t, contextDir, []string{shouldStayFilename, DefaultDockerfileName, dockerignoreFilename})

}
