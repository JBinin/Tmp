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
package libcontainerd

import (
	"io"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/docker/pkg/ioutils"
)

// process keeps the state for both main container process and exec process.
type process struct {
	processCommon

	// Platform specific fields are below here.
	hcsProcess hcsshim.Process
}

type autoClosingReader struct {
	io.ReadCloser
	sync.Once
}

func (r *autoClosingReader) Read(b []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(b)
	if err == io.EOF {
		r.Once.Do(func() { r.ReadCloser.Close() })
	}
	return
}

func createStdInCloser(pipe io.WriteCloser, process hcsshim.Process) io.WriteCloser {
	return ioutils.NewWriteCloserWrapper(pipe, func() error {
		if err := pipe.Close(); err != nil {
			return err
		}

		err := process.CloseStdin()
		if err != nil && !hcsshim.IsNotExist(err) && !hcsshim.IsAlreadyClosed(err) {
			// This error will occur if the compute system is currently shutting down
			if perr, ok := err.(*hcsshim.ProcessError); ok && perr.Err != hcsshim.ErrVmcomputeOperationInvalidState {
				return err
			}
		}

		return nil
	})
}
