package debounce_test

import (
	"testing"
	"time"

	"github.com/Emyrk/strava/internal/debounce"
	"github.com/stretchr/testify/require"
)

func TestDebouncer(t *testing.T) {
	t.Parallel()

	d := debounce.New(time.Second)

	start := time.Now()
	c := 0
	for {
		d.Do(func() {
			c++
		})
		if time.Since(start) > time.Second*4 {
			break
		}
	}

	require.Equal(t, c, 4)
}
