package gencache

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConcurrentLoad(t *testing.T) {
	for _, stale := range []time.Duration{0, time.Minute} {
		t.Run(stale.String(), func(t *testing.T) {
			calls := 0 // Fetch is serialized by the cache; only inspect after all callers finish.
			c := New(context.Background(), stale, func(context.Context) (int, error) {
				calls++
				return 42, nil
			})
			const readers = 20
			start := make(chan struct{})
			results := make(chan error, readers)
			var wg sync.WaitGroup
			for i := 0; i < readers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					for j := 0; j < 50; j++ {
						value, err := c.Load(context.Background())
						if err != nil {
							results <- err
							return
						}
						if value != 42 {
							results <- errors.New("incorrect cached value")
							return
						}
					}
				}()
			}
			close(start)
			wg.Wait()
			close(results)
			for err := range results {
				require.NoError(t, err)
			}
			if stale > 0 {
				require.Equal(t, 1, calls, "concurrent cold loads must share one fetch")
			}
		})
	}
}

func TestTouchFreshnessAndErrorRetry(t *testing.T) {
	ctx := context.Background()
	failure := errors.New("temporary fetch failure")
	calls := 0
	c := New(ctx, time.Minute, func(context.Context) (int, error) {
		calls++
		if calls == 1 {
			return 0, failure
		}
		return 42, nil
	})
	_, err := c.Load(ctx)
	require.ErrorIs(t, err, failure)
	// Set the age explicitly rather than sleeping for the five-second retry interval.
	c.fetched = time.Now().Add(-time.Second * 6)
	value, err := c.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, 42, value)
	c.fetched = time.Now().Add(-time.Second * 40)
	c.Touch(ctx, 0)
	require.Equal(t, 2, calls, "fresh values should remain cached")
	c.Touch(ctx, time.Second*30)
	require.Equal(t, 3, calls, "eager window should refresh nearly stale values")
}
