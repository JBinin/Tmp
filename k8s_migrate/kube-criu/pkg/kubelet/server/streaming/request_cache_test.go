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
/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package streaming

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/clock"
)

func TestInsert(t *testing.T) {
	c, _ := newTestCache()

	// Insert normal
	oldestTok, err := c.Insert(nextRequest())
	require.NoError(t, err)
	assert.Len(t, oldestTok, TokenLen)
	assertCacheSize(t, c, 1)

	// Insert until full
	for i := 0; i < MaxInFlight-2; i++ {
		tok, err := c.Insert(nextRequest())
		require.NoError(t, err)
		assert.Len(t, tok, TokenLen)
	}
	assertCacheSize(t, c, MaxInFlight-1)

	newestReq := nextRequest()
	newestTok, err := c.Insert(newestReq)
	require.NoError(t, err)
	assert.Len(t, newestTok, TokenLen)
	assertCacheSize(t, c, MaxInFlight)
	require.Contains(t, c.tokens, oldestTok, "oldest request should still be cached")

	// Consume newest token.
	req, ok := c.Consume(newestTok)
	assert.True(t, ok, "newest request should still be cached")
	assert.Equal(t, newestReq, req)
	require.Contains(t, c.tokens, oldestTok, "oldest request should still be cached")

	// Insert again (still full)
	tok, err := c.Insert(nextRequest())
	require.NoError(t, err)
	assert.Len(t, tok, TokenLen)
	assertCacheSize(t, c, MaxInFlight)

	// Insert again (should evict)
	_, err = c.Insert(nextRequest())
	assert.Error(t, err, "should reject further requests")
	errResponse := httptest.NewRecorder()
	require.NoError(t, WriteError(err, errResponse))
	assert.Equal(t, errResponse.Code, http.StatusTooManyRequests)
	assert.Equal(t, strconv.Itoa(int(CacheTTL.Seconds())), errResponse.HeaderMap.Get("Retry-After"))

	assertCacheSize(t, c, MaxInFlight)
	_, ok = c.Consume(oldestTok)
	assert.True(t, ok, "oldest request should be valid")
}

func TestConsume(t *testing.T) {
	c, clock := newTestCache()

	{ // Insert & consume.
		req := nextRequest()
		tok, err := c.Insert(req)
		require.NoError(t, err)
		assertCacheSize(t, c, 1)

		cachedReq, ok := c.Consume(tok)
		assert.True(t, ok)
		assert.Equal(t, req, cachedReq)
		assertCacheSize(t, c, 0)
	}

	{ // Insert & consume out of order
		req1 := nextRequest()
		tok1, err := c.Insert(req1)
		require.NoError(t, err)
		assertCacheSize(t, c, 1)

		req2 := nextRequest()
		tok2, err := c.Insert(req2)
		require.NoError(t, err)
		assertCacheSize(t, c, 2)

		cachedReq2, ok := c.Consume(tok2)
		assert.True(t, ok)
		assert.Equal(t, req2, cachedReq2)
		assertCacheSize(t, c, 1)

		cachedReq1, ok := c.Consume(tok1)
		assert.True(t, ok)
		assert.Equal(t, req1, cachedReq1)
		assertCacheSize(t, c, 0)
	}

	{ // Consume a second time
		req := nextRequest()
		tok, err := c.Insert(req)
		require.NoError(t, err)
		assertCacheSize(t, c, 1)

		cachedReq, ok := c.Consume(tok)
		assert.True(t, ok)
		assert.Equal(t, req, cachedReq)
		assertCacheSize(t, c, 0)

		_, ok = c.Consume(tok)
		assert.False(t, ok)
		assertCacheSize(t, c, 0)
	}

	{ // Consume without insert
		_, ok := c.Consume("fooBAR")
		assert.False(t, ok)
		assertCacheSize(t, c, 0)
	}

	{ // Consume expired
		tok, err := c.Insert(nextRequest())
		require.NoError(t, err)
		assertCacheSize(t, c, 1)

		clock.Step(2 * CacheTTL)

		_, ok := c.Consume(tok)
		assert.False(t, ok)
		assertCacheSize(t, c, 0)
	}
}

func TestGC(t *testing.T) {
	c, clock := newTestCache()

	// When empty
	c.gc()
	assertCacheSize(t, c, 0)

	tok1, err := c.Insert(nextRequest())
	require.NoError(t, err)
	assertCacheSize(t, c, 1)
	clock.Step(10 * time.Second)
	tok2, err := c.Insert(nextRequest())
	require.NoError(t, err)
	assertCacheSize(t, c, 2)

	// expired: tok1, tok2
	// non-expired: tok3, tok4
	clock.Step(2 * CacheTTL)
	tok3, err := c.Insert(nextRequest())
	require.NoError(t, err)
	assertCacheSize(t, c, 1)
	clock.Step(10 * time.Second)
	tok4, err := c.Insert(nextRequest())
	require.NoError(t, err)
	assertCacheSize(t, c, 2)

	_, ok := c.Consume(tok1)
	assert.False(t, ok)
	_, ok = c.Consume(tok2)
	assert.False(t, ok)
	_, ok = c.Consume(tok3)
	assert.True(t, ok)
	_, ok = c.Consume(tok4)
	assert.True(t, ok)

	// When full, nothing is expired.
	for i := 0; i < MaxInFlight; i++ {
		_, err := c.Insert(nextRequest())
		require.NoError(t, err)
	}
	assertCacheSize(t, c, MaxInFlight)

	// When everything is expired
	clock.Step(2 * CacheTTL)
	_, err = c.Insert(nextRequest())
	require.NoError(t, err)
	assertCacheSize(t, c, 1)
}

func newTestCache() (*requestCache, *clock.FakeClock) {
	c := newRequestCache()
	fakeClock := clock.NewFakeClock(time.Now())
	c.clock = fakeClock
	return c, fakeClock
}

func assertCacheSize(t *testing.T, cache *requestCache, expectedSize int) {
	tokenLen := len(cache.tokens)
	llLen := cache.ll.Len()
	assert.Equal(t, tokenLen, llLen, "inconsistent cache size! len(tokens)=%d; len(ll)=%d", tokenLen, llLen)
	assert.Equal(t, expectedSize, tokenLen, "unexpected cache size!")
}

var requestUID = 0

func nextRequest() interface{} {
	requestUID++
	return requestUID
}
