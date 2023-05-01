package gencache

import (
	"context"
	"sync"
	"time"
)

type LazyCache[T any] struct {
	Stale time.Duration
	Fetch func(ctx context.Context) (T, error)

	sync.RWMutex

	lastVal   T
	lastError error
	fetched   time.Time
}

func New[T any](stale time.Duration, fetch func(ctx context.Context) (T, error)) *LazyCache[T] {
	return &LazyCache[T]{
		Stale: stale,
		Fetch: fetch,
	}
}

func (c *LazyCache[T]) Load(ctx context.Context) (T, error) {
	if time.Since(c.fetched) > c.Stale {
		c.updateCache(ctx)
	}

	c.RLock()
	defer c.RUnlock()
	return c.lastVal, c.lastError
}

func (c *LazyCache[T]) updateCache(ctx context.Context) {
	c.Lock()
	defer c.Unlock()

	val, err := c.Fetch(ctx)
	c.lastError = err
	c.lastVal = val
	c.fetched = time.Now()
}
