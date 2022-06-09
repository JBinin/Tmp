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
package netlink

import (
	"errors"
)

var (
	// ErrAttrHeaderTruncated is returned when a netlink attribute's header is
	// truncated.
	ErrAttrHeaderTruncated = errors.New("attribute header truncated")
	// ErrAttrBodyTruncated is returned when a netlink attribute's body is
	// truncated.
	ErrAttrBodyTruncated = errors.New("attribute body truncated")
)

type Fou struct {
	Family    int
	Port      int
	Protocol  int
	EncapType int
}
