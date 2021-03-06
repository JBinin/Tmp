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
package rice

// LocateMethod defines how a box is located.
type LocateMethod int

const (
	LocateFS               = LocateMethod(iota) // Locate on the filesystem according to package path.
	LocateAppended                              // Locate boxes appended to the executable.
	LocateEmbedded                              // Locate embedded boxes.
	LocateWorkingDirectory                      // Locate on the binary working directory
)

// Config allows customizing the box lookup behavior.
type Config struct {
	// LocateOrder defines the priority order that boxes are searched for. By
	// default, the package global FindBox searches for embedded boxes first,
	// then appended boxes, and then finally boxes on the filesystem.  That
	// search order may be customized by provided the ordered list here. Leaving
	// out a particular method will omit that from the search space. For
	// example, []LocateMethod{LocateEmbedded, LocateAppended} will never search
	// the filesystem for boxes.
	LocateOrder []LocateMethod
}

// FindBox searches for boxes using the LocateOrder of the config.
func (c *Config) FindBox(boxName string) (*Box, error) {
	return findBox(boxName, c.LocateOrder)
}

// MustFindBox searches for boxes using the LocateOrder of the config, like
// FindBox does.  It does not return an error, instead it panics when an error
// occurs.
func (c *Config) MustFindBox(boxName string) *Box {
	box, err := findBox(boxName, c.LocateOrder)
	if err != nil {
		panic(err)
	}
	return box
}
