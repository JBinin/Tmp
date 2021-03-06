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
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/pkg/archive"
	"github.com/go-check/check"
)

type fileType uint32

const (
	ftRegular fileType = iota
	ftDir
	ftSymlink
)

type fileData struct {
	filetype fileType
	path     string
	contents string
}

func (fd fileData) creationCommand() string {
	var command string

	switch fd.filetype {
	case ftRegular:
		// Don't overwrite the file if it already exists!
		command = fmt.Sprintf("if [ ! -f %s ]; then echo %q > %s; fi", fd.path, fd.contents, fd.path)
	case ftDir:
		command = fmt.Sprintf("mkdir -p %s", fd.path)
	case ftSymlink:
		command = fmt.Sprintf("ln -fs %s %s", fd.contents, fd.path)
	}

	return command
}

func mkFilesCommand(fds []fileData) string {
	commands := make([]string, len(fds))

	for i, fd := range fds {
		commands[i] = fd.creationCommand()
	}

	return strings.Join(commands, " && ")
}

var defaultFileData = []fileData{
	{ftRegular, "file1", "file1"},
	{ftRegular, "file2", "file2"},
	{ftRegular, "file3", "file3"},
	{ftRegular, "file4", "file4"},
	{ftRegular, "file5", "file5"},
	{ftRegular, "file6", "file6"},
	{ftRegular, "file7", "file7"},
	{ftDir, "dir1", ""},
	{ftRegular, "dir1/file1-1", "file1-1"},
	{ftRegular, "dir1/file1-2", "file1-2"},
	{ftDir, "dir2", ""},
	{ftRegular, "dir2/file2-1", "file2-1"},
	{ftRegular, "dir2/file2-2", "file2-2"},
	{ftDir, "dir3", ""},
	{ftRegular, "dir3/file3-1", "file3-1"},
	{ftRegular, "dir3/file3-2", "file3-2"},
	{ftDir, "dir4", ""},
	{ftRegular, "dir4/file3-1", "file4-1"},
	{ftRegular, "dir4/file3-2", "file4-2"},
	{ftDir, "dir5", ""},
	{ftSymlink, "symlinkToFile1", "file1"},
	{ftSymlink, "symlinkToDir1", "dir1"},
	{ftSymlink, "brokenSymlinkToFileX", "fileX"},
	{ftSymlink, "brokenSymlinkToDirX", "dirX"},
	{ftSymlink, "symlinkToAbsDir", "/root"},
}

func defaultMkContentCommand() string {
	return mkFilesCommand(defaultFileData)
}

func makeTestContentInDir(c *check.C, dir string) {
	for _, fd := range defaultFileData {
		path := filepath.Join(dir, filepath.FromSlash(fd.path))
		switch fd.filetype {
		case ftRegular:
			c.Assert(ioutil.WriteFile(path, []byte(fd.contents+"\n"), os.FileMode(0666)), checker.IsNil)
		case ftDir:
			c.Assert(os.Mkdir(path, os.FileMode(0777)), checker.IsNil)
		case ftSymlink:
			c.Assert(os.Symlink(fd.contents, path), checker.IsNil)
		}
	}
}

type testContainerOptions struct {
	addContent bool
	readOnly   bool
	volumes    []string
	workDir    string
	command    string
}

func makeTestContainer(c *check.C, options testContainerOptions) (containerID string) {
	if options.addContent {
		mkContentCmd := defaultMkContentCommand()
		if options.command == "" {
			options.command = mkContentCmd
		} else {
			options.command = fmt.Sprintf("%s && %s", defaultMkContentCommand(), options.command)
		}
	}

	if options.command == "" {
		options.command = "#(nop)"
	}

	args := []string{"run", "-d"}

	for _, volume := range options.volumes {
		args = append(args, "-v", volume)
	}

	if options.workDir != "" {
		args = append(args, "-w", options.workDir)
	}

	if options.readOnly {
		args = append(args, "--read-only")
	}

	args = append(args, "busybox", "/bin/sh", "-c", options.command)

	out, _ := dockerCmd(c, args...)

	containerID = strings.TrimSpace(out)

	out, _ = dockerCmd(c, "wait", containerID)

	exitCode := strings.TrimSpace(out)
	if exitCode != "0" {
		out, _ = dockerCmd(c, "logs", containerID)
	}
	c.Assert(exitCode, checker.Equals, "0", check.Commentf("failed to make test container: %s", out))

	return
}

func makeCatFileCommand(path string) string {
	return fmt.Sprintf("if [ -f %s ]; then cat %s; fi", path, path)
}

