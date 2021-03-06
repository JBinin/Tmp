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
// +build !selinux !linux

package label

// InitLabels returns the process label and file labels to be used within
// the container.  A list of options can be passed into this function to alter
// the labels.
func InitLabels(options []string) (string, string, error) {
	return "", "", nil
}

func GetROMountLabel() string {
	return ""
}

func GenLabels(options string) (string, string, error) {
	return "", "", nil
}

func FormatMountLabel(src string, mountLabel string) string {
	return src
}

func SetProcessLabel(processLabel string) error {
	return nil
}

func GetFileLabel(path string) (string, error) {
	return "", nil
}

func SetFileLabel(path string, fileLabel string) error {
	return nil
}

func SetFileCreateLabel(fileLabel string) error {
	return nil
}

func Relabel(path string, fileLabel string, shared bool) error {
	return nil
}

func GetPidLabel(pid int) (string, error) {
	return "", nil
}

func Init() {
}

func ReserveLabel(label string) error {
	return nil
}

func ReleaseLabel(label string) error {
	return nil
}

// DupSecOpt takes a process label and returns security options that
// can be used to set duplicate labels on future container processes
func DupSecOpt(src string) []string {
	return nil
}

// DisableSecOpt returns a security opt that can disable labeling
// support for future container processes
func DisableSecOpt() []string {
	return nil
}

// Validate checks that the label does not include unexpected options
func Validate(label string) error {
	return nil
}

// RelabelNeeded checks whether the user requested a relabel
func RelabelNeeded(label string) bool {
	return false
}

// IsShared checks that the label includes a "shared" mark
func IsShared(label string) bool {
	return false
}
