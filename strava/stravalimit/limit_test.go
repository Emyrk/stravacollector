package stravalimit

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInterval(t *testing.T) {
	t.Parallel()

	ti := time.Now()
	start := GetInterval(ti)
	for i := 0; i < 50; i++ {
		ti = ti.Add(time.Minute * 15)

		next := GetInterval(ti)
		require.Truef(t, next == start+1, "next: %d, start: %d", next, start)
		start = next
	}
}

func TestDay(t *testing.T) {
	t.Parallel()

	ti := time.Now()
	start := GetDay(ti)
	for i := 0; i < 400; i++ {
		ti = ti.Add(time.Hour * 24)

		next := GetDay(ti)
		require.Truef(t, next != start, "next: %d, start: %d", next, start)
		start = next
	}
}

func TestTimeToMidnight(t *testing.T) {
	now := time.Now()
	fmt.Println(NextDailyReset(now))
}
