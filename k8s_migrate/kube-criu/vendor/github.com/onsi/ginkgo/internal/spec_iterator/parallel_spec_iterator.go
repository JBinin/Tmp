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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/onsi/ginkgo/internal/spec"
)

type ParallelIterator struct {
	specs  []*spec.Spec
	host   string
	client *http.Client
}

func NewParallelIterator(specs []*spec.Spec, host string) *ParallelIterator {
	return &ParallelIterator{
		specs:  specs,
		host:   host,
		client: &http.Client{},
	}
}

func (s *ParallelIterator) Next() (*spec.Spec, error) {
	resp, err := s.client.Get(s.host + "/counter")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unexpected status code %d", resp.StatusCode))
	}

	var counter Counter
	err = json.NewDecoder(resp.Body).Decode(&counter)
	if err != nil {
		return nil, err
	}

	if counter.Index >= len(s.specs) {
		return nil, ErrClosed
	}

	return s.specs[counter.Index], nil
}

func (s *ParallelIterator) NumberOfSpecsPriorToIteration() int {
	return len(s.specs)
}

func (s *ParallelIterator) NumberOfSpecsToProcessIfKnown() (int, bool) {
	return -1, false
}

func (s *ParallelIterator) NumberOfSpecsThatWillBeRunIfKnown() (int, bool) {
	return -1, false
}
