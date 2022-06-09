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
package dhcp4

type OptionCode byte

type Option struct {
	Code  OptionCode
	Value []byte
}

// Map of DHCP options
type Options map[OptionCode][]byte

// SelectOrderOrAll has same functionality as SelectOrder, except if the order
// param is nil, whereby all options are added (in arbitrary order).
func (o Options) SelectOrderOrAll(order []byte) []Option {
	if order == nil {
		opts := make([]Option, 0, len(o))
		for i, v := range o {
			opts = append(opts, Option{Code: i, Value: v})
		}
		return opts
	}
	return o.SelectOrder(order)
}

// SelectOrder returns a slice of options ordered and selected by a byte array
// usually defined by OptionParameterRequestList.  This result is expected to be
// used in ReplyPacket()'s []Option parameter.
func (o Options) SelectOrder(order []byte) []Option {
	opts := make([]Option, 0, len(order))
	for _, v := range order {
		if data, ok := o[OptionCode(v)]; ok {
			opts = append(opts, Option{Code: OptionCode(v), Value: data})
		}
	}
	return opts
}

type OpCode byte
type MessageType byte // Option 53
