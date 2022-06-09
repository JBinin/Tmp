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
package matchers

import (
	"fmt"
	"os"

	"github.com/onsi/gomega/format"
)

type notADirectoryError struct {
	os.FileInfo
}

func (t notADirectoryError) Error() string {
	fileInfo := os.FileInfo(t)
	switch {
	case fileInfo.Mode().IsRegular():
		return "file is a regular file"
	default:
		return fmt.Sprintf("file mode is: %s", fileInfo.Mode().String())
	}
}

type BeADirectoryMatcher struct {
	expected interface{}
	err      error
}

func (matcher *BeADirectoryMatcher) Match(actual interface{}) (success bool, err error) {
	actualFilename, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeADirectoryMatcher matcher expects a file path")
	}

	fileInfo, err := os.Stat(actualFilename)
	if err != nil {
		matcher.err = err
		return false, nil
	}

	if !fileInfo.Mode().IsDir() {
		matcher.err = notADirectoryError{fileInfo}
		return false, nil
	}
	return true, nil
}

func (matcher *BeADirectoryMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to be a directory: %s", matcher.err))
}

func (matcher *BeADirectoryMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not be a directory"))
}
