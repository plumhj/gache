package gache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGache_SetGetDel(t *testing.T) {

	t.Run("int", func(t *testing.T) {
		cache := New[int]()
		cache.Set("A", 1)
		v, ok := cache.Get("A")
		assert.Equal(t, true, ok)
		assert.Equal(t, 1, v)
		cache.Del("A")
		v, ok = cache.Get("A")
		assert.Equal(t, false, ok)
	})

	t.Run("string", func(t *testing.T) {
		cache := New[string]()
		cache.Set("A", "A")
		v, ok := cache.Get("A")
		assert.Equal(t, true, ok)
		assert.Equal(t, "A", v)
		cache.Del("A")
		v, ok = cache.Get("A")
		assert.Equal(t, false, ok)
	})
}

func TestGache_Inc(t *testing.T) {

	t.Run("int", func(t *testing.T) {
		cache := New[int]()
		v := cache.Inc("A", 1)
		assert.Equal(t, 1, v)
		v = cache.Inc("A", 1)
		assert.Equal(t, 2, v)
		v = cache.Inc("A", 2)
		assert.Equal(t, 4, v)
		v, ok := cache.Get("A")
		assert.Equal(t, true, ok)
		assert.Equal(t, 4, v)
	})

	t.Run("float", func(t *testing.T) {
		cache := New[float64]()
		v := cache.Inc("A", 1)
		assert.Equal(t, 1.0, v)
		v = cache.Inc("A", 1)
		assert.Equal(t, 2.0, v)
		v = cache.Inc("A", 2)
		assert.Equal(t, 4.0, v)
	})

	t.Run("string", func(t *testing.T) {
		cache := New[string]()
		v := cache.Inc("A", "1")
		assert.Equal(t, "1", v)
		v = cache.Inc("A", "1")
		assert.Equal(t, "11", v)
		v = cache.Inc("A", "2")
		assert.Equal(t, "112", v)
	})
}

func TestGache_TTL(t *testing.T) {
	cache := New[int](OptionTTLUnit[int](time.Millisecond))
	cache.Set("A", 1, 10)
	v, ok := cache.Get("A")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)
	time.Sleep(time.Millisecond * 10)
	v, ok = cache.Get("A")
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, v)
}

func TestGache_CleanUpAndEviction(t *testing.T) {

	t.Run("clean up", func(t *testing.T) {
		cache := New[int](OptionTTLUnit[int](time.Millisecond), OptionCleanup[int](time.Millisecond))
		cache.Set("A", 1, 1)
		time.Sleep(time.Millisecond * 10)
		_, ok := cache.items["A"]
		assert.Equal(t, false, ok)
	})

	t.Run("eviction", func(t *testing.T) {
		cache := New[int](OptionTTLUnit[int](time.Millisecond), OptionCleanup[int](time.Millisecond))
		evicted := false
		cache.Set("A", 1, 1)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		cache.OnEviction = func(k string, v int) {
			assert.Equal(t, "A", k)
			assert.Equal(t, 1, v)
			evicted = true
			cancel()
		}

		<-ctx.Done()
		assert.Equal(t, true, evicted)
	})
}
