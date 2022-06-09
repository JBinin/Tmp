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
package testsuite

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TestSuite struct {
	Path        string
	PackageName string
	IsGinkgo    bool
	Precompiled bool
}

func PrecompiledTestSuite(path string) (TestSuite, error) {
	info, err := os.Stat(path)
	if err != nil {
		return TestSuite{}, err
	}

	if info.IsDir() {
		return TestSuite{}, errors.New("this is a directory, not a file")
	}

	if filepath.Ext(path) != ".test" {
		return TestSuite{}, errors.New("this is not a .test binary")
	}

	if info.Mode()&0111 == 0 {
		return TestSuite{}, errors.New("this is not executable")
	}

	dir := relPath(filepath.Dir(path))
	packageName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return TestSuite{
		Path:        dir,
		PackageName: packageName,
		IsGinkgo:    true,
		Precompiled: true,
	}, nil
}

func SuitesInDir(dir string, recurse bool) []TestSuite {
	suites := []TestSuite{}

	if vendorExperimentCheck(dir) {
		return suites
	}

	files, _ := ioutil.ReadDir(dir)
	re := regexp.MustCompile(`_test\.go$`)
	for _, file := range files {
		if !file.IsDir() && re.Match([]byte(file.Name())) {
			suites = append(suites, New(dir, files))
			break
		}
	}

	if recurse {
		re = regexp.MustCompile(`^[._]`)
		for _, file := range files {
			if file.IsDir() && !re.Match([]byte(file.Name())) {
				suites = append(suites, SuitesInDir(dir+"/"+file.Name(), recurse)...)
			}
		}
	}

	return suites
}

func relPath(dir string) string {
	dir, _ = filepath.Abs(dir)
	cwd, _ := os.Getwd()
	dir, _ = filepath.Rel(cwd, filepath.Clean(dir))
	dir = "." + string(filepath.Separator) + dir
	return dir
}

func New(dir string, files []os.FileInfo) TestSuite {
	return TestSuite{
		Path:        relPath(dir),
		PackageName: packageNameForSuite(dir),
		IsGinkgo:    filesHaveGinkgoSuite(dir, files),
	}
}

func packageNameForSuite(dir string) string {
	path, _ := filepath.Abs(dir)
	return filepath.Base(path)
}

func filesHaveGinkgoSuite(dir string, files []os.FileInfo) bool {
	reTestFile := regexp.MustCompile(`_test\.go$`)
	reGinkgo := regexp.MustCompile(`package ginkgo|\/ginkgo"`)

	for _, file := range files {
		if !file.IsDir() && reTestFile.Match([]byte(file.Name())) {
			contents, _ := ioutil.ReadFile(dir + "/" + file.Name())
			if reGinkgo.Match(contents) {
				return true
			}
		}
	}

	return false
}
