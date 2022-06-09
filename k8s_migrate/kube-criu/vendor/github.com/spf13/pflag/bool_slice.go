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
package pflag

import (
	"io"
	"strconv"
	"strings"
)

// -- boolSlice Value
type boolSliceValue struct {
	value   *[]bool
	changed bool
}

func newBoolSliceValue(val []bool, p *[]bool) *boolSliceValue {
	bsv := new(boolSliceValue)
	bsv.value = p
	*bsv.value = val
	return bsv
}

// Set converts, and assigns, the comma-separated boolean argument string representation as the []bool value of this flag.
// If Set is called on a flag that already has a []bool assigned, the newly converted values will be appended.
func (s *boolSliceValue) Set(val string) error {

	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	boolStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse boolean values into slice
	out := make([]bool, 0, len(boolStrSlice))
	for _, boolStr := range boolStrSlice {
		b, err := strconv.ParseBool(strings.TrimSpace(boolStr))
		if err != nil {
			return err
		}
		out = append(out, b)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *boolSliceValue) Type() string {
	return "boolSlice"
}

// String defines a "native" format for this boolean slice flag value.
func (s *boolSliceValue) String() string {

	boolStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		boolStrSlice[i] = strconv.FormatBool(b)
	}

	out, _ := writeAsCSV(boolStrSlice)

	return "[" + out + "]"
}

func boolSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []bool{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]bool, len(ss))
	for i, t := range ss {
		var err error
		out[i], err = strconv.ParseBool(t)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// GetBoolSlice returns the []bool value of a flag with the given name.
func (f *FlagSet) GetBoolSlice(name string) ([]bool, error) {
	val, err := f.getFlagType(name, "boolSlice", boolSliceConv)
	if err != nil {
		return []bool{}, err
	}
	return val.([]bool), nil
}

// BoolSliceVar defines a boolSlice flag with specified name, default value, and usage string.
// The argument p points to a []bool variable in which to store the value of the flag.
func (f *FlagSet) BoolSliceVar(p *[]bool, name string, value []bool, usage string) {
	f.VarP(newBoolSliceValue(value, p), name, "", usage)
}

// BoolSliceVarP is like BoolSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string) {
	f.VarP(newBoolSliceValue(value, p), name, shorthand, usage)
}

// BoolSliceVar defines a []bool flag with specified name, default value, and usage string.
// The argument p points to a []bool variable in which to store the value of the flag.
func BoolSliceVar(p *[]bool, name string, value []bool, usage string) {
	CommandLine.VarP(newBoolSliceValue(value, p), name, "", usage)
}

// BoolSliceVarP is like BoolSliceVar, but accepts a shorthand letter that can be used after a single dash.
func BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string) {
	CommandLine.VarP(newBoolSliceValue(value, p), name, shorthand, usage)
}

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func (f *FlagSet) BoolSlice(name string, value []bool, usage string) *[]bool {
	p := []bool{}
	f.BoolSliceVarP(&p, name, "", value, usage)
	return &p
}

// BoolSliceP is like BoolSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool {
	p := []bool{}
	f.BoolSliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func BoolSlice(name string, value []bool, usage string) *[]bool {
	return CommandLine.BoolSliceP(name, "", value, usage)
}

// BoolSliceP is like BoolSlice, but accepts a shorthand letter that can be used after a single dash.
func BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool {
	return CommandLine.BoolSliceP(name, shorthand, value, usage)
}
