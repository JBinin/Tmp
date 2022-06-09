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
package storage

import (
	"errors"
	"fmt"
)

var (
	// ErrPathOutsideStore indicates that the returned path would be
	// outside the store
	ErrPathOutsideStore = errors.New("path outside file store")
)

// ErrMetaNotFound indicates we did not find a particular piece
// of metadata in the store
type ErrMetaNotFound struct {
	Resource string
}

func (err ErrMetaNotFound) Error() string {
	return fmt.Sprintf("%s trust data unavailable.  Has a notary repository been initialized?", err.Resource)
}
