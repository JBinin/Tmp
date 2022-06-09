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
package archive

import (
	"os"

	"github.com/docker/docker/pkg/system"
)

func statDifferent(oldStat *system.StatT, newStat *system.StatT) bool {

	// Don't look at size for dirs, its not a good measure of change
	if oldStat.ModTime() != newStat.ModTime() ||
		oldStat.Mode() != newStat.Mode() ||
		oldStat.Size() != newStat.Size() && !oldStat.IsDir() {
		return true
	}
	return false
}

func (info *FileInfo) isDir() bool {
	return info.parent == nil || info.stat.IsDir()
}

func getIno(fi os.FileInfo) (inode uint64) {
	return
}

func hasHardlinks(fi os.FileInfo) bool {
	return false
}
