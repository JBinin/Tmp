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
package nl

import "encoding/binary"

const (
	MPLS_LS_LABEL_SHIFT = 12
	MPLS_LS_S_SHIFT     = 8
)

func EncodeMPLSStack(labels ...int) []byte {
	b := make([]byte, 4*len(labels))
	for idx, label := range labels {
		l := label << MPLS_LS_LABEL_SHIFT
		if idx == len(labels)-1 {
			l |= 1 << MPLS_LS_S_SHIFT
		}
		binary.BigEndian.PutUint32(b[idx*4:], uint32(l))
	}
	return b
}

func DecodeMPLSStack(buf []byte) []int {
	if len(buf)%4 != 0 {
		return nil
	}
	stack := make([]int, 0, len(buf)/4)
	for len(buf) > 0 {
		l := binary.BigEndian.Uint32(buf[:4])
		buf = buf[4:]
		stack = append(stack, int(l)>>MPLS_LS_LABEL_SHIFT)
		if (l>>MPLS_LS_S_SHIFT)&1 > 0 {
			break
		}
	}
	return stack
}
