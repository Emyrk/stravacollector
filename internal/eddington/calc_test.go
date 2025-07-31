package eddington_test

import (
	"fmt"
	"testing"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/internal/eddington"
	"github.com/stretchr/testify/require"
)

func TestDistance(t *testing.T) {
	mi := database.DistanceToMiles(144840)
	fmt.Println(int(mi))
}

func TestEddingtonNumber(t *testing.T) {
	t.Parallel()

	t.Run("Simple", func(t *testing.T) {
		e := eddington.Sums{}

		e.Add(5)
		require.Equal(t, e, eddington.Sums{1, 1, 1, 1, 1})
		require.Equal(t, e.Current(), int32(1))

		e.Add(3)
		require.Equal(t, e, eddington.Sums{2, 2, 2, 1, 1})
		require.Equal(t, e.Current(), int32(2))

		e.Add(7)
		require.Equal(t, e, eddington.Sums{3, 3, 3, 2, 2, 1, 1})
		require.Equal(t, e.Current(), int32(3))

		e.Add(0)
		require.Equal(t, e, eddington.Sums{3, 3, 3, 2, 2, 1, 1})

		e.Add(7)
		require.Equal(t, e, eddington.Sums{4, 4, 4, 3, 3, 2, 2})
		require.Equal(t, e.Current(), int32(3))

		e.Add(4)
		require.Equal(t, e, eddington.Sums{5, 5, 5, 4, 3, 2, 2})
		require.Equal(t, e.Current(), int32(4))
	})
}
