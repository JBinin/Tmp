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
package bundlefile

import (
	"encoding/json"
	"fmt"
	"io"
)

// Bundlefile stores the contents of a bundlefile
type Bundlefile struct {
	Version  string
	Services map[string]Service
}

// Service is a service from a bundlefile
type Service struct {
	Image      string
	Command    []string          `json:",omitempty"`
	Args       []string          `json:",omitempty"`
	Env        []string          `json:",omitempty"`
	Labels     map[string]string `json:",omitempty"`
	Ports      []Port            `json:",omitempty"`
	WorkingDir *string           `json:",omitempty"`
	User       *string           `json:",omitempty"`
	Networks   []string          `json:",omitempty"`
}

// Port is a port as defined in a bundlefile
type Port struct {
	Protocol string
	Port     uint32
}

// LoadFile loads a bundlefile from a path to the file
func LoadFile(reader io.Reader) (*Bundlefile, error) {
	bundlefile := &Bundlefile{}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(bundlefile); err != nil {
		switch jsonErr := err.(type) {
		case *json.SyntaxError:
			return nil, fmt.Errorf(
				"JSON syntax error at byte %v: %s",
				jsonErr.Offset,
				jsonErr.Error())
		case *json.UnmarshalTypeError:
			return nil, fmt.Errorf(
				"Unexpected type at byte %v. Expected %s but received %s.",
				jsonErr.Offset,
				jsonErr.Type,
				jsonErr.Value)
		}
		return nil, err
	}

	return bundlefile, nil
}

// Print writes the contents of the bundlefile to the output writer
// as human readable json
func Print(out io.Writer, bundle *Bundlefile) error {
	bytes, err := json.MarshalIndent(*bundle, "", "    ")
	if err != nil {
		return err
	}

	_, err = out.Write(bytes)
	return err
}
