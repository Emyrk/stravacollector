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
	once      sync.Once
}

func New[T any](ctx context.Context, stale time.Duration, fetch func(ctx context.Context) (T, error)) *LazyCache[T] {
	c := &LazyCache[T]{
		Stale: stale,
		Fetch: fetch,
		once:  sync.Once{},
	}
	c.RunEagerLoader(ctx)
	return c
}

func (c *LazyCache[T]) RunEagerLoader(ctx context.Context) {
	c.once.Do(func() {
		ticker := time.NewTicker(c.Stale / 2)
		go func() {
			for {
				select {
				case <-ticker.C:
					c.Touch(ctx)
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()
	})
}

func (c *LazyCache[T]) Touch(ctx context.Context) {
	stale := c.Stale
	// Errors should not persist for very long
	if c.lastError != nil {
		stale = time.Second * 5
	}

	if time.Since(c.fetched) > stale {
		c.updateCache(ctx)
	}
}

func (c *LazyCache[T]) Load(ctx context.Context) (T, error) {
	c.Touch(ctx)

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
