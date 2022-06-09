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
package random

import (
	cryptorand "crypto/rand"
	"io"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// Rand is a global *rand.Rand instance, which initialized with NewSource() source.
var Rand = rand.New(NewSource())

// Reader is a global, shared instance of a pseudorandom bytes generator.
// It doesn't consume entropy.
var Reader io.Reader = &reader{rnd: Rand}

// copypaste from standard math/rand
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

// NewSource returns math/rand.Source safe for concurrent use and initialized
// with current unix-nano timestamp
func NewSource() rand.Source {
	var seed int64
	if cryptoseed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64)); err != nil {
		// This should not happen, but worst-case fallback to time-based seed.
		seed = time.Now().UnixNano()
	} else {
		seed = cryptoseed.Int64()
	}
	return &lockedSource{
		src: rand.NewSource(seed),
	}
}

type reader struct {
	rnd *rand.Rand
}

func (r *reader) Read(b []byte) (int, error) {
	i := 0
	for {
		val := r.rnd.Int63()
		for val > 0 {
			b[i] = byte(val)
			i++
			if i == len(b) {
				return i, nil
			}
			val >>= 8
		}
	}
}
