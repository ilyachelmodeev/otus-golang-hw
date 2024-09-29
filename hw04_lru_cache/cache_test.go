package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		// cache - aaa bbb ccc
		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		// cache - bbb ccc ddd
		_, ok := c.Get("aaa")
		require.False(t, ok)

		// cache - bbb ccc ddd
		wasInCache = c.Set("bbb", 200000)
		require.True(t, wasInCache)

		// cache - ddd bbb ccc
		wasInCache = c.Set("ddd", 400000)
		require.True(t, wasInCache)

		val, ok := c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400000, val)

		// cache - eee ddd bbb
		wasInCache = c.Set("eee", 100000)
		require.False(t, wasInCache)

		// cache - eee ddd bbbb
		_, ok = c.Get("ccc")
		require.False(t, ok)

		// cache - empty
		c.Clear()

		// cache - eee
		wasInCache = c.Set("eee", 100000)
		require.False(t, wasInCache)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
