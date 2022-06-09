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
// +build !windows

package chrootarchive

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	reexec.Register("docker-applyLayer", applyLayer)
	reexec.Register("docker-untar", untar)
}

func fatal(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

// flush consumes all the bytes from the reader discarding
// any errors
func flush(r io.Reader) (bytes int64, err error) {
	return io.Copy(ioutil.Discard, r)
}
