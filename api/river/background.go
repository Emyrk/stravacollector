package river

import (
	"context"
	"time"
)

func (m *Manager) background(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		detailed, err := m.db.TotalActivityDetailsCount(ctx)
		if err == nil {
			m.rideActivityDetails.Set(float64(detailed))
		}
		summaries, err := m.db.TotalRideActivitySummariesCount(ctx)
		if err == nil {
			m.rideActivitySummaries.Set(float64(summaries))
		}
		time.Sleep(time.Minute * 8)
	}()
}
