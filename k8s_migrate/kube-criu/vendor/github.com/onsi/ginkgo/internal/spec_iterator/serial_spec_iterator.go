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
package spec_iterator

import (
	"github.com/onsi/ginkgo/internal/spec"
)

type SerialIterator struct {
	specs []*spec.Spec
	index int
}

func NewSerialIterator(specs []*spec.Spec) *SerialIterator {
	return &SerialIterator{
		specs: specs,
		index: 0,
	}
}

func (s *SerialIterator) Next() (*spec.Spec, error) {
	if s.index >= len(s.specs) {
		return nil, ErrClosed
	}

	spec := s.specs[s.index]
	s.index += 1
	return spec, nil
}

func (s *SerialIterator) NumberOfSpecsPriorToIteration() int {
	return len(s.specs)
}

func (s *SerialIterator) NumberOfSpecsToProcessIfKnown() (int, bool) {
	return len(s.specs), true
}

func (s *SerialIterator) NumberOfSpecsThatWillBeRunIfKnown() (int, bool) {
	count := 0
	for _, s := range s.specs {
		if !s.Skipped() && !s.Pending() {
			count += 1
		}
	}
	return count, true
}
