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
package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

const dataStructuresLogNameTemplate = "daemon-data-%s.log"

// dumpDaemon appends the daemon datastructures into file in dir and returns full path
// to that file.
func (d *Daemon) dumpDaemon(dir string) (string, error) {
	// Ensure we recover from a panic as we are doing this without any locking
	defer func() {
		recover()
	}()

	path := filepath.Join(dir, fmt.Sprintf(dataStructuresLogNameTemplate, strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1)))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file to write the daemon datastructure dump")
	}
	defer f.Close()

	dump := struct {
		containers      interface{}
		names           interface{}
		links           interface{}
		execs           interface{}
		volumes         interface{}
		images          interface{}
		layers          interface{}
		imageReferences interface{}
		downloads       interface{}
		uploads         interface{}
		registry        interface{}
		plugins         interface{}
	}{
		containers:      d.containers,
		execs:           d.execCommands,
		volumes:         d.volumes,
		images:          d.imageStore,
		layers:          d.layerStore,
		imageReferences: d.referenceStore,
		downloads:       d.downloadManager,
		uploads:         d.uploadManager,
		registry:        d.RegistryService,
		plugins:         d.PluginStore,
		names:           d.nameIndex,
		links:           d.linkIndex,
	}

	spew.Fdump(f, dump) // Does not return an error
	f.Sync()
	return path, nil
}
