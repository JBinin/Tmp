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
package govalidator

// Errors is an array of multiple errors and conforms to the error interface.
type Errors []error

// Errors returns itself.
func (es Errors) Errors() []error {
	return es
}

func (es Errors) Error() string {
	var err string
	for _, e := range es {
		err += e.Error() + ";"
	}
	return err
}

// Error encapsulates a name, an error and whether there's a custom error message or not.
type Error struct {
	Name                     string
	Err                      error
	CustomErrorMessageExists bool
}

func (e Error) Error() string {
	if e.CustomErrorMessageExists {
		return e.Err.Error()
	}
	return e.Name + ": " + e.Err.Error()
}
