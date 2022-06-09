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
package watch

import (
	"path/filepath"
	"sync"
)

type PackageHashes struct {
	PackageHashes map[string]*PackageHash
	usedPaths     map[string]bool
	lock          *sync.Mutex
}

func NewPackageHashes() *PackageHashes {
	return &PackageHashes{
		PackageHashes: map[string]*PackageHash{},
		usedPaths:     nil,
		lock:          &sync.Mutex{},
	}
}

func (p *PackageHashes) CheckForChanges() []string {
	p.lock.Lock()
	defer p.lock.Unlock()

	modified := []string{}

	for _, packageHash := range p.PackageHashes {
		if packageHash.CheckForChanges() {
			modified = append(modified, packageHash.path)
		}
	}

	return modified
}

func (p *PackageHashes) Add(path string) *PackageHash {
	p.lock.Lock()
	defer p.lock.Unlock()

	path, _ = filepath.Abs(path)
	_, ok := p.PackageHashes[path]
	if !ok {
		p.PackageHashes[path] = NewPackageHash(path)
	}

	if p.usedPaths != nil {
		p.usedPaths[path] = true
	}
	return p.PackageHashes[path]
}

func (p *PackageHashes) Get(path string) *PackageHash {
	p.lock.Lock()
	defer p.lock.Unlock()

	path, _ = filepath.Abs(path)
	if p.usedPaths != nil {
		p.usedPaths[path] = true
	}
	return p.PackageHashes[path]
}

func (p *PackageHashes) StartTrackingUsage() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.usedPaths = map[string]bool{}
}

func (p *PackageHashes) StopTrackingUsageAndPrune() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for path := range p.PackageHashes {
		if !p.usedPaths[path] {
			delete(p.PackageHashes, path)
		}
	}

	p.usedPaths = nil
}
