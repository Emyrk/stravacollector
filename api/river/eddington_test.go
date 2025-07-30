package river_test

import (
	"testing"

	"github.com/Emyrk/strava/api/river"
	"github.com/stretchr/testify/require"
)

func TestEddingtonNumber(t *testing.T) {
	t.Parallel()

	t.Run("Simple", func(t *testing.T) {
		e := river.EddingtonNumbers{}

		e.Add(5)
		require.Equal(t, e, river.EddingtonNumbers{1, 1, 1, 1, 1})
		require.Equal(t, e.Current(), int32(1))

		e.Add(3)
		require.Equal(t, e, river.EddingtonNumbers{2, 2, 2, 1, 1})
		require.Equal(t, e.Current(), int32(2))

		e.Add(7)
		require.Equal(t, e, river.EddingtonNumbers{3, 3, 3, 2, 2, 1, 1})
		require.Equal(t, e.Current(), int32(3))

		e.Add(0)
		require.Equal(t, e, river.EddingtonNumbers{3, 3, 3, 2, 2, 1, 1})

		e.Add(7)
		require.Equal(t, e, river.EddingtonNumbers{4, 4, 4, 3, 3, 2, 2})
		require.Equal(t, e.Current(), int32(3))

		e.Add(4)
		require.Equal(t, e, river.EddingtonNumbers{5, 5, 5, 4, 3, 2, 2})
		require.Equal(t, e.Current(), int32(4))
	})
}
