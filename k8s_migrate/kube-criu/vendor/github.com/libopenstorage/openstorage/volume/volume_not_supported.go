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
package volume

import (
	"github.com/libopenstorage/openstorage/api"
)

var (
	// BlockNotSupported is a default (null) block driver implementation.  This can be
	// used by drivers that do not want to (or care about) implementing the attach,
	// format and detach interfaces.
	BlockNotSupported = &blockNotSupported{}
	// SnapshotNotSupported is a null snapshot driver implementation. This can be used
	// by drivers that do not want to implement the snapshot interface
	SnapshotNotSupported = &snapshotNotSupported{}
	// IONotSupported is a null IODriver interface
	IONotSupported = &ioNotSupported{}
	// StatsNotSupported is a null stats driver implementation. This can be used
	// by drivers that do not want to implement the stats interface.
	StatsNotSupported = &statsNotSupported{}
)

type blockNotSupported struct{}

func (b *blockNotSupported) Attach(volumeID string, attachOptions map[string]string) (string, error) {
	return "", ErrNotSupported
}

func (b *blockNotSupported) Detach(volumeID string, unmountBeforeDetach bool) error {
	return ErrNotSupported
}

type snapshotNotSupported struct{}

func (s *snapshotNotSupported) Snapshot(volumeID string, readonly bool, locator *api.VolumeLocator) (string, error) {
	return "", ErrNotSupported
}

func (s *snapshotNotSupported) Restore(volumeID, snapshotID string) error {
	return ErrNotSupported
}

type ioNotSupported struct{}

func (i *ioNotSupported) Read(volumeID string, buffer []byte, size uint64, offset int64) (int64, error) {
	return 0, ErrNotSupported
}

func (i *ioNotSupported) Write(volumeID string, buffer []byte, size uint64, offset int64) (int64, error) {
	return 0, ErrNotSupported
}

func (i *ioNotSupported) Flush(volumeID string) error {
	return ErrNotSupported
}

type statsNotSupported struct{}

// Stats returns stats
func (s *statsNotSupported) Stats(
	volumeID string,
	cumulative bool,
) (*api.Stats, error) {
	return nil, ErrNotSupported
}

// UsedSize returns allocated size
func (s *statsNotSupported) UsedSize(volumeID string) (uint64, error) {
	return 0, ErrNotSupported
}

// GetActiveRequests gets active requests
func (s *statsNotSupported) GetActiveRequests() (*api.ActiveRequests, error) {
	return nil, nil
}
