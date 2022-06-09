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
package container

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestCalculateBlockIO(t *testing.T) {
	blkio := types.BlkioStats{
		IoServiceBytesRecursive: []types.BlkioStatEntry{{8, 0, "read", 1234}, {8, 1, "read", 4567}, {8, 0, "write", 123}, {8, 1, "write", 456}},
	}
	blkRead, blkWrite := calculateBlockIO(blkio)
	if blkRead != 5801 {
		t.Fatalf("blkRead = %d, want 5801", blkRead)
	}
	if blkWrite != 579 {
		t.Fatalf("blkWrite = %d, want 579", blkWrite)
	}
}