func cpPath(pathElements ...string) string {
	localizedPathElements := make([]string, len(pathElements))
	for i, path := range pathElements {
		localizedPathElements[i] = filepath.FromSlash(path)
	}
	return strings.Join(localizedPathElements, string(filepath.Separator))
}

func cpPathTrailingSep(pathElements ...string) string {
	return fmt.Sprintf("%s%c", cpPath(pathElements...), filepath.Separator)
}

func containerCpPath(containerID string, pathElements ...string) string {
	joined := strings.Join(pathElements, "/")
	return fmt.Sprintf("%s:%s", containerID, joined)
}

func containerCpPathTrailingSep(containerID string, pathElements ...string) string {
	return fmt.Sprintf("%s/", containerCpPath(containerID, pathElements...))
}

func runDockerCp(c *check.C, src, dst string) (err error) {
	c.Logf("running `docker cp %s %s`", src, dst)

	args := []string{"cp", src, dst}

	out, _, err := runCommandWithOutput(exec.Command(dockerBinary, args...))
	if err != nil {
		err = fmt.Errorf("error executing `docker cp` command: %s: %s", err, out)
	}

	return
}

func startContainerGetOutput(c *check.C, containerID string) (out string, err error) {
	c.Logf("running `docker start -a %s`", containerID)

	args := []string{"start", "-a", containerID}

	out, _, err = runCommandWithOutput(exec.Command(dockerBinary, args...))
	if err != nil {
		err = fmt.Errorf("error executing `docker start` command: %s: %s", err, out)
	}

	return
}

func getTestDir(c *check.C, label string) (tmpDir string) {
	var err error

	tmpDir, err = ioutil.TempDir("", label)
	// unable to make temporary directory
	c.Assert(err, checker.IsNil)

	return
}

func isCpNotExist(err error) bool {
	return strings.Contains(err.Error(), "no such file or directory") || strings.Contains(err.Error(), "cannot find the file specified")
}

func isCpDirNotExist(err error) bool {
	return strings.Contains(err.Error(), archive.ErrDirNotExists.Error())
}

func isCpNotDir(err error) bool {
	return strings.Contains(err.Error(), archive.ErrNotDirectory.Error()) || strings.Contains(err.Error(), "filename, directory name, or volume label syntax is incorrect")
}

func isCpCannotCopyDir(err error) bool {
	return strings.Contains(err.Error(), archive.ErrCannotCopyDir.Error())
}

func isCpCannotCopyReadOnly(err error) bool {
	return strings.Contains(err.Error(), "marked read-only")
}

func isCannotOverwriteNonDirWithDir(err error) bool {
	return strings.Contains(err.Error(), "cannot overwrite non-directory")
}

func fileContentEquals(c *check.C, filename, contents string) (err error) {
	c.Logf("checking that file %q contains %q\n", filename, contents)

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	expectedBytes, err := ioutil.ReadAll(strings.NewReader(contents))
	if err != nil {
		return
	}

	if !bytes.Equal(fileBytes, expectedBytes) {
		err = fmt.Errorf("file content not equal - expected %q, got %q", string(expectedBytes), string(fileBytes))
	}

	return
}

func symlinkTargetEquals(c *check.C, symlink, expectedTarget string) (err error) {
	c.Logf("checking that the symlink %q points to %q\n", symlink, expectedTarget)

	actualTarget, err := os.Readlink(symlink)
	if err != nil {
		return
	}

	if actualTarget != expectedTarget {
		err = fmt.Errorf("symlink target points to %q not %q", actualTarget, expectedTarget)
	}

	return
}

func containerStartOutputEquals(c *check.C, containerID, contents string) (err error) {
	c.Logf("checking that container %q start output contains %q\n", containerID, contents)

	out, err := startContainerGetOutput(c, containerID)
	if err != nil {
		return
	}

	if out != contents {
		err = fmt.Errorf("output contents not equal - expected %q, got %q", contents, out)
	}

	return
}

func defaultVolumes(tmpDir string) []string {
	if SameHostDaemon() {
		return []string{
			"/vol1",
			fmt.Sprintf("%s:/vol2", tmpDir),
			fmt.Sprintf("%s:/vol3", filepath.Join(tmpDir, "vol3")),
			fmt.Sprintf("%s:/vol_ro:ro", filepath.Join(tmpDir, "vol_ro")),
		}
	}

	// Can't bind-mount volumes with separate host daemon.
	return []string{"/vol1", "/vol2", "/vol3", "/vol_ro:/vol_ro:ro"}
}
